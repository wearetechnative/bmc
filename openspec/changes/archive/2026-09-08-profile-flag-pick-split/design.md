## Context

See proposal.md — Why. The `-p/--profile` flag currently sets pflag's `NoOptDefVal = " "` in seven commands so that a bare `-p` is a valid, value-less flag that triggers interactive selection. pflag's rule is that a flag with `NoOptDefVal` set will never bind the following space-separated token as its value — so `-p NAME <positional>` is misparsed as two positionals. `ensureAWSProfile()` (in `cmd/profilehelper.go`) distinguishes the three modes by inspecting `globalProfile`: `" "` (whitespace, from `NoOptDefVal`) means interactive, non-empty means "use it", empty means "unset". `console.go` does not use `ensureAWSProfile()` and additionally pulls `args[0]` back into the profile as a workaround for the misparse.

`profsel` is already correct — it never set `NoOptDefVal`, so its `-p/--profile` already takes a value normally and interactive is its default when `-p` is omitted.

## Goals / Non-Goals

**Goals:**
- Make `-p NAME` (space form) work alongside a positional search argument.
- Give "force interactive" its own unambiguous flag (`-P/--pick`) with uniform behavior across all profile-taking commands.
- Remove the `NoOptDefVal` misfeature and the `console.go` `args[0]` workaround it forced.
- Reduce duplication by centralizing profile-flag registration.

**Non-Goals:**
- No change to the interactive picker UI or `selectProfileInteractive` navigation behavior.
- No change to MFA handling or profile resolution logic beyond how the "interactive vs named" decision is signalled.
- No backward-compatibility shim for bare `-p` — the break is accepted and documented.

## Decisions

### Decision: Two single-purpose flags instead of one overloaded flag
`-p/--profile` becomes a plain `StringVarP` (no `NoOptDefVal`); `-P/--pick` becomes a `BoolVarP`. A flag that means exactly one thing lets pflag parse `-p NAME` unambiguously and keeps positional args intact.

- **Alternative — keep `NoOptDefVal`, document `-p=NAME`**: rejected; trains users around a footgun and every muscle-memory `-p NAME` still fails.
- **Alternative — relax `Args` to `MaximumNArgs(2)` and disambiguate in code**: rejected; heuristics on which positional is the profile are fragile and order-dependent.
- **Alternative — `-p=` empty-string sentinel for interactive**: rejected in favor of an explicit, discoverable `-P/--pick` flag that shows up in `--help`.

### Decision: `pick` boolean replaces the whitespace sentinel in `ensureAWSProfile()`
`ensureAWSProfile` gains an explicit signal (a package-level `globalPick` bool, or a parameter) that replaces the `strings.TrimSpace(globalProfile) == ""` interactive branch. Resolution order becomes: `pick` → interactive; else non-empty `profile` → use it; else `AWS_PROFILE` → use it; else interactive.

### Decision: Centralize flag registration in a shared helper
Introduce `addProfileFlags(cmd *cobra.Command)` (in `cmd/profilehelper.go`) that registers both `-p/--profile` and `-P/--pick` against the shared `globalProfile`/`globalPick` vars. Each command calls it in its `init()` instead of duplicating the two `Flags()` lines. This guarantees the seven commands stay consistent and prevents a future command from reintroducing `NoOptDefVal`.

### Decision: `console.go` adopts the shared flags and drops its workaround
`console.go` uses its own `consoleProfile` var and `args[0]` fallback. It moves to the shared `-p/--profile` + `-P/--pick` flags and deletes the `args[0]`-into-profile block, mapping its existing three-way switch onto the `profile` / `pick` / env logic.

### Decision: `profsel` gains `-P/--pick` as an alias only
`profsel` keeps its current `-p/--profile` pre-select semantics and its interactive default. Adding `-P/--pick` gives it the same vocabulary as the other commands (forcing interactive even if a value could otherwise be inferred), for a uniform mental model.

## Risks / Trade-offs

- **Breaking change: bare `-p` now errors** → Mitigation: CHANGELOG breaking-change entry and docs update showing `-P/--pick`; the error message pflag emits (`flag needs an argument: 'p'`) is self-explanatory.
- **console.go has a bespoke flow** → Mitigation: migrate it carefully and verify its three modes (named, pick, env/default) against the new flags; keep its `--service`/`--watch` flags untouched.
- **Muscle-memory `-p` (bare) users** → Mitigation: the replacement `-P` is a one-character relearn and is surfaced in `--help`.

## Migration Plan

1. Add `globalPick` var and `addProfileFlags` helper in `cmd/profilehelper.go`.
2. Update `ensureAWSProfile()` to branch on `globalPick`.
3. Replace per-command flag registration with `addProfileFlags(cmd)` in the seven commands; delete each `NoOptDefVal` line.
4. Migrate `console.go` to the shared flags and remove its `args[0]` workaround.
5. Add `-P/--pick` to `profsel`.
6. Update CHANGELOG (`## NEXT VERSION`, breaking note) and docs (`docs/content/`).

Rollback: revert the change; no persisted state or config migration is involved.
