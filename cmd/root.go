package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bmc",
	Short: "Bill McCloud's Toolbox — AWS CLI helper",
	Long:  `bmc is an AWS toolbox for profile selection, EC2/ECS operations, and console access.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = false
}
