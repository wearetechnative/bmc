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
	consolePick    bool
	consoleService string
	consoleWatch   bool
)

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Open Firefox with AWS account in console",
	Args:  cobra.NoArgs,
	RunE:  runConsole,
}

func init() {
	consoleCmd.Flags().StringVarP(&consoleProfile, "profile", "p", "", "AWS profile name to use")
	consoleCmd.Flags().BoolVarP(&consolePick, "pick", "P", false, "Force interactive profile selection (ignores AWS_PROFILE)")
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

	profileName := strings.TrimSpace(consoleProfile)

	interactive := false

	switch {
	case consolePick:
		// -P/--pick: force interactive selection (ignore AWS_PROFILE)
		selectedProfile, interactive, err = selectProfileWithHistory(profiles)
		if err != nil {
			return err
		}
		if selectedProfile.Name == "" {
			return nil
		}
	case profileName != "":
		// -p <name>: use the given profile directly
		p, ok := awsconfig.FindProfile(profiles, profileName)
		if !ok {
			return fmt.Errorf("profile %q not found", profileName)
		}
		selectedProfile = p
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
	expiry, err := awsops.OpenConsole(selectedProfile.Name, consoleService, cfg)
	if err != nil {
		return err
	}
	if !expiry.IsZero() {
		fmt.Fprintf(os.Stderr, "-- Console session valid until: %s\n", expiry.Local().Format("15:04:05"))
	}

	if interactive {
		_ = history.Save("profile", selectedProfile.Name)
	}

	if consoleWatch {
		registerConsoleSession(selectedProfile.Name, consoleService, expiry)
	}

	return nil
}

// registerConsoleSession registers the opened session with the watcher daemon,
// starting the daemon first if it is not already running. expiry is the real
// credential expiry reported by OpenConsole; guessing it here would make the
// watcher refresh long after the session had already died.
func registerConsoleSession(profile, service string, expiry time.Time) {
	if expiry.IsZero() {
		// Static credentials never expire; assume the shortest session AWS
		// would hand out so the watcher keeps checking.
		expiry = time.Now().Add(awsops.FallbackSessionDuration)
	}
	s := watcher.Session{
		Profile:       profile,
		Service:       service,
		ContainerName: profile,
		Expiry:        expiry,
		RefreshAt:     watcher.RefreshTime(expiry),
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
