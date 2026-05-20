## Why

The original bash-era `ec2connect.sh` supported specifying an SSH identity file via a flag. This functionality was lost during the Go rewrite when the `-i` flag was repurposed for instance ID selection. Users who work with non-default SSH keys (e.g., AWS-generated key pairs, per-environment keys) currently have no way to pass an identity file through `bmc ec2connect`.

## What Changes

- Add `-k`/`--key` flag to `ec2connect` for specifying an SSH identity file path
- When `-k` is provided, SSH is automatically selected as the connection method (no method picker prompt shown)
- The key path is passed directly to `ssh -i <path>` — no file existence validation

## Capabilities

### New Capabilities

- `ec2connect-ssh-key`: Support for specifying an SSH identity file when connecting to EC2 instances via SSH

### Modified Capabilities

(none)

## Impact

- `cmd/ec2connect.go`: add flag declaration and pass key to `connectSSH`
- `docs/content/commands/ec2connect.md`: document the new flag
