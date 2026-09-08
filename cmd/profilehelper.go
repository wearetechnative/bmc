package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsconfig"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/mfa"
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
		selected, err := selectProfileInteractive(profiles)
		if err != nil {
			return "", err
		}
		if selected.Name == "" {
			return "", fmt.Errorf("no profile selected")
		}
		if err := runMFACheck(selected); err != nil {
			return "", err
		}
		os.Setenv("AWS_PROFILE", selected.Name)
		return selected.Name, nil
	}
}
