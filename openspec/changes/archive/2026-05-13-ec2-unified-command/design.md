## Context

The project already has focused EC2 commands (`ec2connect`, `ec2stopstart`, `ec2scheduler`), each requiring the user to repeat instance selection. The implementation logic for instance selection (`selectInstanceID`), connection (`connectSSH`, `connectSSM`), state changes (`startInstance`, `stopInstance`), and scheduler toggling (`awsops.ToggleSchedulerTag`) all live inside the `cmd` package and are accessible within the same package.

`ec2ls` is intentionally display-only and will not be modified.

## Goals / Non-Goals

**Goals:**
- New `bmc ec2 [search]` command: select one instance, then pick an action
- Reuse all existing helpers within the `cmd` package without refactoring them
- Consistent search/filter UX as `ec2connect [search]`
- Register command in `root.go`

**Non-Goals:**
- Modifying `ec2ls`, `ec2connect`, `ec2stopstart`, `ec2scheduler`, `ec2find`
- Adding new AWS operations not already present
- SSH key management or user selection shortcuts in the action menu

## Decisions

### Decision 1: New file `cmd/ec2.go`, no refactoring

All helpers (`connectSSH`, `connectSSM`, `startInstance`, `stopInstance`, `selectInstanceID`, `waitForState`) are package-level functions in the `cmd` package. `ec2.go` lives in the same package and can call them directly — no exports or interface changes needed.

**Alternative considered:** Extract shared logic into an `internal/ec2ops` package. Rejected because it requires significant refactoring of existing files for no functional benefit.

### Decision 2: Action menu as second step after instance selection

After instance selection, show a simple list:
- Connect SSH
- Connect SSM
- Start / Stop
- Toggle scheduler

Navigation: ESC at action menu returns to instance picker (or exits if only one instance matched). The existing `ui.Choose` already handles ESC → returns empty string.

**Alternative considered:** Flags like `bmc ec2 connect`, `bmc ec2 start`. Rejected because it adds subcommand complexity and doesn't match the "select then act" UX goal.

### Decision 3: Search argument identical to `ec2connect`

Same case-insensitive substring match on `InstanceID + Name + PrivateIP + PublicIP`. Single match skips picker. Zero matches → error. This is copy-paste logic from `ec2connect.go`.

### Decision 4: `Start / Stop` label adapts to current state

The action menu item reads `"Start instance"` or `"Stop instance"` depending on the instance's current state, fetched before showing the menu (reuse `awsops.GetInstanceState`). If the instance is in a transitional state (pending, stopping, etc.), the action is omitted with a message.

## Risks / Trade-offs

- **Extra API call**: `ec2.go` calls `awsops.GetInstanceState` to label the start/stop action. This is an extra DescribeInstances call. Acceptable — same as what `ec2stopstart` does.
- **Action functions are stateful**: `connectSSH`/`connectSSM` use `syscall.Exec` (process replacement). `ec2.go` must handle this like `ec2connect.go` — no code after calling them.
- **Duplicate filter logic**: The search fragment logic is duplicated from `ec2connect.go`. Acceptable trade-off to avoid premature abstraction (only two callers).
