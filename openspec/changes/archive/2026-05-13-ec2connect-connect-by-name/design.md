## Context

`ec2connect` currently supports two instance selection modes: providing an exact instance ID via `-i`, or choosing from a full interactive table. Users frequently know a partial name fragment (e.g. `nixhost` for `ec2-nixhost-prod-992382728492`) but must either look up the full ID or scroll through a long list.

The existing codebase already provides all primitives needed: `awsops.ListInstances`, `selectInstanceID` (interactive picker), and a substring filter pattern from `ec2find`.

## Goals / Non-Goals

**Goals:**
- Allow `bmc ec2connect <fragment>` to filter instances before selection
- Zero changes to any file outside `cmd/ec2connect.go`
- No new dependencies

**Non-Goals:**
- Regex or glob matching (plain case-insensitive substring is sufficient)
- Filtering by state (all states shown, existing start-stopped mechanism handles it)
- Multi-word / AND-style queries

## Decisions

### Positional argument over a named flag

**Decision**: Accept the filter as `args[0]` (positional), not as `--name`.

**Rationale**: Positional args require less typing and read naturally in scripts (`bmc ec2connect nixhost`). The command already has `-i`/`--instance` for the fully-qualified ID case; adding `--name` alongside would be redundant and harder to discover.

### Broad search scope

**Decision**: Match against `InstanceID + Name + PrivateIP + PublicIP` (case-insensitive `strings.Contains`), identical to `ec2find`.

**Rationale**: Users sometimes know a partial IP or instance ID prefix rather than a name. Reusing the exact `ec2find` pattern keeps behaviour predictable and avoids a separate filter function.

### `-i` wins over positional arg with a warning

**Decision**: If both `-i <id>` and a positional argument are supplied, `-i` takes precedence. A warning is printed to stderr.

**Rationale**: `-i` is explicit and unambiguous. Silently ignoring the positional arg would confuse users; an error would be too strict for scripting. A warning is the middle ground.

### All instance states shown

**Decision**: No state filter is applied to the filtered list.

**Rationale**: The existing stop→start confirmation mechanism already handles stopped instances. Filtering here would hide valid targets and require duplicating state-check logic.

## Risks / Trade-offs

- **Ambiguous fragment matches many instances** → User sees the picker, same UX as no filter. No risk.
- **No match found** → Clear error returned. Consistent with other bmc error handling.
- **Fragment coincidentally matches instance ID of a different instance than intended** → Broad search is a feature, not a bug. User can be more specific or use the picker.
