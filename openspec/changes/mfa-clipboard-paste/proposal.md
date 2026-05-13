## Why

The Go rewrite collapsed the original `clipboardCopyCommand` and `clipboardPasteCommand` config fields into a single `clipboard_command` (copy only), losing paste functionality. Users who relied on automatic keystroke simulation to fill MFA codes into browser windows must now paste manually. This restores feature parity with the bash version.

## What Changes

- **BREAKING**: `mfa.clipboard_command` is renamed to `mfa.copy_command` in `~/.config/bmc/config.json`
- New config field `mfa.paste_command`: optional command run after a successful copy (keystroke simulation, e.g. `xdotool key ctrl+v`)
- After a successful copy, bmc waits 300ms (hardcoded) then runs `paste_command` if configured
- `paste_command` receives no stdin — it is a pure keystroke simulator
- If copy fails, `paste_command` is skipped
- `bmc doctor` updated to check `copy_command` instead of `clipboard_command`
- `totp_script` is now executed via `sh -c` — quoted arguments with spaces work correctly (e.g. `rbw code "My Entry (new)"`)
- `totp_script` receives the parent's stdin and stderr — interactive TUI selection tools (e.g. `gum filter`) render correctly and accept keyboard input

## Capabilities

### New Capabilities

- `mfa-clipboard-paste`: Automatic paste of TOTP codes via configurable keystroke simulation after clipboard copy

### Modified Capabilities

- `mfa-authentication`: The clipboard integration requirement changes — `clipboard_command` → `copy_command`, paste step added

## Impact

- `internal/config/config.go`: rename `ClipboardCommand` → `CopyCommand`, add `PasteCommand`
- `internal/mfa/session.go`: update `copyToClipboard` to handle copy → delay → paste
- `cmd/doctor.go`: update field reference from `ClipboardCommand` to `CopyCommand`
- `README.md`: update config reference table
- `CHANGELOG.md`: document breaking change
- Linked bean: `.beans/bmc-kx3p--clipboard-copy-paste.md`
