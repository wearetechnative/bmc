package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds all bmc runtime configuration.
type Config struct {
	MFA     MFAConfig     `json:"mfa"`
	EC2     EC2Config     `json:"ec2"`
	Console ConsoleConfig `json:"console"`
	Watcher WatcherConfig `json:"watcher"`
}

// WatcherConfig holds settings for the bmc watcher daemon.
type WatcherConfig struct {
	// FirefoxDebugPort is the CDP remote debugging port Firefox listens on.
	// Set to 0 to disable CDP and always use the tab-based fallback.
	FirefoxDebugPort int `json:"firefox_debug_port"`
}

// ConsoleConfig holds console-related settings.
type ConsoleConfig struct {
	FirefoxContainers bool   `json:"firefox_containers"`
	ChromeProfiles    bool   `json:"chrome_profiles"`
	ChromeBinary      string `json:"chrome_binary"`
}

// MFAConfig holds MFA-related settings.
type MFAConfig struct {
	Enabled        bool              `json:"enabled"`
	TOTPScript     string            `json:"totp_script"`
	ProfileScripts map[string]string `json:"profile_scripts,omitempty"`
	CopyCommand    string            `json:"copy_command"`
	PasteCommand   string            `json:"paste_command"`
}

// EC2Config holds EC2-related settings.
type EC2Config struct {
	AutoStartStopped string   `json:"auto_start_stopped"` // always | never | prompt
	Columns          []string `json:"columns"`
}

// Defaults returns a Config with sensible default values.
func Defaults() Config {
	return Config{
		MFA: MFAConfig{
			Enabled: false,
		},
		EC2: EC2Config{
			AutoStartStopped: "prompt",
			Columns:          []string{"InstanceId", "Name", "PrivateIP", "PublicIP", "State", "Hibernate", "Scheduler"},
		},
		Watcher: WatcherConfig{
			FirefoxDebugPort: 9222,
		},
	}
}

// ConfigPath returns the path to the bmc JSON config file.
func ConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "bmc", "config.json")
}

// LegacyConfigPath returns the path to the legacy bash config file.
func LegacyConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "bmc", "config.env")
}

// tomlConfigPath returns the path to the old TOML config file.
func tomlConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "bmc", "config.toml")
}

// Load reads ~/.config/bmc/config.json. If the file is absent, defaults are returned.
// If config.json is absent but config.toml exists, a migration hint is printed to stderr.
// An invalid JSON file returns an error with the file path.
func Load() (Config, error) {
	cfg := Defaults()
	path := ConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Print migration hint if old TOML config exists
			if _, tomlErr := os.Stat(tomlConfigPath()); tomlErr == nil {
				fmt.Fprintf(os.Stderr, "bmc: config.toml found but config.json is missing.\n")
				fmt.Fprintf(os.Stderr, "bmc: Convert your config to JSON format and save it to %s\n", path)
				fmt.Fprintf(os.Stderr, "bmc: Example: {\"mfa\":{\"enabled\":true},\"console\":{\"firefox_containers\":true}}\n")
			}
			return cfg, nil
		}
		return cfg, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, &ParseError{Path: path, Err: err}
	}

	if cfg.Console.ChromeBinary == "" {
		cfg.Console.ChromeBinary = "google-chrome"
	}

	return cfg, nil
}

// HasLegacyConfig returns true if ~/.config/bmc/config.env exists.
func HasLegacyConfig() bool {
	_, err := os.Stat(LegacyConfigPath())
	return err == nil
}

// ParseError is returned when the config file is present but malformed.
type ParseError struct {
	Path string
	Err  error
}

func (e *ParseError) Error() string {
	return "failed to parse " + e.Path + ": " + e.Err.Error()
}
