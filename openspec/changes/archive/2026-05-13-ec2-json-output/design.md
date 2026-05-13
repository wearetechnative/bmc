## Context

`bmc ec2ls` and `bmc ec2find` render output through `ui.PrintTable` / `ui.ShowTable`, which write formatted text to stdout. There is no machine-readable output path. The `profsel` command already implements a `--json` flag as a precedent pattern in this codebase.

The `Instance` struct in `internal/awsops/ec2.go` holds all relevant fields and is already used by both commands. It currently lacks `json:` struct tags.

`ui.Choose` (bubbletea group picker in `ec2find`) already opens `/dev/tty` for both input and output, meaning it is completely independent of stdout. JSON can therefore be piped from `ec2find` without interfering with interactive group selection.

## Goals / Non-Goals

**Goals:**
- Add `--json` flag to `ec2ls` producing a complete JSON array of all instances
- Add `--json` flag to `ec2find` producing a complete JSON array of matching instances
- Use AWS CLI PascalCase key names consistent with AWS tooling conventions
- Always output all fields regardless of the `ec2.columns` config

**Non-Goals:**
- JSON output for any other command (`ec2connect`, `ec2stopstart`, `ec2scheduler`)
- Filtering or selecting specific JSON fields via flags
- Adding a `--group` flag to `ec2find` (bubbletea via `/dev/tty` is sufficient for interactive use alongside piped output)

## Decisions

### JSON key naming: AWS CLI PascalCase

**Decision**: Use `InstanceId`, `Name`, `PrivateIpAddress`, `PublicIpAddress`, `State`, `Hibernate`, `Scheduler`, `Profile`.

**Rationale**: AWS CLI uses PascalCase for all EC2 field names. Users scripting with `bmc` alongside `aws ec2 describe-instances` will find consistent key names. Go's default JSON marshaling (camelCase) would be unfamiliar in this context.

**Alternative considered**: snake_case (`instance_id`, `private_ip_address`). Rejected — not consistent with AWS tooling.

### Always output all fields, ignore `columns` config

**Decision**: JSON output always includes every field on the `Instance` struct, regardless of `ec2.columns` in config.

**Rationale**: `columns` is a display preference for TUI tables. Scripts consuming JSON should receive complete data and filter with `jq`. Silently omitting fields based on display config would make JSON output unpredictable.

### Struct tags on `Instance`, not manual marshaling

**Decision**: Add `json:` tags to `awsops.Instance` and use `json.Marshal` directly.

**Rationale**: Keeps serialization declarative and avoids a parallel map-building step. The struct already contains exactly the fields needed.

### `ec2find` group selection stays interactive

**Decision**: No `--group` flag. Bubbletea TUI for group selection remains even in `--json` mode.

**Rationale**: `ui.Choose` renders via `/dev/tty`, fully separate from stdout. A user can run `bmc ec2find web --json | jq` and still interact with the group picker in the terminal. A `--group` flag adds complexity that is not needed for this use case.

## Risks / Trade-offs

- **`Profile` field in `ec2ls` JSON**: `ec2ls` operates on a single profile, so `Profile` will always be an empty string in its JSON output. This is consistent (all fields always present) but slightly redundant. → Acceptable; consumers can ignore empty fields.
- **`json:` tags are additive**: Adding struct tags to `Instance` is non-breaking; no existing code serializes this struct to JSON today.

## Migration Plan

No migration required. `--json` is an opt-in flag. All existing behaviour is unchanged when the flag is absent.
