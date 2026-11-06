package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/debug"
)

var rootCmd = &cobra.Command{
	Use:   "bmc",
	Short: "Bill McCloud's Toolbox — AWS CLI helper",
	Long:  `bmc is an AWS toolbox for profile selection, EC2/ECS operations, and console access.`,
}

func Execute() {
	// --version is an alias for the version command. Version is only known
	// here (set in main), so the flag is configured at execution time.
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate(versionBanner())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = false
	rootCmd.PersistentFlags().BoolVar(&debug.Enabled, "debug", false, "Show diagnostic output (e.g. MFA session details) on stderr")
}
