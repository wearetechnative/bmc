## Why

`bmc console` already surfaces recently used profiles at the top of its interactive selector, but `bmc profsel` and all other commands that call `ensureAWSProfile()` (e.g. `ec2connect`, `ec2ls`) do not. Users who frequently switch between a small set of AWS profiles have to scroll the full list every time they use any command other than `console`.

## What Changes

- Introduce a shared history key `"profile"` (file: `~/.local/share/bmc/profile-history.json`) used by all interactive profile selectors
- Refactor `selectProfileForConsoleInteractive()` into a general `selectProfileWithHistory()` in `cmd/profilehelper.go`, reusable by all callers
- Change `bmc console` from history key `"console"` to `"profile"`
- Wire `bmc profsel` to use `selectProfileWithHistory()` and save to history after interactive selection
- Wire `ensureAWSProfile()` to use `selectProfileWithHistory()` and save to history after interactive selection
- Drop the now-unused `selectProfileInteractive()` function

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `profile-selection`: Interactive profile selector now surfaces recently used groups and profiles at the top (with "recent" label), shared across all commands that trigger interactive selection

## Impact

- `cmd/profilehelper.go`: new `selectProfileWithHistory()` replacing `selectProfileInteractive()`; `ensureAWSProfile()` calls it and saves to history
- `cmd/console.go`: switches to `selectProfileWithHistory()`, history key changes from `"console"` to `"profile"`
- `cmd/profsel.go`: switches to `selectProfileWithHistory()`, adds `history.Save("profile", ...)` after interactive selection
- Existing `~/.local/share/bmc/console-history.json` is not migrated; history resets once on upgrade
- No new external dependencies
