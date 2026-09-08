---
# bmc-1p3d
title: profile-flag-noopt-parsing
status: completed
type: bug
priority: high
created_at: 2026-09-08T05:16:06Z
updated_at: 2026-09-08T06:30:00Z
openspec-link: openspec/changes/archive/2026-09-08-profile-flag-pick-split/proposal.md
---

`bmc ec2connect -p TN-Production compute2` fails with `Error: accepts at most 1 arg(s), received 2`.

## Root cause

The `-p/--profile` flag sets `NoOptDefVal = " "` (to allow bare `-p` to force the interactive picker). A side effect of `NoOptDefVal` in pflag is that `-p` NEVER consumes the next space-separated token as its value — the only way to pass a value becomes the `=` form (`-p=TN-Production`). So `-p TN-Production compute2` is parsed as two positional args, which `MaximumNArgs(1)` rejects.

Affected commands (all set `NoOptDefVal = " "`):
ec2, ec2connect, ec2ls, ec2scheduler, ec2stopstart, ecsconnect, console.
`console.go` also carries a hand-rolled `args[0]` fallback workaround that should be removed.
`profsel` is NOT affected (no NoOptDefVal; interactive is already its default).

## Settled design

Split the overloaded flag into two single-purpose flags, applied uniformly:

- `-p / --profile NAME` -> always takes a value (space AND equals both work)
- `-P / --pick`         -> bool, force interactive picker (replaces bare `-p`)
- (neither)             -> AWS_PROFILE if set, else interactive

### Change surface
- Remove `NoOptDefVal` from the 7 affected commands
- Add `-P/--pick` bool flag everywhere (prefer a shared `addProfileFlags(cmd)` helper — flag is defined identically in every file)
- Update `ensureAWSProfile()`: replace the `TrimSpace(globalProfile) == ""` interactive branch with a `pick` bool
- Remove the `args[0]` fallback workaround in `console.go`
- Add `-P/--pick` alias to `profsel` for a uniform mental model (it already works, this is consistency only)

### Breaking change
Bare `-p` no longer forces interactive selection — use `-P/--pick`. Bare `-p` will now error (`flag needs an argument`). Note in CHANGELOG + docs.

### Acceptance
- `bmc ec2connect -p TN-Production compute2` connects using profile TN-Production, searching "compute2"
- `bmc ec2connect --pick compute2` opens the picker, then searches "compute2"
- `-p=NAME` continues to work
