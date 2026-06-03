package watcher

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// WatcherState is the content of ~/.config/bmc/watcher.json.
type WatcherState struct {
	PID       int       `json:"pid"`
	StartedAt time.Time `json:"started_at,omitempty"`
	Port      int       `json:"port,omitempty"`
	Sessions  []Session `json:"sessions,omitempty"`
}

// Session represents one active console session being kept alive by the watcher.
type Session struct {
	Profile       string    `json:"profile"`
	Service       string    `json:"service"`
	ContainerName string    `json:"container_name"`
	Expiry        time.Time `json:"expiry"`
	RefreshAt     time.Time `json:"refresh_at"`
}

// StatePath returns the path to ~/.config/bmc/watcher.json.
func StatePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "bmc", "watcher.json")
}

// ReadState reads the watcher state file. Returns an empty state if the file
// does not exist.
func ReadState() (WatcherState, error) {
	data, err := os.ReadFile(StatePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return WatcherState{}, nil
		}
		return WatcherState{}, err
	}
	var state WatcherState
	if err := json.Unmarshal(data, &state); err != nil {
		return WatcherState{}, err
	}
	return state, nil
}

// WriteState writes the watcher state to ~/.config/bmc/watcher.json,
// creating the directory if needed.
func WriteState(state WatcherState) error {
	path := StatePath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// ClearState removes the watcher state file.
func ClearState() error {
	err := os.Remove(StatePath())
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// IsAlive returns true if a process with the given PID is running.
func IsAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}

// EnsureWatcher checks whether a watcher daemon is already running.
// If the state file contains a dead PID, it clears the state file first.
// Returns true if a daemon is already running.
func EnsureWatcher() (bool, error) {
	state, err := ReadState()
	if err != nil {
		return false, err
	}
	if state.PID != 0 && IsAlive(state.PID) {
		return true, nil
	}
	// No living daemon — clear any stale state.
	if err := ClearState(); err != nil {
		return false, err
	}
	return false, nil
}

// RegisterSession appends a session to the watcher state file.
// Safe to call regardless of whether the daemon is running; the daemon
// picks up new sessions on its next poll cycle.
func RegisterSession(s Session) error {
	state, err := ReadState()
	if err != nil {
		return err
	}
	state.Sessions = append(state.Sessions, s)
	return WriteState(state)
}
