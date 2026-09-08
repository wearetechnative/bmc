## 1. Shared flag plumbing

- [x] 1.1 Add a package-level `globalPick` bool in `cmd/profilehelper.go` alongside `globalProfile`, and verify the package still compiles (`go build ./...`)
- [x] 1.2 Add an `addProfileFlags(cmd *cobra.Command)` helper in `cmd/profilehelper.go` that registers `-p/--profile` (StringVarP → `globalProfile`, no `NoOptDefVal`) and `-P/--pick` (BoolVarP → `globalPick`); verify with `go build ./...`
- [x] 1.3 Update `ensureAWSProfile()` to branch on `globalPick` for interactive selection instead of `strings.TrimSpace(globalProfile) == ""`; resolution order = pick → named profile → `AWS_PROFILE` → interactive. Verify `go build ./...` and that the now-unused `strings` import is handled

## 2. Migrate profile-taking commands

- [x] 2.1 Replace the per-command `--profile` registration and delete the `NoOptDefVal` line in `cmd/ec2.go`, `cmd/ec2ls.go`, `cmd/ec2scheduler.go`, `cmd/ec2stopstart.go`, `cmd/ecsconnect.go` by calling `addProfileFlags(cmd)`; verify `go build ./...` and `grep -rn NoOptDefVal cmd/` returns nothing for these files
- [x] 2.2 Migrate `cmd/ec2connect.go` to `addProfileFlags(cmd)`, remove its `NoOptDefVal` line, and verify `bmc ec2connect -p <profile> <search>` parses without an "accepts at most 1 arg(s)" error (build then run against a real/help invocation)
- [x] 2.3 Migrate `cmd/console.go` to the shared flags, remove its `NoOptDefVal` line AND the `args[0]`-into-profile workaround block, mapping its three-way switch onto profile/pick/env; verify `go build ./...` and that `bmc console -p <profile>` still opens the console
- [x] 2.4 Add `-P/--pick` to `cmd/profsel.go` (keeping existing `-p/--profile` pre-select behavior) and verify `bmc profsel --pick` forces the interactive picker

## 3. Verification

- [x] 3.1 Confirm `grep -rn NoOptDefVal cmd/` returns no matches across the whole `cmd/` package
- [x] 3.2 Verify acceptance scenarios end-to-end: `bmc ec2connect -p <profile> compute2` (space form works), `bmc ec2connect -p=<profile> compute2` (equals form works), `bmc ec2connect --pick compute2` (picker then search), and bare `bmc ec2connect -p` errors with "flag needs an argument"

## 4. Docs & changelog

- [x] 4.1 Add a `## NEXT VERSION` entry to `CHANGELOG.md` under `### Changed` describing the `-p/--profile` fix and the **BREAKING** bare-`-p` change (with `-P/--pick` as the replacement)
- [x] 4.2 Update the docs site (`docs/content/`) pages that document profile selection to describe `-p NAME` vs `-P/--pick`; verify the referenced flags match the implemented behavior
