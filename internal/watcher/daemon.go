package watcher

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/wearetechnative/bmc/internal/awsops"
	"github.com/wearetechnative/bmc/internal/config"
)

const (
	pollInterval  = 30 * time.Second
	refreshWindow = 5 * time.Minute
	// startupGrace is how long the daemon waits for the first session to be
	// registered before self-terminating. This allows `bmc watcher start` to
	// keep the daemon alive while the user opens a console with --watch.
	startupGrace = 60 * time.Second
)

// RunDaemon is the daemon entry point, invoked when BMC_WATCHER_DAEMON=1.
// It starts the HTTP server, registers itself in the state file, and runs
// the poll loop until no sessions remain.
func RunDaemon() {
	bmcCfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "watcher: failed to load config: %v\n", err)
		os.Exit(1)
	}

	srv, err := StartServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "watcher: failed to start HTTP server: %v\n", err)
		os.Exit(1)
	}

	// Read sessions written by the parent before forking, then update PID/port.
	state, err := ReadState()
	if err != nil {
		state = WatcherState{}
	}
	state.PID = os.Getpid()
	state.StartedAt = time.Now().UTC()
	state.Port = srv.Port()
	if err := WriteState(state); err != nil {
		fmt.Fprintf(os.Stderr, "watcher: failed to write state: %v\n", err)
	}

	// Handle SIGTERM and SIGINT: clean up and exit.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigCh
		_ = ClearState()
		os.Exit(0)
	}()

	startTime := time.Now()
	hadSessions := false

	for {
		runPollLoop(srv, bmcCfg)

		st, err := ReadState()
		if err == nil && len(st.Sessions) > 0 {
			hadSessions = true
		}

		// During the startup grace period, keep running even without sessions
		// so that `bmc watcher start` stays alive long enough for the user to
		// open a console with --watch.
		withinGrace := time.Since(startTime) < startupGrace
		if !withinGrace || hadSessions {
			// Self-terminate when no sessions with a future expiry remain.
			if err != nil || len(st.Sessions) == 0 {
				_ = ClearState()
				os.Exit(0)
			}
			hasActive := false
			for _, s := range st.Sessions {
				if s.Expiry.After(time.Now()) {
					hasActive = true
					break
				}
			}
			if !hasActive {
				_ = ClearState()
				os.Exit(0)
			}
		}

		time.Sleep(pollInterval)
	}
}

// Fork starts a detached watcher daemon by re-executing the bmc binary with
// BMC_WATCHER_DAEMON=1 and Setsid set so it survives terminal close.
// Returns the PID of the spawned process.
func Fork() (int, error) {
	exe, err := os.Executable()
	if err != nil {
		return 0, fmt.Errorf("cannot find bmc executable: %w", err)
	}
	cmd := exec.Command(exe, "watcher", "start")
	cmd.Env = append(os.Environ(), "BMC_WATCHER_DAEMON=1")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("failed to fork daemon: %w", err)
	}
	return cmd.Process.Pid, nil
}

func runPollLoop(srv *Server, bmcCfg config.Config) {
	state, err := ReadState()
	if err != nil || len(state.Sessions) == 0 {
		return
	}

	now := time.Now()
	updated := false
	for i := range state.Sessions {
		s := &state.Sessions[i]
		if s.Expiry.Before(now) {
			continue // already expired
		}
		if s.RefreshAt.After(now) {
			continue // not yet time to refresh
		}
		newSession, err := refreshSession(*s, srv, bmcCfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "watcher: refresh failed for %s: %v\n", s.Profile, err)
			continue
		}
		state.Sessions[i] = newSession
		updated = true
	}
	if updated {
		_ = WriteState(state)
	}
}

func refreshSession(s Session, srv *Server, bmcCfg config.Config) (Session, error) {
	signinURL, credExpiry, err := awsops.BuildFederationURL(s.Profile, s.Service, bmcCfg)
	if err != nil {
		return s, fmt.Errorf("failed to build federation URL: %w", err)
	}

	localURL := srv.RefreshURL(signinURL)
	if err := awsops.OpenURLInBrowser(localURL, s.ContainerName, bmcCfg.Console); err != nil {
		// Fallback: open the federation URL directly in the container.
		fmt.Fprintf(os.Stderr, "watcher: local refresh page failed for %s, falling back to direct URL\n", s.Profile)
		if err2 := awsops.OpenURLInBrowser(signinURL, s.ContainerName, bmcCfg.Console); err2 != nil {
			return s, fmt.Errorf("fallback refresh also failed: %w", err2)
		}
	}

	s.Expiry = credExpiry
	s.RefreshAt = credExpiry.Add(-refreshWindow)
	return s, nil
}
