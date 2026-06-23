package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFile is a test helper that writes content to a file, creating
// intermediate directories as needed.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

// --- FindDefaultProfile ---

func TestFindDefaultProfile_ReturnsAbsolutePath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	iniContent := `
[Profile0]
Name=default
IsRelative=1
Path=Profiles/abc123.default
Default=1
`
	writeFile(t, filepath.Join(home, ".mozilla", "firefox", "profiles.ini"), iniContent)

	got, err := FindDefaultProfile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(home, ".mozilla", "firefox", "Profiles", "abc123.default")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFindDefaultProfile_AbsolutePath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	absProfile := filepath.Join(home, "myprofile")
	iniContent := `
[Profile0]
Name=custom
IsRelative=0
Path=` + absProfile + `
Default=1
`
	writeFile(t, filepath.Join(home, ".mozilla", "firefox", "profiles.ini"), iniContent)

	got, err := FindDefaultProfile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != absProfile {
		t.Errorf("got %q, want %q", got, absProfile)
	}
}

func TestFindDefaultProfile_NoDefaultProfile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	iniContent := `
[Profile0]
Name=other
IsRelative=1
Path=Profiles/xyz.other
`
	writeFile(t, filepath.Join(home, ".mozilla", "firefox", "profiles.ini"), iniContent)

	_, err := FindDefaultProfile()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFindDefaultProfile_MissingFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	_, err := FindDefaultProfile()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- IsDebugPortConfigured ---

func TestIsDebugPortConfigured_MarkerPresent(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "user.js"), "// bmc watcher: start Firefox with --remote-debugging-port\n")

	if !IsDebugPortConfigured(dir) {
		t.Error("expected true, got false")
	}
}

func TestIsDebugPortConfigured_MarkerWithConflict(t *testing.T) {
	dir := t.TempDir()
	content := "// bmc watcher: start Firefox with --remote-debugging-port\n" +
		`user_pref("remote.enabled", true);` + "\n"
	writeFile(t, filepath.Join(dir, "user.js"), content)

	if IsDebugPortConfigured(dir) {
		t.Error("expected false (conflicting pref), got true")
	}
}

func TestIsDebugPortConfigured_NotPresent(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "user.js"), `user_pref("browser.startup.homepage", "about:newtab");`)

	if IsDebugPortConfigured(dir) {
		t.Error("expected false, got true")
	}
}

func TestHasConflictingPrefs_Detected(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "user.js"), `user_pref("remote.enabled", true);`)

	if !HasConflictingPrefs(dir) {
		t.Error("expected true, got false")
	}
}

func TestHasConflictingPrefs_Clean(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "user.js"), `user_pref("browser.startup.homepage", "about:newtab");`)

	if HasConflictingPrefs(dir) {
		t.Error("expected false, got true")
	}
}

// --- WriteDebugPortConfig ---

func TestWriteDebugPortConfig_WritesMarker(t *testing.T) {
	dir := t.TempDir()

	if err := WriteDebugPortConfig(dir, 9222); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// After writing, IsDebugPortConfigured should return true.
	if !IsDebugPortConfigured(dir) {
		t.Error("expected IsDebugPortConfigured=true after WriteDebugPortConfig")
	}
}

func TestWriteDebugPortConfig_AppendsToExistingFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "user.js"), `user_pref("browser.startup.homepage", "about:newtab");`)

	if err := WriteDebugPortConfig(dir, 9222); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "user.js"))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	content := string(data)
	// Original pref must still be there.
	if !contains(content, "browser.startup.homepage") {
		t.Error("original content was lost")
	}
	// Marker must be appended.
	if !contains(content, "bmc watcher") {
		t.Error("bmc watcher marker was not appended")
	}
}

func TestWriteDebugPortConfig_NoConflictingPrefs(t *testing.T) {
	dir := t.TempDir()

	if err := WriteDebugPortConfig(dir, 9222); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if HasConflictingPrefs(dir) {
		t.Error("WriteDebugPortConfig should not write conflicting prefs")
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
