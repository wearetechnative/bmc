package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsconfig"
	"github.com/wearetechnative/bmc/internal/awsops"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/history"
	"github.com/wearetechnative/bmc/internal/mfa"
	"github.com/wearetechnative/bmc/internal/ui"
)

var (
	consoleProfile string
	consoleService string
)

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Open Firefox with AWS account in console",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runConsole,
}

func init() {
	consoleCmd.Flags().StringVarP(&consoleProfile, "profile", "p", "", "AWS profile (omit value to force interactive selection)")
	consoleCmd.Flags().Lookup("profile").NoOptDefVal = " "
	consoleCmd.Flags().StringVarP(&consoleService, "service", "s", "", "AWS service to open (e.g. ec2, s3)")
	rootCmd.AddCommand(consoleCmd)
}

func runConsole(cmd *cobra.Command, args []string) error {
	profiles, err := awsconfig.LoadProfiles()
	if err != nil {
		return err
	}

	var selectedProfile awsconfig.Profile

	// When NoOptDefVal fires (-p without =value), cobra puts the next word in args.
	profileName := strings.TrimSpace(consoleProfile)
	if profileName == "" && len(args) > 0 {
		profileName = args[0]
	}

	interactive := false

	switch {
	case profileName != "":
		// -p <name>: use the given profile directly
		p, ok := awsconfig.FindProfile(profiles, profileName)
		if !ok {
			return fmt.Errorf("profile %q not found", profileName)
		}
		selectedProfile = p
	case cmd.Flags().Changed("profile"):
		// -p bare: force interactive selection (ignore AWS_PROFILE)
		selectedProfile, err = selectProfileForConsole(profiles)
		if err != nil {
			return err
		}
		if selectedProfile.Name == "" {
			return nil
		}
		interactive = true
	default:
		// no -p: use AWS_PROFILE if set, otherwise interactive
		envProfile := os.Getenv("AWS_PROFILE")
		if envProfile != "" {
			p, ok := awsconfig.FindProfile(profiles, envProfile)
			if !ok {
				return fmt.Errorf("AWS_PROFILE=%q not found in config", envProfile)
			}
			selectedProfile = p
		} else {
			selectedProfile, err = selectProfileForConsole(profiles)
			if err != nil {
				return err
			}
			if selectedProfile.Name == "" {
				return nil
			}
			interactive = true
		}
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
	if err := awsops.OpenConsole(selectedProfile.Name, consoleService, cfg); err != nil {
		return err
	}

	if interactive {
		_ = history.Save("console", selectedProfile.Name)
	}
	return nil
}

// selectProfileForConsole shows a flat profile list with recent profiles at
// the top (labelled "recent"), bypassing the group-based two-step selector.
func selectProfileForConsole(profiles []awsconfig.Profile) (awsconfig.Profile, error) {
	recent := history.Load("console")
	recentSet := make(map[string]bool, len(recent))
	for _, r := range recent {
		recentSet[r] = true
	}

	var items []ui.Item
	for _, name := range recent {
		items = append(items, ui.Item{Title: name, Desc: "recent"})
	}
	for _, p := range profiles {
		if !recentSet[p.Name] {
			desc := p.AccountID
			if p.RoleName != "" {
				desc += " / " + p.RoleName
			}
			items = append(items, ui.Item{Title: p.Name, Desc: desc})
		}
	}

	selectedName, err := ui.Choose("Select AWS profile", items)
	if err != nil {
		return awsconfig.Profile{}, err
	}
	if selectedName == "" {
		return awsconfig.Profile{}, nil
	}
	p, _ := awsconfig.FindProfile(profiles, selectedName)
	return p, nil
}
