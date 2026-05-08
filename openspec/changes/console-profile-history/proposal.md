## Why

When using `bmc console` interactively, users frequently open the same AWS accounts. Having to scroll through a full profile list each time is unnecessary friction — recently used accounts should appear at the top for quick access.

Relates to: https://github.com/wearetechnative/bmc/issues/28

## What Changes

- `bmc console` saves the selected profile to a local history file after each successful console open
- The interactive profile selector shows a "Recent" section at the top with the last N used profiles
- History is only shown when `AWS_PROFILE` is not set and interactive selection is triggered

## Capabilities

### New Capabilities

- `console-profile-history`: Persists and displays recently used console profiles, showing them at the top of the interactive selector

### Modified Capabilities

- `aws-console-access`: Profile selector now shows a "Recent" section above the full profile list when history entries exist

## Impact

- `cmd/console.go`: write selected profile to history after successful open; pass history to selector
- `internal/history/` (new package): read/write history file (`~/.local/share/bmc/console-history.json`)
- `internal/ui/` or `cmd/profilehelper.go`: extend profile selector to accept and display a "Recent" section
