## Context

`bmc console` runs interactive profile selection via `selectProfileInteractive()`, which calls `ui.Choose()` with the full profile list each time. There is no persistence of which profiles were recently used. Users who repeatedly open the same two or three accounts must scroll or type-filter every time.

The `ui.Choose()` function accepts a flat `[]ui.Item` list. Items have a `Title` and optional `Desc` field rendered dimly beside the title.

## Goals / Non-Goals

**Goals:**
- Persist the last N profiles used with `bmc console` to a local file
- Show recent profiles at the top of the interactive selector with a "recent" label
- Zero overhead when history is empty or file does not exist

**Non-Goals:**
- Shared history across commands (`profsel`, `ec2connect`, etc.)
- History for non-interactive invocations (`-p` flag or `AWS_PROFILE` set)
- Configurable history size (hard-coded at 10)

## Decisions

### Decision: `internal/history` package

A small `internal/history` package encapsulates read/write of the history file. This keeps `cmd/console.go` clean and makes the logic independently testable.

```
internal/history/
  history.go   — Load(name), Save(name, profile)
```

`Load(name)` returns `[]string` (profile names, most recent first). `Save(name, profile)` prepends the profile, deduplicates, caps at 10, and writes atomically. `name` is a key (e.g. `"console"`) allowing future reuse without coupling to a fixed filename.

**Alternative considered**: store history inside `~/.config/bmc/config.toml`. Rejected — config.toml is for user-set preferences, not runtime state. Mixing state into config creates noise when users version-control their config.

### Decision: XDG data dir for history file

Path: `~/.local/share/bmc/<name>-history.json`

Follows XDG Base Directory spec (state/data goes in `XDG_DATA_HOME`, not `XDG_CONFIG_HOME`). Falls back to `~/.local/share/bmc/` if `XDG_DATA_HOME` is unset.

### Decision: Prepend recent items in the same `[]ui.Item` slice

Recent profiles are prepended to the list with `Desc: "recent"`. Items already in the recent list are removed from their original position so each profile appears only once. This approach:
- Requires no changes to `ui.Choose()` or `ui.Item`
- Works with the existing filter (users can still type to filter across all items)
- Is self-contained in `cmd/console.go`

**Alternative considered**: a new `ui.ChooseWithSections()` function with a separator item. Rejected — adds UI complexity for a single use case.

### Decision: Write history only after successful open

`history.Save()` is called only after `awsops.OpenConsole()` returns without error. This prevents failed or cancelled invocations from polluting the history.

## Risks / Trade-offs

- **File I/O on startup**: `history.Load()` adds one file read before the selector appears. The file is small (≤10 profile names) and the read is non-blocking — negligible in practice.
- **History file missing/corrupt**: `Load()` returns an empty slice on any error; `Save()` creates parent directories as needed. Neither failure is fatal.
- **Stale history entries**: If a profile is renamed or deleted, it still appears in "recent". Clicking it will fail at the profile lookup step with a clear error. Acceptable — history is a convenience hint, not a contract.

## Migration Plan

No migration needed. History file is created on first use.
