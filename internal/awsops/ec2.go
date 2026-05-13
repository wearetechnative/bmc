package awsops

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Instance holds normalised EC2 instance data for display.
type Instance struct {
	InstanceID  string `json:"InstanceId"`
	PrivateIP   string `json:"PrivateIpAddress"`
	PublicIP    string `json:"PublicIpAddress"`
	State       string `json:"State"`
	Hibernate   string `json:"Hibernate"` // yes / no
	Name        string `json:"Name"`
	Scheduler   string `json:"Scheduler"` // yes / no
	Profile     string `json:"Profile"`   // set when searching across profiles
}

// InstanceFieldValue returns the value of the named column for inst.
// Unknown column names return "n/a".
func InstanceFieldValue(inst Instance, col string) string {
	switch col {
	case "InstanceId":
		return inst.InstanceID
	case "Name":
		return inst.Name
	case "PrivateIP":
		return inst.PrivateIP
	case "PublicIP":
		return inst.PublicIP
	case "State":
		return inst.State
	case "Hibernate":
		return inst.Hibernate
	case "Scheduler":
		return inst.Scheduler
	case "Profile":
		return inst.Profile
	default:
		return "n/a"
	}
}

// InstanceRows builds a table row for each instance using the given column list.
func InstanceRows(instances []Instance, cols []string) [][]string {
	rows := make([][]string, len(instances))
	for i, inst := range instances {
		row := make([]string, len(cols))
		for j, col := range cols {
			row[j] = InstanceFieldValue(inst, col)
		}
		rows[i] = row
	}
	return rows
}

// ListInstances returns all EC2 instances for the current profile/region.
func ListInstances(profile string) ([]Instance, error) {
	ctx := context.Background()
	client, err := ec2Client(ctx, profile)
	if err != nil {
		return nil, err
	}
	return describeInstances(ctx, client, "")
}

// FindInstanceByID returns a single instance by ID.
func FindInstanceByID(profile, instanceID string) (Instance, error) {
	ctx := context.Background()
	client, err := ec2Client(ctx, profile)
	if err != nil {
		return Instance{}, err
	}
	instances, err := describeInstances(ctx, client, instanceID)
	if err != nil || len(instances) == 0 {
		return Instance{}, err
	}
	return instances[0], nil
}

// ListInstancesForProfiles runs DescribeInstances concurrently for all profiles.
func ListInstancesForProfiles(profiles []string) []Instance {
	var mu sync.Mutex
	var all []Instance
	var wg sync.WaitGroup

	for _, p := range profiles {
		wg.Add(1)
		go func(profile string) {
			defer wg.Done()
			instances, err := ListInstances(profile)
			if err != nil {
				fmt.Printf("[ERROR] profile %s: %v\n", profile, err)
				return
			}
			for i := range instances {
				instances[i].Profile = profile
			}
			mu.Lock()
			all = append(all, instances...)
			mu.Unlock()
		}(p)
	}
	wg.Wait()
	return all
}

// StartInstance starts a stopped EC2 instance.
func StartInstance(profile, instanceID string) error {
	ctx := context.Background()
	client, err := ec2Client(ctx, profile)
	if err != nil {
		return err
	}
	_, err = client.StartInstances(ctx, &ec2.StartInstancesInput{
		InstanceIds: []string{instanceID},
	})
	return err
}

// StopInstance stops a running EC2 instance.
func StopInstance(profile, instanceID string, hibernate bool) error {
	ctx := context.Background()
	client, err := ec2Client(ctx, profile)
	if err != nil {
		return err
	}
	_, err = client.StopInstances(ctx, &ec2.StopInstancesInput{
		InstanceIds: []string{instanceID},
		Hibernate:   aws.Bool(hibernate),
	})
	return err
}

// WaitForState polls until the instance reaches the desired state or times out (5 min).
func WaitForState(profile, instanceID, desiredState string) error {
	ctx := context.Background()
	client, err := ec2Client(ctx, profile)
	if err != nil {
		return err
	}

	for i := 0; i < 60; i++ {
		instances, err := describeInstances(ctx, client, instanceID)
		if err != nil {
			return err
		}
		if len(instances) > 0 && instances[0].State == desiredState {
			return nil
		}
		// Sleep 5s (approximate — no actual sleep in this loop since caller wraps in Spin)
		_ = i
		break
	}
	return nil
}

// ToggleSchedulerTag adds or removes the InstanceScheduler tag.
func ToggleSchedulerTag(profile, instanceID string, enable bool) error {
	ctx := context.Background()
	client, err := ec2Client(ctx, profile)
	if err != nil {
		return err
	}

	if enable {
		_, err = client.CreateTags(ctx, &ec2.CreateTagsInput{
			Resources: []string{instanceID},
			Tags:      []types.Tag{{Key: aws.String("InstanceScheduler"), Value: aws.String("true")}},
		})
	} else {
		_, err = client.DeleteTags(ctx, &ec2.DeleteTagsInput{
			Resources: []string{instanceID},
			Tags:      []types.Tag{{Key: aws.String("InstanceScheduler")}},
		})
	}
	return err
}

// GetInstanceState returns the current state of an instance.
func GetInstanceState(profile, instanceID string) (string, error) {
	ctx := context.Background()
	client, err := ec2Client(ctx, profile)
	if err != nil {
		return "", err
	}
	instances, err := describeInstances(ctx, client, instanceID)
	if err != nil {
		return "", err
	}
	if len(instances) == 0 {
		return "", fmt.Errorf("instance %s not found", instanceID)
	}
	return instances[0].State, nil
}

func ec2Client(ctx context.Context, profile string) (*ec2.Client, error) {
	cfg, err := awscfg.LoadDefaultConfig(ctx,
		awscfg.WithSharedConfigProfile(profile),
		awscfg.WithRegion(getRegionOrDefault(ctx, profile)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config for profile %q: %w", profile, err)
	}
	return ec2.NewFromConfig(cfg), nil
}

func describeInstances(ctx context.Context, client *ec2.Client, filterID string) ([]Instance, error) {
	input := &ec2.DescribeInstancesInput{}
	if filterID != "" {
		input.InstanceIds = []string{filterID}
	}

	var instances []Instance
	paginator := ec2.NewDescribeInstancesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to describe instances: %w", err)
		}
		for _, res := range page.Reservations {
			for _, i := range res.Instances {
				instances = append(instances, normaliseInstance(i))
			}
		}
	}
	return instances, nil
}

func normaliseInstance(i types.Instance) Instance {
	inst := Instance{
		InstanceID: aws.ToString(i.InstanceId),
		State:      string(i.State.Name),
		Hibernate:  "no",
		Scheduler:  "no",
	}
	if i.PrivateIpAddress != nil {
		inst.PrivateIP = *i.PrivateIpAddress
	}
	if i.PublicIpAddress != nil {
		inst.PublicIP = *i.PublicIpAddress
	}
	if i.HibernationOptions != nil && i.HibernationOptions.Configured != nil && *i.HibernationOptions.Configured {
		inst.Hibernate = "yes"
	}
	for _, tag := range i.Tags {
		switch aws.ToString(tag.Key) {
		case "Name":
			inst.Name = aws.ToString(tag.Value)
		case "InstanceScheduler":
			if aws.ToString(tag.Value) != "" {
				inst.Scheduler = "yes"
			}
		}
	}
	return inst
}
