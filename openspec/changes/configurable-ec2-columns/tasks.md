## 1. Shared row builder in awsops

- [x] 1.1 Add `InstanceFieldValue(inst Instance, col string) string` to `internal/awsops/ec2.go` — switch on col, return field value or `"n/a"`
- [x] 1.2 Add `InstanceRows(instances []Instance, cols []string) [][]string` to `internal/awsops/ec2.go` — iterate instances and cols using `InstanceFieldValue`

## 2. Config: add Columns field

- [x] 2.1 Add `Columns []string` to `EC2Config` in `internal/config/config.go` with TOML tag `toml:"columns"`
- [x] 2.2 Update `Defaults()` to set `Columns` to `["InstanceId", "Name", "PrivateIP", "PublicIP", "State", "Hibernate", "Scheduler"]`

## 3. Move and update selectInstanceID

- [x] 3.1 Create `cmd/instancehelper.go` with `selectInstanceID(instances []awsops.Instance, cols []string) (string, error)`
- [x] 3.2 Implement column-index-based InstanceId lookup: find index of `"InstanceId"` in cols, fall back to index 0
- [x] 3.3 Remove `selectInstanceID` from `cmd/ec2stopstart.go`

## 4. Update commands to use shared builder

- [x] 4.1 Update `cmd/ec2ls.go` — load config, use `awsops.InstanceRows` with `cfg.EC2.Columns`
- [x] 4.2 Update `cmd/ec2stopstart.go` — pass `cfg.EC2.Columns` to `selectInstanceID`; load config where needed
- [x] 4.3 Update `cmd/ec2connect.go` — pass `cfg.EC2.Columns` to `selectInstanceID` (config already loaded)
- [x] 4.4 Update `cmd/ec2scheduler.go` — load config, pass `cfg.EC2.Columns` to `selectInstanceID`
- [x] 4.5 Update `cmd/ec2find.go` — use `awsops.InstanceRows` with config columns + append `"Profile"` if not already present

## 5. Documentation

- [x] 5.1 Update `README.md` to document `[ec2] columns` config option and list all available column names

## 6. EC2 table display improvements

- [x] 6.1 Export `PrintTable` in `internal/ui/table.go` for plain stdout output
- [x] 6.2 Update `cmd/ec2ls.go` to use `ui.PrintTable` (plain, pipeable)
- [x] 6.3 Add footer to interactive table: row count + key hints (`↑↓ scroll  enter select  q cancel`)
- [x] 6.4 Make interactive table terminal-height-aware: initial height via `/dev/tty` `term.GetSize`, update on `WindowSizeMsg`
