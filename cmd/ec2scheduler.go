package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsops"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/ui"
)

var ec2schedulerCmd = &cobra.Command{
	Use:   "ec2scheduler",
	Short: "Toggle InstanceScheduler tag to enable/disable scheduling",
	RunE:  runEC2Scheduler,
}

func init() {
	ec2schedulerCmd.Flags().StringVarP(&globalProfile, "profile", "p", "", "AWS profile to use (omit value to force interactive selection)")
	ec2schedulerCmd.Flags().Lookup("profile").NoOptDefVal = " "
	rootCmd.AddCommand(ec2schedulerCmd)
}

func runEC2Scheduler(cmd *cobra.Command, args []string) error {
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

	// Determine current scheduler state
	var currentScheduler string
	for _, i := range instances {
		if i.InstanceID == instanceID {
			currentScheduler = i.Scheduler
			break
		}
	}

	var enable bool
	if currentScheduler == "yes" {
		confirm, err := ui.Confirm(fmt.Sprintf("Disable InstanceScheduler tag for %s?", instanceID))
		if err != nil {
			return err
		}
		if !confirm {
			return nil
		}
		enable = false
	} else {
		confirm, err := ui.Confirm(fmt.Sprintf("Enable InstanceScheduler tag for %s?", instanceID))
		if err != nil {
			return err
		}
		if !confirm {
			return nil
		}
		enable = true
	}

	if err := awsops.ToggleSchedulerTag(profile, instanceID, enable); err != nil {
		return err
	}

	action := "enabled"
	if !enable {
		action = "disabled"
	}
	fmt.Printf("InstanceScheduler %s for %s\n", action, instanceID)
	return nil
}
