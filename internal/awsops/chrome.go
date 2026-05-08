package awsops

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/wearetechnative/bmc/internal/config"
)

var invalidPathChars = regexp.MustCompile(`[/:\\*?"<>|]`)

func sanitizeProfileName(name string) string {
	return invalidPathChars.ReplaceAllString(name, "-")
}

func chromeProfileDir(profileName string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "bmc", "chrome", "profiles", sanitizeProfileName(profileName))
}

func defaultChromeProfilePath() string {
	home, _ := os.UserHomeDir()

	var candidates []string
	switch runtime.GOOS {
	case "linux":
		candidates = []string{
			filepath.Join(home, ".config", "google-chrome", "Default"),
			filepath.Join(home, ".config", "chromium", "Default"),
		}
	case "darwin":
		candidates = []string{
			filepath.Join(home, "Library", "Application Support", "Google", "Chrome", "Default"),
			filepath.Join(home, "Library", "Application Support", "Chromium", "Default"),
		}
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func seedChromeProfile(destDir, srcDir string) error {
	entries := []string{"Extensions", "Local Extension Settings", "Preferences"}
	defaultDir := filepath.Join(destDir, "Default")
	if err := os.MkdirAll(defaultDir, 0700); err != nil {
		return err
	}
	srcDefault := srcDir

	for _, entry := range entries {
		src := filepath.Join(srcDefault, entry)
		dst := filepath.Join(defaultDir, entry)

		info, err := os.Stat(src)
		if err != nil {
			continue // skip missing entries silently
		}

		if info.IsDir() {
			_ = copyDir(src, dst)
		} else {
			_ = copyFile(src, dst)
		}
	}
	return nil
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0700); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		s := filepath.Join(src, entry.Name())
		d := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyDir(s, d); err != nil {
				continue
			}
		} else {
			_ = copyFile(s, d)
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func openChromeProfile(u, profileName string, cfg config.ConsoleConfig) error {
	profileDir := chromeProfileDir(profileName)

	if _, err := os.Stat(profileDir); os.IsNotExist(err) {
		srcDir := defaultChromeProfilePath()
		if srcDir != "" {
			_ = seedChromeProfile(profileDir, srcDir)
		} else {
			if err := os.MkdirAll(filepath.Join(profileDir, "Default"), 0700); err != nil {
				return err
			}
		}
	}

	binary := cfg.ChromeBinary
	chromePath, err := exec.LookPath(binary)
	if err != nil {
		return fmt.Errorf("%s not found in PATH (required for chrome_profiles): install it or set chrome_binary in config", binary)
	}
	cmd := exec.Command(chromePath,
		"--user-data-dir="+profileDir,
		"--no-first-run",
		"--no-default-browser-check",
		u,
	)
	cmd.Stderr = nil
	return cmd.Start()
}
