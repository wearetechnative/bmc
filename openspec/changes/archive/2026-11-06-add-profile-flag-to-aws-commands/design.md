## Context

Six AWS operation commands (`ec2connect`, `ec2ls`, `ec2`, `ec2scheduler`, `ec2stopstart`, `ecsconnect`) all delegate profile resolution to `ensureAWSProfile()` in `cmd/profilehelper.go`. That function currently reads `AWS_PROFILE` from the environment, then falls back to interactive selection.

`console` and `profsel` each have their own `-p`/`--profile` flag with custom logic. A Cobra persistent flag on root cannot coexist with these local flags of the same name, so a global persistent flag is not viable.

## Goals / Non-Goals

**Goals:**
- Allow users to pass `-p`/`--profile <name>` directly to the 6 AWS operation commands
- Centralise the flag's effect in `ensureAWSProfile()` so no per-command logic is duplicated
- Leave `console` and `profsel` completely unchanged

**Non-Goals:**
- Making `-p` a root-level persistent flag
- Changing the profile resolution logic for `console` or `profsel`
- Adding interactive fuzzy-matching or partial-name resolution for the flag value

## Decisions

### Shared variable in profilehelper.go

A package-level `globalProfile string` variable is declared in `profilehelper.go`. Each of the 6 commands registers `-p`/`--profile` pointing at this variable via `Flags().StringVarP(...)` in their `init()` function. `ensureAWSProfile()` reads `globalProfile` first in its resolution chain.

**Alternatives considered:**
- Persistent flag on root: rejected — conflicts with existing `-p` on `console` and `profsel`
- Pass profile as argument to `ensureAWSProfile()`: would require changing all 6 call sites and every future caller; shared var is simpler

### Resolution priority in ensureAWSProfile()

```
globalProfile (flag) → AWS_PROFILE (env) → interactive picker
```

If `globalProfile` is set, skip env var and interactive selection entirely, run MFA check, and return.

**Rationale:** Flag is the most explicit signal; it should win. Env var is implicit session context; interactive is the fallback.

## Risks / Trade-offs

- **MFA still runs when `-p` is used** — this is correct behaviour; the flag bypasses the picker but not security checks
- **No validation of profile name at flag-parse time** — invalid profile names are caught inside `ensureAWSProfile()` the same way env var profiles are, with a clear error message
- **Repetition in init() across 6 files** — minimal; each is a single `Flags().StringVarP` line pointing to the same var

## Migration Plan

No migration needed. Existing behaviour (env var, interactive) is preserved. The flag is purely additive.
