## Context

The current `MFAConfig` struct has `ClipboardCommand string` used only for copying. The `copyToClipboard` function in `mfa/session.go` pipes the TOTP code via stdin to the configured command. The paste step (keystroke simulation) was present in the original bash version but never ported.

Current flow:
```
TOTP code → copyToClipboard(cmd, code) → pipe code to cmd stdin → done
```

Target flow:
```
TOTP code → copyToClipboard(copy_cmd, code) → pipe code to cmd stdin
                → success? → time.Sleep(300ms) → run paste_cmd (no stdin)
```

## Goals / Non-Goals

**Goals:**
- Restore paste functionality with minimal code change
- Clean field rename (`clipboard_command` → `copy_command`) as part of this change
- Hardcoded 300ms delay — no configuration needed

**Non-Goals:**
- Configurable delay
- Platform-specific paste logic (user provides the right command for their setup)
- Fallback paste if copy fails

## Decisions

### Hardcoded 300ms delay

**Decision**: `time.Sleep(300 * time.Millisecond)` between copy and paste.

**Rationale**: Clipboard managers and Wayland compositors need a brief moment to register the new clipboard content before a paste keystroke is processed. 300ms is enough for all common setups without being perceptible to the user. Making it configurable adds surface area for no real benefit.

### paste_command receives no stdin

**Decision**: `exec.Command(parts[0], parts[1:]...).Run()` — no stdin attachment.

**Rationale**: The paste command is a keystroke simulator (`xdotool key ctrl+v`, `wtype -k ctrl+v`, `osascript ...`). It doesn't need the TOTP code — the code is already in the clipboard. Attaching stdin would be confusing and potentially harmful.

### paste only runs after successful copy

**Decision**: Only call paste if `copyToClipboard` reports success (no error).

**Rationale**: Pasting an empty or stale clipboard into a focused MFA field would be worse than doing nothing. If copy fails, the user is already warned and can enter the code manually.

### Breaking rename: clipboard_command → copy_command

**Decision**: No backward compatibility shim. Old configs silently use zero value (empty string).

**Rationale**: The field is optional — if not set, clipboard integration is simply disabled. Users who had `clipboard_command` will see no clipboard copy until they update their config, which is the correct signal that a migration is needed. A deprecation warning would require reading both fields which adds unnecessary complexity.

### totp_script executed via sh -c

**Decision**: `exec.Command("sh", "-c", script)` instead of `strings.Fields(script)`.

**Rationale**: `strings.Fields` splits on whitespace, breaking arguments that contain spaces (e.g. `rbw code "Technative AWS (new)"` becomes four separate tokens instead of two). Using `sh -c` delegates argument parsing to the shell, which correctly handles quoting, escaping, and compound commands.

### totp_script receives os.Stdin and os.Stderr

**Decision**: Set `cmd.Stdin = os.Stdin` and `cmd.Stderr = os.Stderr` when running the TOTP script.

**Rationale**: Go's `exec.Command` defaults both to `/dev/null` when unset. Interactive TUI scripts (e.g. using `gum filter`) need stderr to render their UI on the terminal and stdin (or `/dev/tty`) for keyboard input. Without these, such scripts hang silently. Passing the parent's stdin/stderr allows TUI-based TOTP scripts to work correctly while stdout is still captured for the TOTP code.

## Risks / Trade-offs

- **Breaking change for existing configs** → Documented in CHANGELOG. Users need to rename the field.
- **paste_command targeting wrong window** → The focused window at the time of paste receives the keystroke. This is inherent to keystroke simulation and documented in README. Not a bmc concern.
- **300ms too short on slow systems** → Unlikely for clipboard operations. If it is, users can adjust by wrapping the paste command with a `sleep` prefix: `sleep 0.2 && xdotool key ctrl+v`.
