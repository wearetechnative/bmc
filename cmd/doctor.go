package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsconfig"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/prereqs"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check bmc prerequisites and configuration",
	RunE:  runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nbmc v%s — system check\n\n", Version)

	allOK := true

	// --- Core ---
	fmt.Println("Core (always required)")
	allOK = doctorCheck("~/.aws/config", checkAWSConfig()) && allOK
	allOK = doctorCheck("~/.aws/credentials", checkAWSCredentials()) && allOK
	allOK = doctorCheck("~/.config/bmc/config.json", checkBMCConfig()) && allOK

	fmt.Println()

	// --- Optional ---
	fmt.Println("ec2connect / ecsconnect (optional)")
	sshPath := prereqs.FindPath("ssh")
	allOK = doctorCheck("ssh", checkBinary("ssh", sshPath)) && allOK

	awsPath := prereqs.FindPath("aws")
	awsVersion := prereqs.AWSVersion()
	allOK = doctorCheck("aws CLI", checkAWSBinary(awsPath, awsVersion)) && allOK

	smpPath := prereqs.FindPath("session-manager-plugin")
	allOK = doctorCheck("session-manager-plugin", checkBinary("session-manager-plugin", smpPath)) && allOK

	fmt.Println()

	// --- MFA ---
	fmt.Println("MFA")
	cfg, _ := config.Load()
	allOK = doctorCheck("MFA enabled", checkMFAEnabled(cfg)) && allOK
	if cfg.MFA.Enabled {
		allOK = doctorCheck("totp_script configured", checkTOTPScript(cfg)) && allOK
		if len(cfg.MFA.ProfileScripts) > 0 {
			fmt.Printf("  ✓ profile_scripts: %d override(s) configured\n", len(cfg.MFA.ProfileScripts))
			for profile := range cfg.MFA.ProfileScripts {
				fmt.Printf("      %s\n", profile)
			}
		}
		allOK = doctorCheck("copy_command", checkClipboard(cfg)) && allOK
	}

	fmt.Println()

	// --- Shell integration ---
	fmt.Println("Shell integration")
	allOK = doctorCheck("profsel wrapper installed", checkShellWrapper()) && allOK

	// --- Legacy ---
	if config.HasLegacyConfig() {
		fmt.Println()
		fmt.Println("⚠  Legacy config.env found:")
		fmt.Printf("   %s\n", config.LegacyConfigPath())
		fmt.Println("   Please migrate to ~/.config/bmc/config.json")
		fmt.Println("   See: https://github.com/wearetechnative/bmc#configuration")
		allOK = false
	}

	fmt.Println()

	if !allOK {
		return fmt.Errorf("some checks failed — see above for details")
	}
	fmt.Println("All checks passed ✓")
	return nil
}

// doctorCheck prints a check result and returns true if passed.
func doctorCheck(label, detail string) bool {
	if detail == "" {
		fmt.Printf("  ✓ %s\n", label)
		return true
	}
	fmt.Printf("  ✗ %s\n", label)
	for _, line := range strings.Split(detail, "\n") {
		if line != "" {
			fmt.Printf("    %s\n", line)
		}
	}
	return false
}

func checkAWSConfig() string {
	profiles, err := awsconfig.LoadProfiles()
	if err != nil {
		return err.Error()
	}
	// Print count inline
	fmt.Printf("  ✓ ~/.aws/config found (%d profiles)\n", len(profiles))
	return "skip"
}

func checkAWSCredentials() string {
	path := os.Getenv("HOME") + "/.aws/credentials"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "~/.aws/credentials not found — MFA refresh will not be available"
	}
	return ""
}

func checkBMCConfig() string {
	path := config.ConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "~/.config/bmc/config.json not found — defaults will be used\n" +
			"Create it to configure MFA and EC2 behaviour."
	}
	_, err := config.Load()
	if err != nil {
		return err.Error()
	}
	return ""
}

func checkBinary(name, path string) string {
	if path == "" {
		return name + " not found\n" + formatInstallInstructions(name)
	}
	fmt.Printf("  ✓ %s found (%s)\n", name, path)
	return "skip"
}

func checkAWSBinary(path, version string) string {
	if path == "" {
		return "aws CLI not found\n" + formatInstallInstructions("aws")
	}
	v1warning := ""
	if strings.Contains(version, "aws-cli/1.") {
		v1warning = " ⚠ v1 detected — v2 required"
	}
	fmt.Printf("  ✓ aws CLI found (%s)%s\n", version, v1warning)
	return "skip"
}

func checkMFAEnabled(cfg config.Config) string {
	if cfg.MFA.Enabled {
		return ""
	}
	return "MFA disabled in config (mfa.enabled = false)"
}

func checkTOTPScript(cfg config.Config) string {
	if cfg.MFA.TOTPScript != "" {
		return ""
	}
	return "totp_script not configured — MFA codes must be entered manually"
}

func checkClipboard(cfg config.Config) string {
	if cfg.MFA.CopyCommand == "" {
		return "copy_command not configured (optional)"
	}
	path := prereqs.FindPath(strings.Fields(cfg.MFA.CopyCommand)[0])
	if path == "" {
		return "copy_command binary not found: " + cfg.MFA.CopyCommand + "\n" +
			"Install: sudo apt install xclip  (clipboard copy will be skipped)"
	}
	return ""
}

func checkShellWrapper() string {
	shell := os.Getenv("SHELL")
	var rcFile string
	if strings.HasSuffix(shell, "zsh") {
		rcFile = os.Getenv("HOME") + "/.zshrc"
	} else {
		rcFile = os.Getenv("HOME") + "/.bashrc"
	}
	data, err := os.ReadFile(rcFile)
	if err != nil {
		return "could not read " + rcFile
	}
	if strings.Contains(string(data), "bmc shell integration") {
		return ""
	}
	return "Shell wrapper not installed\n" +
		"Run: bmc install-shell-integration"
}

func formatInstallInstructions(binary string) string {
	switch binary {
	case "ssh":
		return "apt:          sudo apt install openssh-client\n" +
			"brew:         brew install openssh\n" +
			"nix-env:      nix-env -iA nixpkgs.openssh\n" +
			"nix profile:  nix profile add nixpkgs#openssh\n" +
			"NixOS config: environment.systemPackages = [ pkgs.openssh ];"
	case "aws":
		return "brew:         brew install awscli\n" +
			"nix-env:      nix-env -iA nixpkgs.awscli2\n" +
			"nix profile:  nix profile add nixpkgs#awscli2\n" +
			"NixOS config: environment.systemPackages = [ pkgs.awscli2 ];"
	case "session-manager-plugin":
		return "brew:         brew install session-manager-plugin\n" +
			"nix-env:      nix-env -iA nixpkgs.session-manager-plugin\n" +
			"nix profile:  nix profile add nixpkgs#session-manager-plugin\n" +
			"NixOS config: environment.systemPackages = [ pkgs.session-manager-plugin ];"
	}
	return ""
}
