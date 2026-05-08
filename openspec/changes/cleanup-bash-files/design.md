## Context

BMC was originally a bash toolbox (`_bmclib.sh`) with supporting scripts. A complete Go rewrite replaced it. The bash files were kept during the transition but are now fully superseded. Additionally, `bmc-go` (a compiled binary) was accidentally committed to git.

## Goals / Non-Goals

**Goals:**
- Remove three bash-era source files from the repository
- Remove the accidentally tracked `bmc-go` binary from git and prevent future accidental commits via `.gitignore`

**Non-Goals:**
- Any changes to Go source code
- Any changes to user-facing behavior
- Archiving files (git history preserves them if ever needed)

## Decisions

**Delete, not archive**
The files will be deleted outright rather than moved to an `archive/` folder. Rationale: git history already serves as the archive. An `archive/` folder in the working tree adds noise without value.

**`.gitignore` for bmc-go**
Add `/bmc-go` to `.gitignore` alongside the existing `/bmc` entry. Both are build artifacts that should never be committed.

## Risks / Trade-offs

- [Low] Someone may have sourced `_bmclib.sh` locally → Mitigation: the Go binary has been the documented install path since the rewrite; no users of the bash scripts are expected.
- [None] No rollback risk — file deletion is trivially reversible via `git checkout` or `git restore`.
