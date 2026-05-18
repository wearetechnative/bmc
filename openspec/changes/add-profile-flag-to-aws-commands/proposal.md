## Why

Users often work with multiple AWS profiles and must currently set `AWS_PROFILE` as an environment variable before running commands like `bmc ec2connect` or `bmc ec2ls`. A `-p`/`--profile` flag directly on these commands avoids the extra shell step and makes one-off cross-profile operations more ergonomic.

## What Changes

- Add a shared `globalProfile` variable to `cmd/profilehelper.go`
- Register `-p`/`--profile` as a local flag on each of the 6 commands that call `ensureAWSProfile()`: `ec2connect`, `ec2ls`, `ec2`, `ec2scheduler`, `ec2stopstart`, `ecsconnect`
- `ensureAWSProfile()` reads `globalProfile` first, before falling back to `AWS_PROFILE` env var, then interactive selection
- `console` and `profsel` retain their existing `-p` flag implementations unchanged

## Capabilities

### New Capabilities

- `aws-command-profile-flag`: `-p`/`--profile` flag on all AWS operation commands, backed by a shared variable in `profilehelper.go`, allowing users to specify an AWS profile directly on the command line without setting `AWS_PROFILE`

### Modified Capabilities

- `profile-selection`: `ensureAWSProfile()` gains a new resolution step — flag value takes priority over env var and interactive selection

## Impact

- `cmd/profilehelper.go`: new `globalProfile` var, `ensureAWSProfile()` updated
- `cmd/ec2connect.go`, `cmd/ec2ls.go`, `cmd/ec2.go`, `cmd/ec2scheduler.go`, `cmd/ec2stopstart.go`, `cmd/ecsconnect.go`: each registers the flag locally
- No impact on `console` or `profsel`
- No breaking changes
