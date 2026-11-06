# Proposal

## Why

`bmc ec2ls` only offers a bordered table or JSON, which makes quick shell processing (`grep`, `awk`, `cut`) awkward, and every AWS command prints MFA status lines ("-- Using AWS source-profile", "Current MFA Session Valid") that are debug information rather than something the user needs to act on. In addition, SSM is the connection method used most often with `ec2connect`, yet SSH is the default in the picker.

## What Changes

- `ec2connect`: the connection method picker lists `ssm` first, so pressing Enter starts an SSM session. Behaviour when SSH flags are given is unchanged.
- Printed tables (`ec2ls`, `ec2find`) pad every cell with one space between the border and the content.
- New `bmc ec2ls --txt` flag: outputs instances as aligned plain text without border lines, using the configured `ec2.columns`. Mutually exclusive with `--json`.
- New `bmc ec2ls --no-headers` flag (requires `--txt`): omits the header line, so every output line is an instance.
- New global `--debug` flag on the root command.
- **BREAKING** (output only): the MFA status messages "-- Using AWS source-profile: {profile}" and "Current MFA Session Valid, until: {time}", and the echo of the TOTP code produced by the TOTP script, are no longer shown by default. They are shown when `--debug` is given. Messages that inform about or require user action (session refresh, TOTP script execution, manual TOTP entry, refreshed-until confirmation) remain always visible.
- The ec2ls-specific `mfa.Quiet` switch that was added during exploration is replaced by the global `--debug` mechanism.

## Capabilities

### New Capabilities
- `ec2-text-output`: plain-text (borderless, aligned) output for `ec2ls`, including `--no-headers`.
- `debug-output`: global `--debug` flag that enables diagnostic output across commands.

### Modified Capabilities
- `mfa-authentication`: the source-profile message and the "session valid" message become debug-only; action-related messages remain always visible.
- `ec2connect`: when the connection method picker is shown, SSM is listed first and is the default selection.
- `ec2ls`: the outdated `gum table` table requirement is replaced by the current rendered table, including cell padding.

## Impact

- Code: `cmd/root.go` (persistent `--debug` flag), `cmd/ec2ls.go` (`--txt`, `--no-headers`), `cmd/ec2connect.go` (picker order), `internal/ui/table.go` (`PrintText`, padding in `PrintTable`), `internal/mfa/session.go` (debug-gated messages).
- Users and scripts that parse the MFA status lines from stderr will no longer see them unless `--debug` is passed.
- Docs: `docs/content/commands/ec2ls.md`, `docs/content/commands/ec2connect.md`, a mention of `--debug` in the docs site, and `CHANGELOG.md`.
- No new dependencies.
