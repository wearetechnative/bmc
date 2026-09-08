package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsconfig"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/history"
	"github.com/wearetechnative/bmc/internal/mfa"
	"github.com/wearetechnative/bmc/internal/ui"
)

var (
	globalProfile string
	globalPick    bool
)

// addProfileFlags registers the standard profile-selection flags on cmd:
//
//	-p / --profile NAME  always takes a value (space and equals forms both work)
//	-P / --pick          forces the interactive picker (ignores AWS_PROFILE)
//
// Both flags bind to the shared globalProfile/globalPick variables, so every
// command that calls this gets identical, predictable parsing behavior.
func addProfileFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&globalProfile, "profile", "p", "", "AWS profile name to use")
	cmd.Flags().BoolVarP(&globalPick, "pick", "P", false, "Force interactive profile selection (ignores AWS_PROFILE)")
}

// profileMode is the resolved intent of the profile-selection flags.
type profileMode int

const (
	// modeInteractive: show the interactive profile picker.
	modeInteractive profileMode = iota
	// modeNamed: use the explicitly named profile.
	modeNamed
	// modeEnv: use the AWS_PROFILE environment variable.
	modeEnv
)

// resolveProfileMode maps the flag/env state to a selection mode.
//
// Precedence: --pick forces the picker; otherwise a named --profile wins;
// otherwise AWS_PROFILE is used; otherwise fall back to the picker.
func resolveProfileMode(profile string, pick bool, envProfile string) (profileMode, string) {
	if pick {
		return modeInteractive, ""
	}
	if p := strings.TrimSpace(profile); p != "" {
		return modeNamed, p
	}
	if strings.TrimSpace(envProfile) != "" {
		return modeEnv, envProfile
	}
	return modeInteractive, ""
}

// runMFACheck resolves the source profile and validates the MFA session for p.
func runMFACheck(p awsconfig.Profile) error {
	sourceProfile, err := awsconfig.ResolveSourceProfile(p)
	if err != nil {
		return err
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	return mfa.EnsureValid(sourceProfile, cfg, os.Stderr)
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

// selectProfileWithHistory shows a two-step group-aware selector with recently
// used groups and profiles surfaced at the top (marked "recent"). History is
// shared across all commands via the "profile" store. Returns the selected
// profile, true if a profile was interactively selected (false if cancelled),
// and any error.
func selectProfileWithHistory(profiles []awsconfig.Profile) (awsconfig.Profile, bool, error) {
	recent := history.Load("profile")
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
		return awsconfig.Profile{}, false, fmt.Errorf("no profile groups found in ~/.aws/config")
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
			return awsconfig.Profile{}, false, err
		}
		if selectedGroup == "" {
			return awsconfig.Profile{}, false, nil
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
			return awsconfig.Profile{}, false, err
		}
		if selectedName == "" {
			return awsconfig.Profile{}, false, nil
		}

		p, _ := awsconfig.FindProfile(profiles, selectedName)
		return p, true, nil
	}
}

// ensureAWSProfile returns the active AWS_PROFILE, prompting if needed.
// It also runs the MFA check for the resolved profile.
func ensureAWSProfile() (string, error) {
	profiles, err := awsconfig.LoadProfiles()
	if err != nil {
		return "", err
	}

	envProfile := os.Getenv("AWS_PROFILE")
	mode, name := resolveProfileMode(globalProfile, globalPick, envProfile)

	switch mode {
	case modeNamed:
		matched, ok := awsconfig.FindProfile(profiles, name)
		if !ok {
			return "", fmt.Errorf("profile %q not found", name)
		}
		if err := runMFACheck(matched); err != nil {
			return "", err
		}
		return name, nil

	case modeEnv:
		matched, ok := awsconfig.FindProfile(profiles, name)
		if !ok {
			fmt.Fprintf(os.Stderr, "warning: AWS_PROFILE=%q not found in ~/.aws/config, skipping MFA check\n", name)
			return name, nil
		}
		if err := runMFACheck(matched); err != nil {
			return "", err
		}
		return name, nil

	default: // modeInteractive
		selected, _, err := selectProfileWithHistory(profiles)
		if err != nil {
			return "", err
		}
		if selected.Name == "" {
			return "", fmt.Errorf("no profile selected")
		}
		if err := runMFACheck(selected); err != nil {
			return "", err
		}
		_ = history.Save("profile", selected.Name)
		os.Setenv("AWS_PROFILE", selected.Name)
		return selected.Name, nil
	}
}
