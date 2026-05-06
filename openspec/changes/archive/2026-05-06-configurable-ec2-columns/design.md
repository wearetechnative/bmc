## Context

EC2 instance tables are rendered in four commands: `ec2ls` (display), `ec2connect` / `ec2stopstart` / `ec2scheduler` (interactive selection via shared `selectInstanceID`), and `ec2find` (multi-profile display). All four build `[][]string` rows with an identical hardcoded column list. `selectInstanceID` returns `row[0]` assuming InstanceId is always first. The config system uses TOML and already has an `EC2Config` struct.

## Goals / Non-Goals

**Goals:**
- Single row-building path used by all EC2 table commands
- User-configurable column list and order via `config.toml`
- Default column order puts Name second (after InstanceId)
- Unknown column names in config silently render as `"n/a"`
- `selectInstanceID` moved to `cmd/instancehelper.go`

**Non-Goals:**
- Per-command column lists (one global `[ec2] columns` setting)
- Arbitrary tag columns (only the fields on the `Instance` struct are supported)
- Column width or alignment configuration

## Decisions

### 1. Row builder lives in `internal/awsops`

`InstanceFieldValue(inst Instance, col string) string` and `InstanceRows(instances []Instance, cols []string) [][]string` are placed in `internal/awsops/ec2.go`. This keeps the mapping close to the data type it maps.

**Alternative considered**: helper in `cmd/` package. Rejected because `cmd` is the CLI layer; the mapping is data logic.

### 2. Column names are PascalCase strings matching table headers

Config uses `columns = ["InstanceId", "Name", "PrivateIP", ...]`. The same strings are used as table column headers. This makes the connection obvious and avoids a translation layer.

**Alternative considered**: snake_case (`instance_id`, `private_ip`). Rejected to avoid a mapping between config names and display names.

### 3. `selectInstanceID` resolves InstanceId by column index

After retrieving the selected row from `ui.SelectFromTable`, the function finds the index of `"InstanceId"` in the configured column list and reads `row[idx]`. If `"InstanceId"` is not in the configured columns, it falls back to `row[0]`.

**Alternative considered**: always force `"InstanceId"` as the first column in the selection path. Rejected because it would silently diverge from the user's configured order.

### 4. `ec2find` always appends `"Profile"`

`ec2find` uses the shared builder with the configured columns, then appends `"Profile"` if it is not already present. This ensures cross-profile results always show which profile an instance belongs to.

### 5. Default column order

`["InstanceId", "Name", "PrivateIP", "PublicIP", "State", "Hibernate", "Scheduler"]` — Name moved from position 6 to position 2 as requested.

## Risks / Trade-offs

- [InstanceId removed from config] → `selectInstanceID` cannot return an ID. Mitigation: fall back to `row[0]`; document that InstanceId should not be removed from the column list when using interactive commands.
- [Config file absent] → `Defaults()` returns the new default column list; behaviour unchanged for users without a config file.

## Migration Plan

No migration needed. Config file addition is opt-in; existing users without a `columns` entry get the new default order automatically.
