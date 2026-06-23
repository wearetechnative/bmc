package watcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

// FindDefaultProfile returns the absolute path to the default Firefox profile
// directory by reading ~/.mozilla/firefox/profiles.ini.
func FindDefaultProfile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	iniPath := filepath.Join(home, ".mozilla", "firefox", "profiles.ini")

	data, err := os.ReadFile(iniPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("Firefox profiles.ini not found at %s", iniPath)
		}
		return "", fmt.Errorf("cannot read %s: %w", iniPath, err)
	}

	cfg, err := ini.Load(data)
	if err != nil {
		return "", fmt.Errorf("cannot parse profiles.ini: %w", err)
	}

	for _, section := range cfg.Sections() {
		if !strings.HasPrefix(section.Name(), "Profile") {
			continue
		}
		if section.Key("Default").Value() != "1" {
			continue
		}
		path := section.Key("Path").Value()
		if path == "" {
			continue
		}
		isRelative := section.Key("IsRelative").Value()
		if isRelative == "1" {
			path = filepath.Join(home, ".mozilla", "firefox", path)
		}
		return path, nil
	}
	return "", fmt.Errorf("no default Firefox profile found in %s", iniPath)
}

// bmcWatcherMarker is written to user.js to mark that bmc watcher setup was run.
// No browser preferences are needed — Firefox BiDi is enabled solely via
// --remote-debugging-port at startup.
const bmcWatcherMarker = "// bmc watcher: start Firefox with --remote-debugging-port"

// conflictingPrefs are user.js entries written by older bmc versions that
// cause Firefox to start a bare httpd.js server that blocks BiDi registration.
var conflictingPrefs = []string{
	`remote.enabled`,
	`remote.force-local`,
	`devtools.debugger.remote-enabled`,
	`devtools.debugger.remote-port`,
	`devtools.debugger.remote-host`,
	`devtools.debugger.prompt-connection`,
}

// IsDebugPortConfigured returns true when the bmc watcher marker is present
// in user.js and no conflicting prefs exist.
func IsDebugPortConfigured(profileDir string) bool {
	data, err := os.ReadFile(filepath.Join(profileDir, "user.js"))
	if err != nil {
		return false
	}
	content := string(data)
	if !strings.Contains(content, bmcWatcherMarker) {
		return false
	}
	for _, p := range conflictingPrefs {
		if strings.Contains(content, p) {
			return false // marker present but conflicting prefs still exist
		}
	}
	return true
}

// HasConflictingPrefs returns true if user.js contains prefs that interfere
// with Firefox BiDi (written by older bmc versions).
func HasConflictingPrefs(profileDir string) bool {
	data, err := os.ReadFile(filepath.Join(profileDir, "user.js"))
	if err != nil {
		return false
	}
	content := string(data)
	for _, p := range conflictingPrefs {
		if strings.Contains(content, p) {
			return true
		}
	}
	return false
}

// WriteDebugPortConfig writes only a marker comment to user.js.
// Firefox BiDi needs no preferences — it is activated solely via
// --remote-debugging-port=<port> at startup.
func WriteDebugPortConfig(profileDir string, _ int) error {
	userJS := filepath.Join(profileDir, "user.js")
	marker := "\n" + bmcWatcherMarker + "\n"

	f, err := os.OpenFile(userJS, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("cannot write %s: %w", userJS, err)
	}
	defer f.Close()
	_, err = f.WriteString(marker)
	return err
}

// FirefoxIsRunning returns true if a Firefox process is currently running.
func FirefoxIsRunning() bool {
	if err := exec.Command("pgrep", "-x", "firefox").Run(); err == nil {
		return true
	}
	// Also check for "firefox-bin" (some Linux packaging uses this name).
	return exec.Command("pgrep", "-x", "firefox-bin").Run() == nil
}
