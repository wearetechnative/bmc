package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsconfig"
	"github.com/wearetechnative/bmc/internal/awsops"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/mfa"
)

var (
	consoleProfile string
	consoleService string
)

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Open Firefox with AWS account in console",
	RunE:  runConsole,
}

func init() {
	consoleCmd.Flags().StringVarP(&consoleProfile, "profile", "p", "", "AWS profile to use")
	consoleCmd.Flags().StringVarP(&consoleService, "service", "s", "", "AWS service to open (e.g. ec2, s3)")
	rootCmd.AddCommand(consoleCmd)
}

func runConsole(cmd *cobra.Command, args []string) error {
	profiles, err := awsconfig.LoadProfiles()
	if err != nil {
		return err
	}

	var selectedProfile awsconfig.Profile

	// Use existing AWS_PROFILE if set and no -p flag
	if consoleProfile == "" {
		envProfile := os.Getenv("AWS_PROFILE")
		if envProfile != "" {
			p, ok := awsconfig.FindProfile(profiles, envProfile)
			if !ok {
				return fmt.Errorf("AWS_PROFILE=%q not found in config", envProfile)
			}
			selectedProfile = p
		} else {
			selectedProfile, err = selectProfileInteractive(profiles)
			if err != nil {
				return err
			}
			if selectedProfile.Name == "" {
				return nil
			}
		}
	} else {
		p, ok := awsconfig.FindProfile(profiles, consoleProfile)
		if !ok {
			return fmt.Errorf("profile %q not found", consoleProfile)
		}
		selectedProfile = p
	}

	sourceProfile, err := awsconfig.ResolveSourceProfile(selectedProfile)
	if err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if err := mfa.EnsureValid(sourceProfile, cfg, os.Stderr); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "-- Opening console for profile: %s\n", selectedProfile.Name)
	return awsops.OpenConsole(selectedProfile.Name, consoleService)
}
