package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsconfig"
	"github.com/wearetechnative/bmc/internal/awsops"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/history"
	"github.com/wearetechnative/bmc/internal/mfa"
	"github.com/wearetechnative/bmc/internal/watcher"
)

var (
	consoleProfile string
	consoleService string
	consoleWatch   bool
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
	consoleCmd.Flags().BoolVarP(&consoleWatch, "watch", "w", false, "Keep session alive via background watcher")
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
		selectedProfile, interactive, err = selectProfileWithHistory(profiles)
		if err != nil {
			return err
		}
		if selectedProfile.Name == "" {
			return nil
		}
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
			selectedProfile, interactive, err = selectProfileWithHistory(profiles)
			if err != nil {
				return err
			}
			if selectedProfile.Name == "" {
				return nil
			}
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
		_ = history.Save("profile", selectedProfile.Name)
	}

	if consoleWatch {
		registerConsoleSession(selectedProfile.Name, consoleService)
	}

	return nil
}

// registerConsoleSession registers the opened session with the watcher daemon,
// starting the daemon first if it is not already running.
func registerConsoleSession(profile, service string) {
	// AWS federation sessions are capped at 1 hour for role-assumed credentials.
	expiry := time.Now().Add(time.Hour)
	s := watcher.Session{
		Profile:       profile,
		Service:       service,
		ContainerName: profile,
		Expiry:        expiry,
		RefreshAt:     expiry.Add(-5 * time.Minute),
	}

	alreadyRunning, err := watcher.EnsureWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "-- watcher: %v\n", err)
		return
	}

	if err := watcher.RegisterSession(s); err != nil {
		fmt.Fprintf(os.Stderr, "-- watcher: failed to register session: %v\n", err)
		return
	}

	if !alreadyRunning {
		pid, err := watcher.Fork()
		if err != nil {
			fmt.Fprintf(os.Stderr, "-- watcher: failed to start daemon: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stderr, "-- watcher started (PID %d)\n", pid)
	} else {
		state, _ := watcher.ReadState()
		fmt.Fprintf(os.Stderr, "-- watcher: session registered (PID %d)\n", state.PID)
	}
}
