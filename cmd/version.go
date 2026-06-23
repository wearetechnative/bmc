package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// Version is set from the embedded VERSION-bmc file in main.go
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show bmc version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(`
    bmc v%s
    Bill McCloud's Toolbox

    http://github.com/wearetechnative/bmc

    by Wouter, Pim, et al.
    © Technative 2024-%d

`, Version, time.Now().Year())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
