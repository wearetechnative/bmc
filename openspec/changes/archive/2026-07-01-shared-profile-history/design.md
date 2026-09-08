## Context

`bmc console` has a fully-featured interactive profile selector (`selectProfileForConsoleInteractive()` in `cmd/console.go`) that surfaces recently used groups and profiles at the top of the picker. It stores history in `~/.local/share/bmc/console-history.json` using the `internal/history` package.

`bmc profsel` and all commands using `ensureAWSProfile()` (ec2connect, ec2ls, ec2, ec2scheduler, ec2stopstart, ecsconnect) use a separate, simpler `selectProfileInteractive()` in `cmd/profilehelper.go` that has no history awareness.

The history logic currently lives entirely in console.go — duplicating it per-command is not the right path.

## Goals / Non-Goals

**Goals:**
- One shared history file (`profile-history.json`) that all interactive selectors read and write
- One shared function (`selectProfileWithHistory()`) that all callers use
- Remove the now-redundant `selectProfileInteractive()` function

**Non-Goals:**
- Migrating existing `~/.local/share/bmc/console-history.json` data
- Adding history to non-interactive code paths (e.g. `-p` flag, `AWS_PROFILE` env var)

## Decisions

### 1. Single shared history key: `"profile"`

**Decision**: All interactive profile selectors use history key `"profile"` (file: `~/.local/share/bmc/profile-history.json`).

**Rationale**: A profile recently used in `console` is almost certainly relevant when running `ec2connect` next. Separate per-command history files would fragment what should be a single signal of "which profiles does this user actually use."

**Alternative considered**: Per-command keys (`"console"`, `"profsel"`, `"ec2connect"` etc.) — rejected because it defeats the purpose; a user switching from console to ec2connect would see no recent profiles until they've used ec2connect interactively before.

### 2. Shared function in `profilehelper.go`

**Decision**: Extract `selectProfileForConsoleInteractive()` logic into `selectProfileWithHistory()` in `cmd/profilehelper.go`. Console, profsel, and ensureAWSProfile all call this function.

**Rationale**: profilehelper.go is already the home for shared profile selection logic. Keeping it there avoids cross-package imports and matches the existing pattern.

**Alternative considered**: Moving to a new internal package — unnecessary indirection for what is purely a UI composition function.

### 3. Drop `selectProfileInteractive()` entirely

**Decision**: Delete `selectProfileInteractive()` from profilehelper.go and replace all call sites with `selectProfileWithHistory()`.

**Rationale**: Having two similar functions invites regression. History is strictly additive UX; there is no reason to keep a history-less version.

### 4. No history migration

**Decision**: Existing `~/.local/share/bmc/console-history.json` is not migrated. Users see an empty recent list once after upgrading.

**Rationale**: History is UX convenience only. The migration complexity (reading old file, merging, deleting) is not worth the one-time inconvenience of an empty list.

## Risks / Trade-offs

- **Lost console history on upgrade** → Accepted; history repopulates quickly with normal use.
- **History key rename is a one-way change** → If a future command wants command-specific history, it adds its own key; the shared `"profile"` key remains for the shared selector.

## Migration Plan

No deployment steps required. The new history file is created automatically on first interactive use. The old `console-history.json` is left in place but no longer written to; it can be deleted manually.
