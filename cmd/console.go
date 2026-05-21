package cmd

import (
	"errors"
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
		selectedProfile, err = selectProfileForConsoleInteractive(profiles)
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
			selectedProfile, err = selectProfileForConsoleInteractive(profiles)
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

// recentGroups returns the groups of recently-used profiles, in order of most
// recent use, with duplicates removed. Groups whose profiles no longer exist in
// the config are silently skipped.
func recentGroups(profiles []awsconfig.Profile, recentProfiles []string) []string {
	profileGroup := make(map[string]string, len(profiles))
	for _, p := range profiles {
		profileGroup[p.Name] = p.Group
	}
	seen := make(map[string]bool)
	var groups []string
	for _, name := range recentProfiles {
		g := profileGroup[name]
		if g != "" && !seen[g] {
			seen[g] = true
			groups = append(groups, g)
		}
	}
	return groups
}

// selectProfileForConsoleInteractive shows a two-step group-aware selector.
// Step 1: account groups, with recently-used groups surfaced at the top.
// Step 2: profiles within the selected group, with recently-used profiles at the top.
// Pressing back in step 2 returns to step 1.
func selectProfileForConsoleInteractive(profiles []awsconfig.Profile) (awsconfig.Profile, error) {
	recent := history.Load("console")
	recentSet := make(map[string]bool, len(recent))
	for _, r := range recent {
		recentSet[r] = true
	}

	recentGroupList := recentGroups(profiles, recent)
	recentGroupSet := make(map[string]bool, len(recentGroupList))
	for _, g := range recentGroupList {
		recentGroupSet[g] = true
	}

	allGroups := awsconfig.Groups(profiles)
	if len(allGroups) == 0 {
		return awsconfig.Profile{}, fmt.Errorf("no profile groups found in ~/.aws/config")
	}

	var groupItems []ui.Item
	for _, g := range recentGroupList {
		groupItems = append(groupItems, ui.Item{Title: g, Desc: "recent"})
	}
	for _, g := range allGroups {
		if !recentGroupSet[g] {
			groupItems = append(groupItems, ui.Item{Title: g})
		}
	}

	for {
		selectedGroup, err := ui.Choose("Select AWS account group", groupItems)
		if err != nil {
			return awsconfig.Profile{}, err
		}
		if selectedGroup == "" {
			return awsconfig.Profile{}, nil
		}

		groupProfiles := awsconfig.ByGroup(profiles, selectedGroup)
		var profileItems []ui.Item
		for _, p := range groupProfiles {
			if recentSet[p.Name] {
				profileItems = append(profileItems, ui.Item{Title: p.Name, Desc: "recent"})
			}
		}
		for _, p := range groupProfiles {
			if !recentSet[p.Name] {
				desc := p.AccountID
				if p.RoleName != "" {
					desc += " / " + p.RoleName
				}
				profileItems = append(profileItems, ui.Item{Title: p.Name, Desc: desc})
			}
		}

		selectedName, err := ui.Choose("Select profile  (group: "+selectedGroup+")", profileItems)
		if errors.Is(err, ui.ErrBack) {
			continue
		}
		if err != nil {
			return awsconfig.Profile{}, err
		}
		if selectedName == "" {
			return awsconfig.Profile{}, nil
		}

		p, _ := awsconfig.FindProfile(profiles, selectedName)
		return p, nil
	}
}
