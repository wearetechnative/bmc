package cmd

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/watcher"
)

var watcherCmd = &cobra.Command{
	Use:   "watcher",
	Short: "Manage the console session keep-alive daemon",
}

var watcherStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the session keep-alive daemon",
	RunE:  runWatcherStart,
}

var watcherStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the session keep-alive daemon",
	RunE:  runWatcherStop,
}

var watcherStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show active console sessions being kept alive",
	RunE:  runWatcherStatus,
}

func init() {
	watcherCmd.AddCommand(watcherStartCmd)
	watcherCmd.AddCommand(watcherStopCmd)
	watcherCmd.AddCommand(watcherStatusCmd)
	rootCmd.AddCommand(watcherCmd)
}

func runWatcherStart(_ *cobra.Command, _ []string) error {
	// Daemon path: BMC_WATCHER_DAEMON=1 is set by the parent fork.
	if os.Getenv("BMC_WATCHER_DAEMON") == "1" {
		watcher.RunDaemon()
		return nil
	}
	return ensureAndStartWatcher()
}

// ensureAndStartWatcher checks if the daemon is running and starts it if not.
// Shared by runWatcherStart and the console --watch flag.
func ensureAndStartWatcher() error {
	alreadyRunning, err := watcher.EnsureWatcher()
	if err != nil {
		return fmt.Errorf("watcher: %w", err)
	}
	if alreadyRunning {
		state, _ := watcher.ReadState()
		fmt.Fprintf(os.Stderr, "-- watcher already running (PID %d)\n", state.PID)
		return nil
	}
	pid, err := watcher.Fork()
	if err != nil {
		return fmt.Errorf("watcher: %w", err)
	}
	fmt.Fprintf(os.Stderr, "-- watcher started (PID %d)\n", pid)
	return nil
}

func runWatcherStop(_ *cobra.Command, _ []string) error {
	state, err := watcher.ReadState()
	if err != nil {
		return fmt.Errorf("watcher: %w", err)
	}
	if state.PID == 0 || !watcher.IsAlive(state.PID) {
		fmt.Fprintln(os.Stderr, "watcher is not running")
		return nil
	}
	proc, err := os.FindProcess(state.PID)
	if err != nil {
		return fmt.Errorf("watcher: cannot find process %d: %w", state.PID, err)
	}
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("watcher: failed to stop daemon: %w", err)
	}
	_ = watcher.ClearState()
	fmt.Fprintln(os.Stderr, "-- watcher stopped")
	return nil
}

func runWatcherStatus(_ *cobra.Command, _ []string) error {
	state, err := watcher.ReadState()
	if err != nil {
		return fmt.Errorf("watcher: %w", err)
	}
	if state.PID == 0 || !watcher.IsAlive(state.PID) {
		fmt.Fprintln(os.Stderr, "watcher is not running")
		return nil
	}
	fmt.Fprintf(os.Stderr, "-- watcher running (PID %d)\n", state.PID)
	if len(state.Sessions) == 0 {
		fmt.Fprintln(os.Stderr, "   no active sessions")
		return nil
	}
	now := time.Now()
	for _, s := range state.Sessions {
		until := s.RefreshAt.Sub(now)
		if until < 0 {
			until = 0
		}
		service := s.Service
		if service == "" {
			service = "console"
		}
		fmt.Fprintf(os.Stderr, "   %-30s  %-20s  refreshes in %s\n",
			s.Profile, s.ContainerName+"/"+service, formatWatcherDuration(until))
	}
	return nil
}

func formatWatcherDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh%dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}
