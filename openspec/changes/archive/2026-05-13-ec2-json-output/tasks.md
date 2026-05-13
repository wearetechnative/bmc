## 1. Instance struct JSON tags

- [x] 1.1 Add `json:` struct tags to `Instance` in `internal/awsops/ec2.go` using AWS PascalCase names (`InstanceId`, `Name`, `PrivateIpAddress`, `PublicIpAddress`, `State`, `Hibernate`, `Scheduler`, `Profile`)

## 2. ec2ls JSON flag

- [x] 2.1 Add `ec2lsJSON bool` variable and register `--json` flag on `ec2lsCmd` in `cmd/ec2ls.go`
- [x] 2.2 In `runEC2ls`, when `--json` is set, marshal instances to JSON and print to stdout instead of calling `ui.PrintTable`

## 3. ec2find JSON flag

- [x] 3.1 Add `ec2findJSON bool` variable and register `--json` flag on `ec2findCmd` in `cmd/ec2find.go`
- [x] 3.2 In `runEC2Find`, when `--json` is set, marshal matched instances to JSON and print to stdout instead of calling `ui.ShowTable`

## 4. Update documentation

- [x] 4.1 Add `--json` flag to `ec2ls` and `ec2find` entries in `README.md`
- [x] 4.2 Update CHANGELOG under `## NEXT VERSION`
