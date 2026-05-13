## Context

`cmd/install_shell.go` currently handles zsh and bash via a switch on `$SHELL`. It has one shared `shellWrapper` const and one `manualInstructions` const (shown on permission errors). The Fish manual snippet already exists inside `manualInstructions` but is never shown to Fish users hitting the `default` case.

Fish differs from bash/zsh in two ways relevant here:
1. Functions live in `~/.config/fish/functions/<name>.fish` and are auto-loaded — no `source` needed
2. Syntax differs (`function`/`end`, no `[[ ]]`, `eval (...)` not `eval "$(...)"`

NixOS users cannot write to system-managed config paths; they configure shell integration declaratively via home-manager.

## Goals / Non-Goals

**Goals:**
- Fish users get automatic installation to `~/.config/fish/functions/bmc.fish`
- NixOS Fish users get clear manual instructions with a `programs.fish.functions` home-manager snippet
- Already-installed detection for Fish via file existence check
- No changes to zsh/bash paths

**Non-Goals:**
- Supporting other shells (nush, xonsh, etc.)
- NixOS detection for zsh/bash (the existing permission-denied path already handles this)
- Modifying the Fish wrapper content beyond the already-proven snippet

## Decisions

### Write to functions directory, not config.fish

**Decision**: `~/.config/fish/functions/bmc.fish`

**Rationale**: This is the idiomatic Fish way — functions auto-load without any `source` or shell restart. It also makes already-installed detection trivial (file existence). The alternative (`config.fish` append) would be consistent with bash/zsh but less correct for Fish.

### NixOS detection via /etc/nixos/

**Decision**: Check if `/etc/nixos/` directory exists. If so, treat as NixOS regardless of `$SHELL`.

**Rationale**: `/etc/nixos/` is present on all NixOS systems with a configuration. It's simpler and more reliable than parsing `/etc/os-release`. The check is done once before the shell switch — if NixOS is detected and `$SHELL` is fish, show manual instructions and return early.

**Alternative considered**: Check `/etc/os-release` for `ID=nixos`. More portable but requires file parsing. `/etc/nixos/` existence is sufficient and already the convention in the codebase's manual instructions context.

### Already-installed check: file existence

**Decision**: Check if `~/.config/fish/functions/bmc.fish` exists.

**Rationale**: Unlike bash/zsh (which check for a marker string in a shared rc file), Fish functions each have their own file. File existence is unambiguous. This check is skipped on NixOS since we never write the file there anyway.

### Manual instructions content

The home-manager snippet uses `programs.fish.functions` (the declarative function approach), not `programs.fish.shellInit`. This is the correct home-manager idiom for Fish functions and keeps the function isolated from other shell init code.

## Risks / Trade-offs

- **`/etc/nixos/` false positive** → Unlikely: the directory is NixOS-specific. A non-NixOS user with `/etc/nixos/` would get manual instructions instead of auto-install — conservative, not harmful.
- **Fish functions dir doesn't exist yet** → `os.MkdirAll` ensures the directory is created before writing.
- **File already exists with different content** → We report "already installed" and don't overwrite. User can delete the file to reinstall. Acceptable trade-off for safety.
