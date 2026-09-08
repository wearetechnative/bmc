package cmd

import (
	"io"
	"testing"

	"github.com/spf13/cobra"
)

func TestResolveProfileMode(t *testing.T) {
	tests := []struct {
		name       string
		profile    string
		pick       bool
		envProfile string
		wantMode   profileMode
		wantValue  string
	}{
		{"pick wins over everything", "TN-Production", true, "TN-Env", modeInteractive, ""},
		{"pick with nothing else", "", true, "", modeInteractive, ""},
		{"named profile", "TN-Production", false, "", modeNamed, "TN-Production"},
		{"named profile beats env", "TN-Production", false, "TN-Env", modeNamed, "TN-Production"},
		{"named profile trimmed", "  TN-Production  ", false, "", modeNamed, "TN-Production"},
		{"env used when no profile", "", false, "TN-Env", modeEnv, "TN-Env"},
		{"whitespace profile falls through to env", "   ", false, "TN-Env", modeEnv, "TN-Env"},
		{"nothing set falls back to interactive", "", false, "", modeInteractive, ""},
		{"whitespace everywhere is interactive", "  ", false, "  ", modeInteractive, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, value := resolveProfileMode(tt.profile, tt.pick, tt.envProfile)
			if mode != tt.wantMode {
				t.Errorf("mode = %v, want %v", mode, tt.wantMode)
			}
			if value != tt.wantValue {
				t.Errorf("value = %q, want %q", value, tt.wantValue)
			}
		})
	}
}

// newProfileFlagTestCmd builds a throwaway command wired with addProfileFlags,
// capturing the positional args it receives. It resets the shared flag globals.
func newProfileFlagTestCmd(gotArgs *[]string) *cobra.Command {
	globalProfile = ""
	globalPick = false
	cmd := &cobra.Command{
		Use:           "test [search]",
		Args:          cobra.MaximumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			*gotArgs = args
			return nil
		},
	}
	addProfileFlags(cmd)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	return cmd
}

func TestProfileFlagParsing(t *testing.T) {
	t.Run("space form: -p NAME with positional search", func(t *testing.T) {
		var args []string
		cmd := newProfileFlagTestCmd(&args)
		cmd.SetArgs([]string{"-p", "TN-Production", "compute2"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if globalProfile != "TN-Production" {
			t.Errorf("globalProfile = %q, want %q", globalProfile, "TN-Production")
		}
		if globalPick {
			t.Errorf("globalPick = true, want false")
		}
		if len(args) != 1 || args[0] != "compute2" {
			t.Errorf("args = %v, want [compute2]", args)
		}
	})

	t.Run("equals form: -p=NAME with positional search", func(t *testing.T) {
		var args []string
		cmd := newProfileFlagTestCmd(&args)
		cmd.SetArgs([]string{"-p=TN-Production", "compute2"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if globalProfile != "TN-Production" {
			t.Errorf("globalProfile = %q, want %q", globalProfile, "TN-Production")
		}
		if len(args) != 1 || args[0] != "compute2" {
			t.Errorf("args = %v, want [compute2]", args)
		}
	})

	t.Run("--pick with positional search", func(t *testing.T) {
		var args []string
		cmd := newProfileFlagTestCmd(&args)
		cmd.SetArgs([]string{"--pick", "compute2"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !globalPick {
			t.Errorf("globalPick = false, want true")
		}
		if globalProfile != "" {
			t.Errorf("globalProfile = %q, want empty", globalProfile)
		}
		if len(args) != 1 || args[0] != "compute2" {
			t.Errorf("args = %v, want [compute2]", args)
		}
	})

	t.Run("-P short form sets pick", func(t *testing.T) {
		var args []string
		cmd := newProfileFlagTestCmd(&args)
		cmd.SetArgs([]string{"-P"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !globalPick {
			t.Errorf("globalPick = false, want true")
		}
	})

	t.Run("bare -p errors (needs an argument)", func(t *testing.T) {
		var args []string
		cmd := newProfileFlagTestCmd(&args)
		cmd.SetArgs([]string{"-p"})
		err := cmd.Execute()
		if err == nil {
			t.Fatalf("expected error for bare -p, got nil")
		}
	})

	t.Run("no flags: only positional search", func(t *testing.T) {
		var args []string
		cmd := newProfileFlagTestCmd(&args)
		cmd.SetArgs([]string{"compute2"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if globalProfile != "" || globalPick {
			t.Errorf("expected no profile flags set, got profile=%q pick=%v", globalProfile, globalPick)
		}
		if len(args) != 1 || args[0] != "compute2" {
			t.Errorf("args = %v, want [compute2]", args)
		}
	})
}
