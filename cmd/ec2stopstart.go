package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsops"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/ui"
)

var ec2stopstartCmd = &cobra.Command{
	Use:   "ec2stopstart",
	Short: "Stop/start an EC2 instance",
	RunE:  runEC2StopStart,
}

func init() {
	ec2stopstartCmd.Flags().StringVarP(&globalProfile, "profile", "p", "", "AWS profile to use (omit value to force interactive selection)")
	ec2stopstartCmd.Flags().Lookup("profile").NoOptDefVal = " "
	rootCmd.AddCommand(ec2stopstartCmd)
}

func runEC2StopStart(cmd *cobra.Command, args []string) error {
	profile, err := ensureAWSProfile()
	if err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	instances, err := awsops.ListInstances(profile)
	if err != nil {
		return err
	}

	instanceID, err := selectInstanceID(instances, cfg.EC2.Columns)
	if err != nil {
		return err
	}
	if instanceID == "" {
		return nil
	}

	state, err := awsops.GetInstanceState(profile, instanceID)
	if err != nil {
		return err
	}

	switch state {
	case "stopped":
		return startInstance(profile, instanceID)
	case "running":
		return stopInstance(profile, instanceID, instances)
	default:
		fmt.Fprintf(os.Stderr, "Instance %s not in running/stopped state.\nCurrent state: %s\n", instanceID, state)
		return nil
	}
}

func startInstance(profile, instanceID string) error {
	if err := awsops.StartInstance(profile, instanceID); err != nil {
		return err
	}
	return ui.Spin(fmt.Sprintf("Starting instance %s", instanceID), func() error {
		return waitForState(profile, instanceID, "running")
	})
}

func stopInstance(profile, instanceID string, instances []awsops.Instance) error {
	// Check if hibernation is supported
	var hibernateAvailable bool
	for _, i := range instances {
		if i.InstanceID == instanceID && i.Hibernate == "yes" {
			hibernateAvailable = true
			break
		}
	}

	options := []ui.Item{{Title: "stop"}, {Title: "exit menu"}}
	if hibernateAvailable {
		options = append([]ui.Item{{Title: "hibernate"}}, options...)
	}

	method, err := ui.Choose(fmt.Sprintf("Choose stop-method for instance: %s", instanceID), options)
	if err != nil {
		return err
	}

	switch method {
	case "hibernate":
		if err := awsops.StopInstance(profile, instanceID, true); err != nil {
			return err
		}
		return ui.Spin(fmt.Sprintf("Hibernating instance %s", instanceID), func() error {
			return waitForState(profile, instanceID, "stopped")
		})
	case "stop":
		if err := awsops.StopInstance(profile, instanceID, false); err != nil {
			return err
		}
		return ui.Spin(fmt.Sprintf("Stopping instance %s", instanceID), func() error {
			return waitForState(profile, instanceID, "stopped")
		})
	}
	return nil
}

func waitForState(profile, instanceID, desiredState string) error {
	for i := 0; i < 60; i++ {
		state, err := awsops.GetInstanceState(profile, instanceID)
		if err != nil {
			return err
		}
		if state == desiredState {
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("instance %s did not reach state %s within 5 minutes", instanceID, desiredState)
}

