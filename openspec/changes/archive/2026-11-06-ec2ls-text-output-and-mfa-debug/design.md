# Design

## Context

Part of this change was implemented exploratively before the proposal existed and is already in the working tree: SSM first in the `ec2connect` picker, cell padding in `ui.PrintTable`, `ui.PrintText` with `ec2ls --txt`/`--no-headers`, and an ec2ls-only `mfa.Quiet` package variable that suppressed two MFA status lines. This design keeps the first three and replaces `mfa.Quiet` with a global debug mechanism (see proposal.md - What Changes).

MFA messages are written by `mfa.EnsureValid`, which is reached from three call sites: `runMFACheck` in `cmd/profilehelper.go` (all AWS operation commands), `cmd/profsel.go` and `cmd/console.go`. `bmc profsel` is evaluated by the shell wrapper (`eval "$(command bmc profsel ...)"`), so anything on stdout must stay shell code; all MFA output already goes to stderr.

## Goals / Non-Goals

**Goals:**
- One switch (`--debug`) that controls diagnostic output for every command, without threading a parameter through every call site.
- Default output contains only results and messages that need or report user action.

**Non-Goals:**
- A logging framework, log levels or log files.
- Classifying non-MFA messages (e.g. "-- Opening console for profile", watcher messages) as debug. They can opt in later using the same mechanism.
- `--txt`/`--no-headers` for `ec2find` or other commands.

## Decisions

### Global persistent flag bound to a package-level variable
`rootCmd.PersistentFlags().BoolVar(&debug.Enabled, "debug", false, ...)` in `cmd/root.go`, with `Enabled` living in a tiny new package `internal/debug`. Cobra parses persistent flags before `RunE`, so the value is set before any MFA call.

- *Alternative: pass a `debug bool` into `mfa.EnsureValid`.* Rejected: it changes the signature at three call sites and every future diagnostic would need the same threading.
- *Alternative: keep the variable in `internal/mfa`.* Rejected: debug output is cross-cutting; a neutral package lets other packages use it without importing `mfa`.

### Helper `debug.Printf(w io.Writer, format, args...)`
Writes only when `Enabled` is true. MFA uses it with its existing `outfd`, so the destination (stderr) stays where it is today.

### Classification of MFA messages
| Message                                   | Shown           |
| ----------------------------------------- | --------------- |
| `-- Using AWS source-profile: X`          | debug only      |
| `Current MFA Session Valid, until: T`     | debug only      |
| TOTP code echoed after the TOTP script    | debug only      |
| `!! AWS MFA Device not found...`          | always          |
| `-- Refreshing MFA session for X...`      | always          |
| `-- Executing TOTP script...`             | always          |
| `-- No TOTP script configured...`         | always          |
| `-- MFA session refreshed, valid until`   | always          |

Rule: a message is shown by default when the user has to act on it, has to wait because of it, or something changed. Pure state reporting is debug. The TOTP code echo is debug because it is sensitive and adds nothing when the script succeeded.

### `ec2ls --no-headers` no longer touches MFA output
With MFA status lines hidden by default, the ec2ls-specific `mfa.Quiet` assignment is removed. `--no-headers` only controls the header line.

### Plain-text rendering
`ui.PrintText(columns, rows, headers)` computes column widths with `lipgloss.Width` (handles multi-byte characters), pads each cell to width + 2 spaces and leaves the last column unpadded to avoid trailing whitespace. When `headers` is false the header is excluded from the width calculation. `text/tabwriter` was considered but pads the last column as well and measures bytes, not display width.

## Risks / Trade-offs

- [Scripts or users relying on the "Using AWS source-profile" line to know which profile is active] → Documented in CHANGELOG as an output change; `--debug` restores it. `bmc profsel --json` still exposes `source_profile`.
- [Package-level mutable state for `debug.Enabled`] → Acceptable for a CLI process with a single command per invocation; tests that need it set and reset it explicitly.
- [A silent first run when the MFA session is valid may feel like nothing happened] → Commands still show their own spinner/result, so feedback remains.
