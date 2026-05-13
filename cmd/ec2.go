package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsops"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/ui"
)

var ec2Cmd = &cobra.Command{
	Use:   "ec2 [search]",
	Short: "Select an EC2 instance and perform an action",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runEC2,
}

func init() {
	rootCmd.AddCommand(ec2Cmd)
}

func runEC2(cmd *cobra.Command, args []string) error {
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

	// Select instance — with optional search filter
	var instanceID string
	if len(args) > 0 {
		fragment := strings.ToLower(args[0])
		var filtered []awsops.Instance
		for _, inst := range instances {
			combined := strings.ToLower(inst.InstanceID + inst.Name + inst.PrivateIP + inst.PublicIP)
			if strings.Contains(combined, fragment) {
				filtered = append(filtered, inst)
			}
		}
		switch len(filtered) {
		case 0:
			return fmt.Errorf("no instances found matching %q", args[0])
		case 1:
			instanceID = filtered[0].InstanceID
		default:
			instanceID, err = selectInstanceID(filtered, cfg.EC2.Columns)
			if err != nil {
				return err
			}
			if instanceID == "" {
				return nil
			}
		}
	} else {
		instanceID, err = selectInstanceID(instances, cfg.EC2.Columns)
		if err != nil {
			return err
		}
		if instanceID == "" {
			return nil
		}
	}

	// Determine current state to label the start/stop action
	state, err := awsops.GetInstanceState(profile, instanceID)
	if err != nil {
		return err
	}

	// Build action menu
	actions := []ui.Item{
		{Title: "Connect SSH"},
		{Title: "Connect SSM"},
	}
	switch state {
	case "running":
		actions = append(actions, ui.Item{Title: "Stop instance"})
	case "stopped":
		actions = append(actions, ui.Item{Title: "Start instance"})
	default:
		fmt.Fprintf(os.Stderr, "Note: instance is in state %q — start/stop not available\n", state)
	}
	actions = append(actions, ui.Item{Title: "Toggle scheduler"})

	action, err := ui.Choose(fmt.Sprintf("Action for %s", instanceID), actions)
	if err != nil {
		return err
	}
	if action == "" {
		return nil
	}

	switch action {
	case "Connect SSH":
		return connectSSH(instanceID)
	case "Connect SSM":
		return connectSSM(instanceID)
	case "Start instance":
		return startInstance(profile, instanceID)
	case "Stop instance":
		return stopInstance(profile, instanceID, instances)
	case "Toggle scheduler":
		return ec2ToggleScheduler(profile, instanceID, instances)
	}
	return nil
}

func ec2ToggleScheduler(profile, instanceID string, instances []awsops.Instance) error {
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
