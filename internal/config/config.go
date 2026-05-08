package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds all bmc runtime configuration.
type Config struct {
	MFA     MFAConfig     `toml:"mfa"`
	EC2     EC2Config     `toml:"ec2"`
	Console ConsoleConfig `toml:"console"`
}

// ConsoleConfig holds console-related settings.
type ConsoleConfig struct {
	FirefoxContainers bool `toml:"firefox_containers"`
}

// MFAConfig holds MFA-related settings.
type MFAConfig struct {
	Enabled          bool   `toml:"enabled"`
	TOTPScript       string `toml:"totp_script"`
	ClipboardCommand string `toml:"clipboard_command"`
}

// EC2Config holds EC2-related settings.
type EC2Config struct {
	AutoStartStopped string   `toml:"auto_start_stopped"` // always | never | prompt
	Columns          []string `toml:"columns"`
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
	}
}

// ConfigPath returns the path to the bmc TOML config file.
func ConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "bmc", "config.toml")
}

// LegacyConfigPath returns the path to the legacy bash config file.
func LegacyConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "bmc", "config.env")
}

// Load reads ~/.config/bmc/config.toml. If the file is absent, defaults are returned.
// An invalid TOML file returns an error with the file path.
func Load() (Config, error) {
	cfg := Defaults()
	path := ConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, err
	}

	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return cfg, &ParseError{Path: path, Err: err}
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
