## Why

Bean: [.beans/bmc-1p3d--profile-flag-noopt-parsing.md](../../../.beans/bmc-1p3d--profile-flag-noopt-parsing.md)

`bmc ec2connect -p TN-Production compute2` fails with `Error: accepts at most 1 arg(s), received 2`. The `-p/--profile` flag is overloaded to mean both "use this profile" (with a value) and "force interactive selection" (bare `-p`). To make bare `-p` valid, the flag sets pflag's `NoOptDefVal`, which has the side effect that `-p` can **never** consume the next space-separated token as its value — so `-p NAME <search>` is misparsed as two positional arguments and rejected. Users must currently use the awkward `-p=NAME` equals form, which is undiscoverable and breaks all muscle-memory `-p NAME` usage.

## What Changes

- Split the overloaded profile flag into two single-purpose flags, applied uniformly across all commands that select an AWS profile:
  - `-p / --profile NAME` — always takes a value; both `-p NAME` (space) and `-p=NAME` (equals) work
  - `-P / --pick` — boolean; forces the interactive profile picker (replaces bare `-p`)
- Remove the `NoOptDefVal = " "` setting from all affected commands: `ec2`, `ec2connect`, `ec2ls`, `ec2scheduler`, `ec2stopstart`, `ecsconnect`, `console`.
- Remove the hand-rolled `args[0]` fallback workaround in `console` that compensated for the misparse.
- Add the `-P/--pick` alias to `profsel` for a consistent mental model (it already works correctly; this is consistency only).
- **BREAKING**: Bare `-p` (with no value) no longer forces interactive selection. It now errors with `flag needs an argument`. Use `-P/--pick` instead.

## Capabilities

### New Capabilities
<!-- None. -->

### Modified Capabilities
- `profile-selection`: The command-line contract for selecting a profile changes. `-p/--profile` becomes a plain value flag (space and equals forms both valid), and a new `-P/--pick` flag becomes the sole trigger for forced interactive selection. Bare `-p` no longer triggers interactive selection.
- `aws-console-access`: The `bmc console` "Force Profile Selection" requirement changes from bare `-p` to `-P/--pick`; `-p/--profile` becomes value-only.

## Impact

- **Affected code**:
  - `cmd/ec2.go`, `cmd/ec2connect.go`, `cmd/ec2ls.go`, `cmd/ec2scheduler.go`, `cmd/ec2stopstart.go`, `cmd/ecsconnect.go`, `cmd/console.go` — remove `NoOptDefVal`, add `-P/--pick` flag.
  - `cmd/profsel.go` — add `-P/--pick` alias.
  - `cmd/profilehelper.go` — `ensureAWSProfile()` replaces the `TrimSpace(globalProfile) == ""` interactive branch with an explicit `pick` boolean.
  - Consider a shared `addProfileFlags(cmd)` helper since the flag registration is currently duplicated identically across files.
- **User-facing / docs**: CHANGELOG breaking-change note; docs site (`docs/content/`) updated to describe `-p NAME` vs `-P/--pick`.
- **Dependencies**: none (cobra/pflag only).
