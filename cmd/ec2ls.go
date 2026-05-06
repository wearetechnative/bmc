package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsops"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/ui"
)

var ec2lsCmd = &cobra.Command{
	Use:   "ec2ls",
	Short: "List running EC2 instances",
	RunE:  runEC2ls,
}

func init() {
	rootCmd.AddCommand(ec2lsCmd)
}

func runEC2ls(cmd *cobra.Command, args []string) error {
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

	return ui.ShowTable(cfg.EC2.Columns, awsops.InstanceRows(instances, cfg.EC2.Columns))
}
