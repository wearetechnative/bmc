## Why

`bmc ec2ls` and `bmc ec2find` currently only produce human-readable table output, making them unusable in scripts and automation pipelines. Adding `--json` output enables composability with tools like `jq`, shell scripts, and CI workflows.

Bean: [bmc-wko6](../../../.beans/bmc-wko6--json-output.md)

## What Changes

- `bmc ec2ls --json` outputs all EC2 instances as a JSON array to stdout, ignoring the `columns` config
- `bmc ec2find <search> --json` outputs matching instances as a JSON array; group selection remains interactive via bubbletea TUI (rendered on `/dev/tty`, separate from stdout)
- JSON keys follow AWS CLI PascalCase convention: `InstanceId`, `Name`, `PrivateIpAddress`, `PublicIpAddress`, `State`, `Hibernate`, `Scheduler`
- `ec2find` JSON output always includes the `Profile` field (which identifies the source AWS profile)
- All fields are always present in JSON output regardless of the `ec2.columns` config setting

## Capabilities

### New Capabilities

- `ec2-json-output`: `--json` flag for `ec2ls` and `ec2find` commands producing machine-readable JSON output

### Modified Capabilities

_(none — existing table output behaviour is unchanged)_

## Impact

- `cmd/ec2ls.go`: add `--json` flag, branch output path
- `cmd/ec2find.go`: add `--json` flag, branch output path
- `internal/awsops/ec2.go`: add `json:` struct tags to `Instance` with AWS PascalCase field names
- No breaking changes; existing behaviour unchanged when `--json` is not passed
