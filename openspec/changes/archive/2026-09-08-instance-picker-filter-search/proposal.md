## Why

Bean: [.beans/bmc-p6m8--instance-table-filter-search.md](../../../.beans/bmc-p6m8--instance-table-filter-search.md)

The interactive EC2 instance picker offers no type-to-filter/search — users can only scroll with the arrow keys to locate an instance. The profile picker (`ui.Choose`, built on `bubbles/list`) already filters as you type, so the experience is inconsistent. The instance picker uses `ui.SelectFromTable` (built on `bubbles/table`), which has no filtering. On accounts with many instances this makes the picker slow and frustrating to use.

## What Changes

- Add real-time type-to-filter to the interactive instance table selector (`ui.SelectFromTable` / `selectInstanceID`), so the user can narrow the visible rows by typing.
- The filter SHALL match against the same fields the positional `[search]` fragment already uses (`InstanceID`, `Name`, `PrivateIP`, `PublicIP`), case-insensitively, keeping behavior consistent between the two entry points.
- Because the selector is shared, the improvement applies uniformly to all commands that use it: `ec2connect`, `ec2`, `ec2stopstart`, `ec2scheduler`.
- Preserve existing selection, cancel (Esc/Ctrl+C), and scrolling behavior; filtering is additive.

## Capabilities

### New Capabilities
- `instance-picker-search`: Defines that the interactive instance table selector supports real-time type-to-filter over instance identifying fields, shared across all commands that present the instance picker.

### Modified Capabilities
<!-- None. -->

## Impact

- **Affected code**:
  - `internal/ui/table.go` — add a filter input to the `tableModel` (or provide a filterable variant), matching rows in real time.
  - `cmd/instancehelper.go` — `selectInstanceID` passes the filterable fields/values through so the widget can match against `InstanceID`/`Name`/`PrivateIP`/`PublicIP`, not just rendered columns.
  - Beneficiaries (no call-site logic change expected): `cmd/ec2connect.go`, `cmd/ec2.go`, `cmd/ec2stopstart.go`, `cmd/ec2scheduler.go`.
- **Dependencies**: `charmbracelet/bubbles` (already vendored — `table`, `textinput`, `list`). No new dependency.
- **Docs**: docs site (`docs/content/`) note that the instance picker is filterable.
