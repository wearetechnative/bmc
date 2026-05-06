package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const shellWrapper = `
# bmc shell integration — added by bmc install-shell-integration
bmc() {
  if [[ "$1" == "profsel" ]]; then
    eval "$(command bmc profsel "$@")"
  else
    command bmc "$@"
  fi
}
`

const manualInstructions = `Cannot write to %s (permission denied).
This usually means the file is managed by home-manager or another tool.

Add the wrapper manually using one of these methods:

── home-manager zsh (home.nix) ──────────────────────────────
  programs.zsh.initContent = ''
    bmc() {
      if [[ "$1" == "profsel" ]]; then
        eval "$(command bmc profsel "$@")"
      else
        command bmc "$@"
      fi
    }
  '';

── home-manager bash (home.nix) ─────────────────────────────
  programs.bash.initContent = ''
    bmc() {
      if [[ "$1" == "profsel" ]]; then
        eval "$(command bmc profsel "$@")"
      else
        command bmc "$@"
      fi
    }
  '';

── manual ~/.zshrc or ~/.bashrc ─────────────────────────────
  bmc() {
    if [[ "$1" == "profsel" ]]; then
      eval "$(command bmc profsel "$@")"
    else
      command bmc "$@"
    fi
  }

── fish (~/.config/fish/config.fish) ────────────────────────
  function bmc
    if test "$argv[1]" = "profsel"
      eval (command bmc profsel $argv)
    else
      command bmc $argv
    end
  end
`

var installShellCmd = &cobra.Command{
	Use:   "install-shell-integration",
	Short: "Install the profsel shell wrapper into your shell rc file",
	RunE:  runInstallShell,
}

func init() {
	rootCmd.AddCommand(installShellCmd)
}

func runInstallShell(cmd *cobra.Command, args []string) error {
	shell := os.Getenv("SHELL")

	var rcFile string
	switch {
	case strings.HasSuffix(shell, "zsh"):
		rcFile = filepath.Join(os.Getenv("HOME"), ".zshrc")
	case strings.HasSuffix(shell, "bash"):
		rcFile = filepath.Join(os.Getenv("HOME"), ".bashrc")
	default:
		fmt.Printf("Unsupported shell: %s\n", shell)
		fmt.Printf("Add the following to your shell rc file manually:\n%s", shellWrapper)
		return nil
	}

	// Check if already installed
	existing, err := os.ReadFile(rcFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read %s: %w", rcFile, err)
	}
	if strings.Contains(string(existing), "bmc shell integration") {
		fmt.Printf("Shell integration already installed in %s\n", rcFile)
		return nil
	}

	f, err := os.OpenFile(rcFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			fmt.Printf(manualInstructions, rcFile)
			return nil
		}
		return fmt.Errorf("failed to open %s: %w", rcFile, err)
	}
	defer f.Close()

	if _, err := f.WriteString(shellWrapper); err != nil {
		return fmt.Errorf("failed to write to %s: %w", rcFile, err)
	}

	fmt.Printf("Shell integration installed in %s\n", rcFile)
	fmt.Printf("Restart your shell or run: source %s\n", rcFile)
	return nil
}
