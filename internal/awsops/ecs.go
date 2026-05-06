package awsops

import (
	"context"
	"fmt"
	"strings"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

// ECSCluster holds a cluster name.
type ECSCluster struct {
	Name string
}

// ECSService holds a service name.
type ECSService struct {
	Name string
}

// ECSTask holds a task short ID and ARN.
type ECSTask struct {
	ShortID string
	ARN     string
}

// ECSContainer holds a container name.
type ECSContainer struct {
	Name string
}

func ecsClient(ctx context.Context, profile string) (*ecs.Client, error) {
	cfg, err := awscfg.LoadDefaultConfig(ctx, awscfg.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config for profile %q: %w", profile, err)
	}
	return ecs.NewFromConfig(cfg), nil
}

// ListClusters returns all ECS cluster names.
func ListClusters(profile string) ([]ECSCluster, error) {
	ctx := context.Background()
	client, err := ecsClient(ctx, profile)
	if err != nil {
		return nil, err
	}

	resp, err := client.ListClusters(ctx, &ecs.ListClustersInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters: %w", err)
	}

	clusters := make([]ECSCluster, len(resp.ClusterArns))
	for i, arn := range resp.ClusterArns {
		clusters[i] = ECSCluster{Name: shortName(arn)}
	}
	return clusters, nil
}

// ListServices returns all services in a cluster.
func ListServices(profile, cluster string) ([]ECSService, error) {
	ctx := context.Background()
	client, err := ecsClient(ctx, profile)
	if err != nil {
		return nil, err
	}

	resp, err := client.ListServices(ctx, &ecs.ListServicesInput{
		Cluster: &cluster,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	services := make([]ECSService, len(resp.ServiceArns))
	for i, arn := range resp.ServiceArns {
		services[i] = ECSService{Name: shortName(arn)}
	}
	return services, nil
}

// ListRunningTasks returns running tasks for a service.
func ListRunningTasks(profile, cluster, service string) ([]ECSTask, error) {
	ctx := context.Background()
	client, err := ecsClient(ctx, profile)
	if err != nil {
		return nil, err
	}

	resp, err := client.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster:       &cluster,
		ServiceName:   &service,
		DesiredStatus: "RUNNING",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	tasks := make([]ECSTask, len(resp.TaskArns))
	for i, arn := range resp.TaskArns {
		tasks[i] = ECSTask{ShortID: shortName(arn), ARN: arn}
	}
	return tasks, nil
}

// ListContainers returns containers for a task.
func ListContainers(profile, cluster, taskARN string) ([]ECSContainer, error) {
	ctx := context.Background()
	client, err := ecsClient(ctx, profile)
	if err != nil {
		return nil, err
	}

	resp, err := client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: &cluster,
		Tasks:   []string{taskARN},
	})
	if err != nil || len(resp.Tasks) == 0 {
		return nil, fmt.Errorf("failed to describe task: %w", err)
	}

	containers := make([]ECSContainer, len(resp.Tasks[0].Containers))
	for i, c := range resp.Tasks[0].Containers {
		if c.Name != nil {
			containers[i] = ECSContainer{Name: *c.Name}
		}
	}
	return containers, nil
}

func shortName(arn string) string {
	parts := strings.Split(arn, "/")
	return parts[len(parts)-1]
}
