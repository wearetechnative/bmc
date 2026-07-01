package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/wearetechnative/bmc/internal/awsconfig"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/history"
	"github.com/wearetechnative/bmc/internal/mfa"
	"github.com/wearetechnative/bmc/internal/ui"
)

var globalProfile string

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
// used groups and profiles surfaced at the top. Returns the selected profile,
// true if a profile was interactively selected (false if cancelled), and any error.
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

// ensureAWSProfile returns the active AWS_PROFILE, prompting if not set.
// It also runs the MFA check.
func ensureAWSProfile() (string, error) {
	if globalProfile != "" {
		profiles, err := awsconfig.LoadProfiles()
		if err != nil {
			return "", err
		}

		// bare -p (no value): force interactive picker, ignore AWS_PROFILE
		if strings.TrimSpace(globalProfile) == "" {
			selected, _, err := selectProfileWithHistory(profiles)
			if err != nil {
				return "", err
			}
			if selected.Name == "" {
				return "", fmt.Errorf("no profile selected")
			}
			sourceProfile, err := awsconfig.ResolveSourceProfile(selected)
			if err != nil {
				return "", err
			}
			cfg, err := config.Load()
			if err != nil {
				return "", err
			}
			if err := mfa.EnsureValid(sourceProfile, cfg, os.Stderr); err != nil {
				return "", err
			}
			_ = history.Save("profile", selected.Name)
			os.Setenv("AWS_PROFILE", selected.Name)
			return selected.Name, nil
		}

		matched, ok := awsconfig.FindProfile(profiles, globalProfile)
		if !ok {
			return "", fmt.Errorf("profile %q not found", globalProfile)
		}
		sourceProfile, err := awsconfig.ResolveSourceProfile(matched)
		if err != nil {
			return "", err
		}
		cfg, err := config.Load()
		if err != nil {
			return "", err
		}
		if err := mfa.EnsureValid(sourceProfile, cfg, os.Stderr); err != nil {
			return "", err
		}
		return strings.TrimSpace(globalProfile), nil
	}

	profile := os.Getenv("AWS_PROFILE")
	if profile != "" {
		profiles, err := awsconfig.LoadProfiles()
		if err != nil {
			return "", err
		}

		var matched *awsconfig.Profile
		for i := range profiles {
			if profiles[i].Name == profile {
				matched = &profiles[i]
				break
			}
		}

		if matched == nil {
			fmt.Fprintf(os.Stderr, "warning: AWS_PROFILE=%q not found in ~/.aws/config, skipping MFA check\n", profile)
			return profile, nil
		}

		sourceProfile, err := awsconfig.ResolveSourceProfile(*matched)
		if err != nil {
			return "", err
		}

		cfg, err := config.Load()
		if err != nil {
			return "", err
		}

		if err := mfa.EnsureValid(sourceProfile, cfg, os.Stderr); err != nil {
			return "", err
		}

		return profile, nil
	}

	profiles, err := awsconfig.LoadProfiles()
	if err != nil {
		return "", err
	}

	selected, _, err := selectProfileWithHistory(profiles)
	if err != nil {
		return "", err
	}
	if selected.Name == "" {
		return "", fmt.Errorf("no profile selected")
	}

	sourceProfile, err := awsconfig.ResolveSourceProfile(selected)
	if err != nil {
		return "", err
	}

	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	if err := mfa.EnsureValid(sourceProfile, cfg, os.Stderr); err != nil {
		return "", err
	}

	_ = history.Save("profile", selected.Name)
	os.Setenv("AWS_PROFILE", selected.Name)
	return selected.Name, nil
}
