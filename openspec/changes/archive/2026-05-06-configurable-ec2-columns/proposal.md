## Why

The EC2 instance table is built three times with identical hardcoded column lists (`ec2ls`, `selectInstanceID`, `ec2find`), making it impossible to change the display without editing multiple files. Users also have no control over which columns are shown or in what order, and the default column order buries the Name column at the end.

## What Changes

- Add `columns` field to `[ec2]` section of `~/.config/bmc/config.toml` to configure which columns appear in EC2 instance tables and in what order
- Default column order changed: `["InstanceId", "Name", "PrivateIP", "PublicIP", "State", "Hibernate", "Scheduler"]` (Name moved from second-to-last to second position)
- New shared `InstanceRows()` and `InstanceFieldValue()` functions in `internal/awsops` replace three duplicate row-building blocks
- `selectInstanceID()` moved from `cmd/ec2stopstart.go` to new `cmd/instancehelper.go`
- `selectInstanceID()` resolves InstanceId by column index rather than hardcoded `row[0]`
- `ec2find` uses the shared builder and always appends `"Profile"` to the configured column list
- Unknown column names in config silently produce `"n/a"` values
- Documentation updated to list all available column names
- `ec2ls` outputs a formatted table to stdout using `lipgloss.NewTable()` — bordered, aligned, bold header; pipeable and scrollable via terminal scrollback
- Interactive selection tables (ec2connect, ec2stopstart, ec2scheduler) show a footer with row count and key hints, and resize to fit the terminal

## Capabilities

### New Capabilities
- `ec2-instance-columns`: Configurable EC2 instance table columns — which fields are shown, in what order, via `config.toml`

### Modified Capabilities

## Impact

- `internal/awsops/ec2.go`: New exported functions `InstanceFieldValue` and `InstanceRows`
- `internal/config/config.go`: `EC2Config` gains `Columns []string` field with defaults
- `cmd/ec2ls.go`: Uses shared builder instead of inline row construction
- `cmd/ec2stopstart.go`: `selectInstanceID` moved out; uses shared builder
- `cmd/ec2connect.go`: No logic change; benefits from shared builder via `selectInstanceID`
- `cmd/ec2scheduler.go`: No logic change; benefits from shared builder via `selectInstanceID`
- `cmd/ec2find.go`: Uses shared builder; always appends `"Profile"` column
- `cmd/instancehelper.go`: New file containing `selectInstanceID`
- `internal/ui/table.go`: `PrintTable` renders a lipgloss-bordered table to stdout; `SelectFromTable` gains footer + terminal-height-aware sizing via `/dev/tty`
- `cmd/ec2ls.go`: Uses `ui.PrintTable` instead of interactive table
- Bean: bmc-599m
