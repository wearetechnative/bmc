## Why

The Go rewrite of BMC is complete and released. Several bash-era files remain in the repository that served the old shell-based implementation and are no longer used. Additionally, the `bmc-go` binary was accidentally committed as a tracked file. Removing these reduces confusion and keeps the repository focused on the current Go codebase.

## What Changes

- Remove `_bmclib.sh` — the original bash library that implemented BMC's core functionality before the Go rewrite
- Remove `_get_var_file.sh` — a terraform variable file helper from the bash era
- Remove `tgselect.sh` — a Toggl time tracker selector script that landed in this repo during the bash era
- Remove `bmc-go` from git tracking and add `/bmc-go` to `.gitignore` (accidentally committed build artifact)
- Keep `release.sh` — still actively used for release automation
- Keep `release-script-prompt.txt` — still used as reference for release tooling

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None. This is a pure cleanup — no functional requirements change.

## Impact

- `.gitignore`: add `/bmc-go` entry
- No Go source code changes
- No user-facing behavior changes
- Related bean: [bmc-scyt](/.beans/bmc-scyt--cleanup-bash-files.md)
