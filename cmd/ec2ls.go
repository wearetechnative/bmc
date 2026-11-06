package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsops"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/ui"
)

var (
	ec2lsJSON      bool
	ec2lsTxt       bool
	ec2lsNoHeaders bool
)

var ec2lsCmd = &cobra.Command{
	Use:   "ec2ls",
	Short: "List running EC2 instances",
	RunE:  runEC2ls,
}

func init() {
	ec2lsCmd.Flags().BoolVar(&ec2lsJSON, "json", false, "Output instances as JSON array")
	ec2lsCmd.Flags().BoolVar(&ec2lsTxt, "txt", false, "Output instances as aligned plain text without borders")
	ec2lsCmd.Flags().BoolVar(&ec2lsNoHeaders, "no-headers", false, "Omit the header line (requires --txt)")
	ec2lsCmd.MarkFlagsMutuallyExclusive("json", "txt")
	addProfileFlags(ec2lsCmd)
	rootCmd.AddCommand(ec2lsCmd)
}

func runEC2ls(cmd *cobra.Command, args []string) error {
	if ec2lsNoHeaders && !ec2lsTxt {
		return fmt.Errorf("--no-headers requires --txt")
	}

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

	rows := awsops.InstanceRows(instances, cfg.EC2.Columns)
	if ec2lsTxt {
		ui.PrintText(cfg.EC2.Columns, rows, !ec2lsNoHeaders)
		return nil
	}
	ui.PrintTable(cfg.EC2.Columns, rows)
	return nil
}
