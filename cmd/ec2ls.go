package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsops"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/ui"
)

var ec2lsJSON bool

var ec2lsCmd = &cobra.Command{
	Use:   "ec2ls",
	Short: "List running EC2 instances",
	RunE:  runEC2ls,
}

func init() {
	ec2lsCmd.Flags().BoolVar(&ec2lsJSON, "json", false, "Output instances as JSON array")
	ec2lsCmd.Flags().StringVarP(&globalProfile, "profile", "p", "", "AWS profile to use (omit value to force interactive selection)")
	ec2lsCmd.Flags().Lookup("profile").NoOptDefVal = " "
	rootCmd.AddCommand(ec2lsCmd)
}

func runEC2ls(cmd *cobra.Command, args []string) error {
	profile, err := ensureAWSProfile()
	if err != nil {
		return err
	}

	instances, err := awsops.ListInstances(profile)
	if err != nil {
		return err
	}

	if ec2lsJSON {
		data, err := json.Marshal(instances)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ui.PrintTable(cfg.EC2.Columns, awsops.InstanceRows(instances, cfg.EC2.Columns))
	return nil
}
