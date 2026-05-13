## Why

`bmc install-shell-integration` only supports zsh and bash. Fish shell users get an "Unsupported shell" message and a bash snippet that doesn't work for them. Fish is increasingly common among DevOps engineers — adding first-class support removes manual friction.

## What Changes

- `bmc install-shell-integration` detects Fish shell via `$SHELL` and writes a Fish wrapper function to `~/.config/fish/functions/bmc.fish` (the idiomatic Fish functions directory, auto-loaded without sourcing)
- On NixOS (detected via `/etc/nixos/` directory or `ID=nixos` in `/etc/os-release`), the command skips auto-write and prints manual instructions with a `programs.fish.functions` home-manager snippet instead
- Already-installed check: if `~/.config/fish/functions/bmc.fish` already exists, report "already installed" and skip (not applicable on NixOS)
- The `default` case for unsupported shells is unchanged

## Capabilities

### New Capabilities

- `fish-shell-integration`: Automatic installation of the bmc profsel wrapper for Fish shell users, with NixOS-aware fallback to manual instructions

### Modified Capabilities

- `shell-integration`: The `install-shell-integration` command gains Fish shell support; existing zsh/bash behaviour is unchanged

## Impact

- `cmd/install_shell.go`: add Fish case, NixOS detection helper, Fish wrapper const, Fish manual instructions
- No other files need changes
- No breaking changes — existing zsh/bash flows are unaffected
- Linked bean: `.beans/bmc-p8w2--fish-support.md`
