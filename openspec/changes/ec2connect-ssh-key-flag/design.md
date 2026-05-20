## Context

`ec2connect` currently supports `-i` (instance ID), `-u` (SSH user), and `-p` (AWS profile). When `-u` is set, the connection method picker is skipped and SSH is used directly. The bash-era script had a `-i` flag for SSH key that was lost when `-i` was repurposed for instance ID in the Go rewrite.

The `connectSSH` function builds an `ssh` invocation via `syscall.Exec`, giving direct control over the arguments passed to the system `ssh` binary.

## Goals / Non-Goals

**Goals:**
- Add `-k`/`--key` string flag to `ec2connect`
- Auto-select SSH when `-k` is set (same pattern as `-u`)
- Pass the key path as `-i <path>` to the `ssh` binary

**Non-Goals:**
- File existence validation (ssh provides its own error)
- Support for multiple keys
- Integration with `~/.ssh/config` management

## Decisions

### Flag name: `-k`/`--key`

`-i` is already used for instance ID. `-k` is the natural short form for "key" and does not conflict with any existing flags. `--key` is self-documenting.

### No file validation

The `ssh` binary produces a clear error when an identity file is missing or unreadable. Adding a redundant check in bmc adds complexity with no user benefit. This matches the approach taken for other flags in bmc (e.g., no validation of AWS profile existence before passing to AWS CLI).

### Auto-select SSH when `-k` is set

Consistent with the existing `-u` behaviour: if any SSH-specific flag is provided, the connection method picker is skipped. Providing a key and then selecting SSM would be nonsensical.

## Risks / Trade-offs

- **Paths with spaces**: `syscall.Exec` passes args as a slice, so quoting is not an issue — the path is a single element. No risk here.
- **No validation**: If the user provides a wrong path, the error comes from `ssh`, not `bmc`. This is acceptable and consistent with the rest of the codebase.

## Open Questions

(none)
