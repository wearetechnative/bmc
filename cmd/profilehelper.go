package cmd

import (
	"fmt"
	"os"

	"github.com/wearetechnative/bmc/internal/awsconfig"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/mfa"
)

// ensureAWSProfile returns the active AWS_PROFILE, prompting if not set.
// It also runs the MFA check.
func ensureAWSProfile() (string, error) {
	profile := os.Getenv("AWS_PROFILE")
	if profile != "" {
		return profile, nil
	}

	profiles, err := awsconfig.LoadProfiles()
	if err != nil {
		return "", err
	}

	selected, err := selectProfileInteractive(profiles)
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

	os.Setenv("AWS_PROFILE", selected.Name)
	return selected.Name, nil
}
