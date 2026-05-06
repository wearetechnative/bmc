package prereqs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Prereq describes an optional binary dependency.
type Prereq struct {
	Binary   string
	Commands []string // commands that require this binary
}

var (
	SSH = Prereq{
		Binary:   "ssh",
		Commands: []string{"ec2connect (SSH method)"},
	}
	AWSCLI = Prereq{
		Binary:   "aws",
		Commands: []string{"ec2connect (SSM)", "ecsconnect"},
	}
	SessionManagerPlugin = Prereq{
		Binary:   "session-manager-plugin",
		Commands: []string{"ec2connect (SSM)", "ecsconnect"},
	}
)

// Check verifies that the given binary is available. Returns nil if found.
// Returns a formatted error with install instructions if not found.
func Check(p Prereq) error {
	if _, err := exec.LookPath(p.Binary); err == nil {
		return nil
	}
	return &MissingError{Prereq: p}
}

// CheckAWSCLIVersion checks for aws CLI and verifies it's v2.
func CheckAWSCLIVersion() error {
	path, err := exec.LookPath("aws")
	if err != nil {
		return &MissingError{Prereq: AWSCLI}
	}
	out, err := exec.Command(path, "--version").Output()
	if err != nil {
		return nil // found but version check failed — continue anyway
	}
	version := string(out)
	if strings.Contains(version, "aws-cli/1.") {
		fmt.Fprintln(os.Stderr, "⚠  Warning: aws CLI v1 detected. bmc requires aws CLI v2.")
		fmt.Fprintln(os.Stderr, "   "+installInstructions("aws")[0])
	}
	return nil
}

// MissingError is returned when a required binary is not found.
type MissingError struct {
	Prereq Prereq
}

func (e *MissingError) Error() string {
	lines := []string{
		fmt.Sprintf("✗ %s not found", e.Prereq.Binary),
		fmt.Sprintf("  Required for: %s", strings.Join(e.Prereq.Commands, ", ")),
		"  Install:",
	}
	for _, inst := range installInstructions(e.Prereq.Binary) {
		lines = append(lines, "    "+inst)
	}
	return strings.Join(lines, "\n")
}

func installInstructions(binary string) []string {
	switch binary {
	case "ssh":
		return []string{
			"apt:          sudo apt install openssh-client",
			"brew:         brew install openssh",
			"nix-env:      nix-env -iA nixpkgs.openssh",
			"nix profile:  nix profile add nixpkgs#openssh",
			"NixOS config: environment.systemPackages = [ pkgs.openssh ];",
		}
	case "aws":
		return []string{
			"See: https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html",
			"brew:         brew install awscli",
			"nix-env:      nix-env -iA nixpkgs.awscli2",
			"nix profile:  nix profile add nixpkgs#awscli2",
			"NixOS config: environment.systemPackages = [ pkgs.awscli2 ];",
		}
	case "session-manager-plugin":
		return []string{
			"See: https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-plugin-install.html",
			"brew:         brew install session-manager-plugin",
			"nix-env:      nix-env -iA nixpkgs.session-manager-plugin",
			"nix profile:  nix profile add nixpkgs#session-manager-plugin",
			"NixOS config: environment.systemPackages = [ pkgs.session-manager-plugin ];",
		}
	}
	return []string{"See distribution docs for install instructions."}
}

// FindPath returns the path to a binary, or "" if not found.
func FindPath(binary string) string {
	path, _ := exec.LookPath(binary)
	return path
}

// AWSVersion returns the aws CLI version string.
func AWSVersion() string {
	out, err := exec.Command("aws", "--version").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
