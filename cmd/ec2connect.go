package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsops"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/prereqs"
	"github.com/wearetechnative/bmc/internal/ui"
)

var (
	ec2connectInstanceID string
	ec2connectUser       string
	ec2connectKey        string
)

var ec2connectCmd = &cobra.Command{
	Use:   "ec2connect [search]",
	Short: "Connect to a running EC2 instance",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runEC2Connect,
}

func init() {
	ec2connectCmd.Flags().StringVarP(&ec2connectInstanceID, "instance", "i", "", "Instance ID to connect to")
	ec2connectCmd.Flags().StringVarP(&ec2connectUser, "user", "u", "", "SSH user")
	ec2connectCmd.Flags().StringVarP(&ec2connectKey, "key", "k", "", "SSH identity file (passed to ssh -i)")
	ec2connectCmd.Flags().StringVarP(&globalProfile, "profile", "p", "", "AWS profile to use (omit value to force interactive selection)")
	ec2connectCmd.Flags().Lookup("profile").NoOptDefVal = " "
	rootCmd.AddCommand(ec2connectCmd)
}

func runEC2Connect(cmd *cobra.Command, args []string) error {
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

	// Select instance
	instanceID := ec2connectInstanceID
	if instanceID != "" && len(args) > 0 {
		fmt.Fprintf(os.Stderr, "Warning: positional argument %q ignored because -i flag is set\n", args[0])
	}
	if instanceID == "" {
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
	}

	// Check instance state
	state, err := awsops.GetInstanceState(profile, instanceID)
	if err != nil {
		return err
	}

	if state != "running" {
		if state == "stopped" {
			switch cfg.EC2.AutoStartStopped {
			case "always":
				if err := startInstance(profile, instanceID); err != nil {
					return err
				}
			case "never":
				return fmt.Errorf("instance %s is stopped (auto_start_stopped=never)", instanceID)
			default: // prompt
				ok, err := ui.Confirm(fmt.Sprintf("Instance %s is stopped. Start it?", instanceID))
				if err != nil {
					return err
				}
				if !ok {
					fmt.Fprintln(os.Stderr, "Instance not started. Exiting.")
					return nil
				}
				if err := startInstance(profile, instanceID); err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf("instance %s is not in running/stopped state (current: %s)", instanceID, state)
		}
	}

	// Choose connection method (skip if -u or known SSH intent)
	connectionMethod := ""
	if ec2connectUser != "" || ec2connectKey != "" {
		connectionMethod = "ssh"
	} else {
		method, err := ui.Choose("Connection method", []ui.Item{{Title: "ssh"}, {Title: "ssm"}})
		if err != nil {
			return err
		}
		connectionMethod = method
	}

	switch connectionMethod {
	case "ssh":
		return connectSSH(instanceID, ec2connectKey)
	case "ssm":
		return connectSSM(instanceID)
	}
	return nil
}

func connectSSH(instanceID, key string) error {
	if err := prereqs.Check(prereqs.SSH); err != nil {
		return err
	}

	user := ec2connectUser
	if user == "" {
		shellUsers := []ui.Item{
			{Title: "root"}, {Title: "ubuntu"}, {Title: "ec2-user"}, {Title: "other"},
		}
		selected, err := ui.Choose("Select SSH user", shellUsers)
		if err != nil {
			return err
		}
		if selected == "other" {
			selected, err = ui.Input("Enter username:", true)
			if err != nil {
				return err
			}
		}
		user = selected
	}

	if user == "" {
		return fmt.Errorf("no user selected")
	}

	sshBin, _ := prereqs.FindPath("ssh"), ""
	if sshBin == "" {
		sshBin = "ssh"
	}

	target := user + "@" + instanceID
	sshArgs := []string{"ssh"}
	if key != "" {
		sshArgs = append(sshArgs, "-i", key)
	}
	sshArgs = append(sshArgs, target)
	fmt.Fprintf(os.Stderr, "-- Executing: ssh %s\n", strings.Join(sshArgs[1:], " "))

	return syscall.Exec(sshBin, sshArgs, os.Environ())
}

func connectSSM(instanceID string) error {
	if err := prereqs.Check(prereqs.AWSCLI); err != nil {
		return err
	}
	if err := prereqs.Check(prereqs.SessionManagerPlugin); err != nil {
		return err
	}

	awsBin := prereqs.FindPath("aws")
	fmt.Fprintf(os.Stderr, "-- Executing: aws ssm start-session --target %s\n", instanceID)

	return syscall.Exec(awsBin, []string{"aws", "ssm", "start-session", "--target", instanceID}, os.Environ())
}
