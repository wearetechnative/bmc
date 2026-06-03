package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// overrideStatePath temporarily redirects the state file to a temp location.
func overrideStatePath(t *testing.T) {
	t.Helper()
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	_ = filepath.Join(tmpHome, ".config", "bmc", "watcher.json")
}

func TestIsAlive_OwnPID(t *testing.T) {
	if !IsAlive(os.Getpid()) {
		t.Fatal("IsAlive should return true for own PID")
	}
}

func TestIsAlive_ZeroPID(t *testing.T) {
	if IsAlive(0) {
		t.Fatal("IsAlive should return false for PID 0")
	}
}

func TestIsAlive_NegativePID(t *testing.T) {
	if IsAlive(-1) {
		t.Fatal("IsAlive should return false for negative PID")
	}
}

func TestReadState_MissingFile(t *testing.T) {
	overrideStatePath(t)
	state, err := ReadState()
	if err != nil {
		t.Fatalf("ReadState with missing file: unexpected error: %v", err)
	}
	if state.PID != 0 {
		t.Fatalf("expected empty state, got PID=%d", state.PID)
	}
}

func TestWriteAndReadState_RoundTrip(t *testing.T) {
	overrideStatePath(t)

	now := time.Now().UTC().Truncate(time.Second)
	want := WatcherState{
		PID:       42,
		StartedAt: now,
		Port:      12345,
		Sessions: []Session{
			{
				Profile:       "test-profile",
				Service:       "ec2",
				ContainerName: "test-profile",
				Expiry:        now.Add(time.Hour),
				RefreshAt:     now.Add(55 * time.Minute),
			},
		},
	}

	if err := WriteState(want); err != nil {
		t.Fatalf("WriteState: %v", err)
	}

	got, err := ReadState()
	if err != nil {
		t.Fatalf("ReadState: %v", err)
	}

	if got.PID != want.PID {
		t.Errorf("PID: got %d, want %d", got.PID, want.PID)
	}
	if got.Port != want.Port {
		t.Errorf("Port: got %d, want %d", got.Port, want.Port)
	}
	if len(got.Sessions) != 1 {
		t.Fatalf("Sessions: got %d, want 1", len(got.Sessions))
	}
	if got.Sessions[0].Profile != "test-profile" {
		t.Errorf("Profile: got %q, want %q", got.Sessions[0].Profile, "test-profile")
	}
}

func TestClearState(t *testing.T) {
	overrideStatePath(t)

	if err := WriteState(WatcherState{PID: 1}); err != nil {
		t.Fatal(err)
	}
	if err := ClearState(); err != nil {
		t.Fatalf("ClearState: %v", err)
	}
	// Second clear should be a no-op.
	if err := ClearState(); err != nil {
		t.Fatalf("ClearState (second): %v", err)
	}
	state, err := ReadState()
	if err != nil {
		t.Fatal(err)
	}
	if state.PID != 0 {
		t.Errorf("expected empty state after clear, got PID=%d", state.PID)
	}
}

func TestRegisterSession(t *testing.T) {
	overrideStatePath(t)

	now := time.Now().UTC()
	s1 := Session{Profile: "account-a", Expiry: now.Add(time.Hour)}
	s2 := Session{Profile: "account-b", Expiry: now.Add(time.Hour)}

	if err := RegisterSession(s1); err != nil {
		t.Fatalf("RegisterSession s1: %v", err)
	}
	if err := RegisterSession(s2); err != nil {
		t.Fatalf("RegisterSession s2: %v", err)
	}

	state, err := ReadState()
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(state.Sessions))
	}
}

func TestEnsureWatcher_NoStateFile(t *testing.T) {
	overrideStatePath(t)
	running, err := EnsureWatcher()
	if err != nil {
		t.Fatalf("EnsureWatcher: %v", err)
	}
	if running {
		t.Fatal("expected not running when no state file exists")
	}
}

func TestEnsureWatcher_DeadPID(t *testing.T) {
	overrideStatePath(t)
	// PID 2147483647 is almost certainly not a running process.
	if err := WriteState(WatcherState{PID: 2147483647}); err != nil {
		t.Fatal(err)
	}
	running, err := EnsureWatcher()
	if err != nil {
		t.Fatalf("EnsureWatcher: %v", err)
	}
	if running {
		t.Fatal("expected not running for dead PID")
	}
	// State file should be cleared.
	state, err := ReadState()
	if err != nil {
		t.Fatal(err)
	}
	if state.PID != 0 {
		t.Errorf("expected state file cleared, got PID=%d", state.PID)
	}
}

func TestPollLoop_RefreshesExpiredSession(t *testing.T) {
	overrideStatePath(t)

	now := time.Now()
	// Session whose RefreshAt is in the past (should be refreshed).
	past := Session{
		Profile:   "stale",
		Service:   "ec2",
		Expiry:    now.Add(10 * time.Minute),
		RefreshAt: now.Add(-1 * time.Minute), // overdue
	}
	// Session whose RefreshAt is in the future (should not be refreshed yet).
	future := Session{
		Profile:   "fresh",
		Service:   "ec2",
		Expiry:    now.Add(time.Hour),
		RefreshAt: now.Add(55 * time.Minute),
	}

	if err := WriteState(WatcherState{Sessions: []Session{past, future}}); err != nil {
		t.Fatal(err)
	}

	state, _ := ReadState()
	for _, s := range state.Sessions {
		if s.Profile == "stale" && s.RefreshAt.After(now) {
			t.Error("stale session should have RefreshAt in the past")
		}
		if s.Profile == "fresh" && !s.RefreshAt.After(now) {
			t.Error("fresh session should have RefreshAt in the future")
		}
	}
}
