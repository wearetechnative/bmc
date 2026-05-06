## Context

BMC is a bash-based AWS toolbox used daily for AWS profile selection, EC2/ECS operations, and MFA session management. It currently depends on: gum (TUI), jq (JSON), awk (text), jsonify-aws-dotfiles (config parsing), assumego/Granted (console), aws-mfa/broamski (STS session tokens). These dependencies make installation complex and cross-platform behaviour inconsistent (macOS `date -j` vs Linux `date -d`, zsh vs bash sourcing quirks).

The rewrite produces a single Go binary. AWS SDK Go v2 replaces all AWS CLI calls. Charmbracelet's bubbletea+bubbles stack replaces gum. AWS config files are parsed natively. Only `ssh`, `aws` CLI, and `session-manager-plugin` remain as optional runtime dependencies — and only for the two interactive session commands.

## Goals / Non-Goals

**Goals:**
- Single binary, zero mandatory runtime deps for all non-session commands
- Full command parity with the bash implementation
- Actionable error messages with per-platform install instructions (Nix × 3, apt, brew)
- `bmc doctor` for proactive system health checks
- Shell wrapper pattern for `profsel` to preserve `export AWS_PROFILE` UX
- Distribution via GoReleaser + GitHub releases + Homebrew tap + Nix flake
- Config migrated to TOML; `~/.aws/credentials` write format stays aws-mfa compatible

**Non-Goals:**
- Replacing session-manager-plugin (WebSocket protocol is out of scope)
- Replacing `ssh` binary
- Automatic config.env → config.toml migration (user migrates manually; doctor warns)
- `tgselect.sh` (Toggl script) is out of scope

## Decisions

### D1: Go + Cobra for CLI framework
**Decision**: Use [cobra](https://github.com/spf13/cobra) for command dispatch.
**Rationale**: Industry standard for Go CLIs. Provides built-in shell completion (replaces `gencompletions`), help generation, flag parsing. Eliminates the manual `CMDS[]`/`DESC[]` bash array pattern.
**Alternatives**: `urfave/cli` — less idiomatic, smaller ecosystem.

### D2: bubbletea + bubbles for TUI (not huh)
**Decision**: Use [bubbletea](https://github.com/charmbracelet/bubbletea) + [bubbles](https://github.com/charmbracelet/bubbles) for all interactive UI.
**Rationale**: bubbles/list has built-in fuzzy filtering (replaces both `gum choose` and `gum filter`). bubbles/table replaces `gum table`. bubbles/spinner replaces `gum spin`. bubbles/textinput replaces `gum input`. Same Charmbracelet ecosystem as gum — similar UX, pure Go.
**Alternatives**: `huh` (Charmbracelet forms) — good for forms but not for filterable list selection of AWS resources.

### D3: AWS SDK Go v2 for all API calls
**Decision**: All EC2, ECS, STS calls via [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2).
**Rationale**: Eliminates aws CLI dependency for core operations. Typed responses eliminate jq parsing. Concurrent multi-profile queries (ec2find) become straightforward with goroutines.
**Alternatives**: Shell out to `aws` CLI — keeps dependency, brittle text parsing.

### D4: Inline MFA logic (replaces aws-mfa)
**Decision**: Implement STS GetSessionToken + `~/.aws/credentials` writing in Go.
**Rationale**: aws-mfa (broamski) writes: `aws_access_key_id`, `aws_secret_access_key`, `aws_session_token`, `expiration` under `[sourceProfile]`. This is straightforward to replicate with AWS SDK Go v2 + `gopkg.in/ini.v1` for ini file read/write. Format stays identical — other tools (Terraform, boto3) are unaffected.
**Key detail**: `expiration` lives in `~/.aws/credentials`, not in bmc config. `[sourceProfile-long-term]` holds the permanent credentials and `aws_mfa_device`.

### D5: Shell wrapper for profsel
**Decision**: `bmc profsel` outputs `export AWS_PROFILE=<name>` to stdout. A shell function wrapper installed via `bmc install-shell-integration` wraps the binary call with `eval`.
**Rationale**: Go subprocesses cannot modify parent shell environment. This is the standard pattern used by direnv, rbenv, nvm.
**Shell wrapper**:
```bash
bmc() {
  if [[ "$1" == "profsel" ]]; then
    eval "$(command bmc profsel "$@")"
  else
    command bmc "$@"
  fi
}
```
**Alternatives**: Fish shell uses `bmc profsel | source` — supported as a variant.

### D6: TOML config (replaces config.env)
**Decision**: `~/.config/bmc/config.toml` replaces `~/.config/bmc/config.env`.
**Rationale**: Sourcing arbitrary shell files is a security risk (shell injection). TOML is typed, structured, and idiomatic for Go tools. `github.com/BurntSushi/toml` is the standard library.
**Migration**: `bmc doctor` warns if `config.env` exists but `config.toml` does not. No automatic migration.

### D7: Lazy prerequisite checks
**Decision**: Prerequisites checked only when the command that needs them is invoked, immediately before the action that requires them.
**Rationale**: Checking all deps at startup penalises users who never use SSM. `ec2connect` checks for `ssh` only after the user selects the SSH method in the TUI; checks for `aws`+`session-manager-plugin` only after SSM is selected.
**Error format**: Binary name, version (if detectable), required-by note, install instructions for apt/brew/nix-env/nix-profile/nixos-config.

### D8: syscall.Exec for interactive sessions
**Decision**: Use `syscall.Exec` (Unix exec) to hand off to `ssh`, `aws ssm start-session`, and `aws ecs execute-command`.
**Rationale**: These commands take over the terminal (PTY). `exec.Command` keeps bmc as parent process which interferes with signal handling and terminal control. `syscall.Exec` replaces the bmc process entirely — cleanest handoff.
**Windows**: Not a target platform; `syscall.Exec` is Unix-only which is acceptable.

### D9: Console command (replaces assumego)
**Decision**: Implement AWS federation URL flow natively.
**Rationale**: assumego (Granted) is a heavyweight dependency. The federation flow is well-documented:
1. STS AssumeRole → temporary credentials
2. POST to `https://signin.aws.amazon.com/federation?Action=getSigninToken`
3. Build sign-in URL with destination service
4. Open with `xdg-open` (Linux) / `open` (macOS)
No external binary needed.

### D11: Custom bubbletea list delegate
**Decision**: Use a custom `itemDelegate` instead of `list.NewDefaultDelegate()` for all `bubbles/list` instances.
**Rationale**: `list.NewDefaultDelegate()` requires items to implement the `list.DefaultItem` interface, which demands a `Title() string` method. Our `Item` struct has a exported field named `Title` — Go does not allow a method and a field with the same name on the same type. Without `Title()`, the default delegate renders every item as blank, making the list appear empty and returning an empty string on selection.
**Solution**: `itemDelegate` implements `list.ItemDelegate` directly and accesses `item.(Item).Title` (the field) for rendering.
**Constraint**: All `Item` values must keep the `Title` field name; adding a `Title()` method is not possible without renaming the field.

### D12: Cobra reserves `-h` for `--help`
**Decision**: Never use `-h` as a shorthand for any command flag; use `-i` for instance ID in `ec2connect`.
**Rationale**: Cobra automatically registers `-h`/`--help` on every command. Attempting to define another flag with shorthand `-h` causes a panic at startup (`unable to redefine 'h' shorthand`). The bash predecessor used `-h` for instance ID; this is a breaking change in flag naming.
**Impact**: `ec2connect -h <id>` (bash) → `ec2connect -i <id>` (Go).

### D13: Filter-aware Enter in bubbles/list
**Decision**: When `bubbles/list` is in `Filtering` state (user is typing a filter), the Enter key commits the filter and returns to browsing — it does NOT immediately select and quit.
**Rationale**: Intercepting Enter unconditionally caused the list to quit with an empty selection when the user typed to narrow results and pressed Enter to confirm the filter. The fix: check `m.list.FilterState() == list.Filtering` before treating Enter as a selection confirmation.

### D10: Distribution pipeline
**Decision**: GoReleaser for all distribution targets.
- GitHub Releases: binaries for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
- Homebrew tap: `wearetechnative/homebrew-tap` auto-updated by GoReleaser post-release
- Nix flake: `buildGoModule` with `vendorHash`; updated on release via `go mod vendor`

## Risks / Trade-offs

- **config.env → config.toml migration**: Users must manually migrate. `bmc doctor` will warn, but there is no auto-migration. → Provide migration docs in README and a `bmc migrate-config` helper command (stretch goal).
- **profsel shell wrapper**: Existing users sourcing `bmc profsel` must run `bmc install-shell-integration` once. → Document prominently in CHANGELOG and README.
- **ini file write safety**: Concurrent writes to `~/.aws/credentials` could corrupt the file. → Use file locking (`flock`) during credential writes.
- **session-manager-plugin still required**: SSM and ECS sessions cannot be zero-dep. → Document clearly; `bmc doctor` checks and provides install instructions.
- **Federation URL format changes**: AWS could change the federation sign-in API. → Isolate in `internal/awsops/console.go`; easy to update.
- **bubbletea in non-interactive terminals**: Piped or non-TTY contexts break TUI. → Detect `!term.IsTerminal(os.Stdin)` and fall back to plain text prompts.

## Migration Plan

1. New binary named `bmc` — replaces the bash script of the same name
2. Nix flake updated: `buildGoModule` replaces script installation
3. Users run `bmc install-shell-integration` once to install profsel wrapper
4. Users create `~/.config/bmc/config.toml` (guided by `bmc doctor` output)
5. `aws_mfa_device` and credentials remain in `~/.aws/credentials` — no change
6. Old bash scripts (`_bmclib.sh`, `ec2connect.sh`, etc.) removed from repo

**Rollback**: Keep previous bash `bmc` script tagged in git. Nix flake can pin to previous version.

## Open Questions

- Should `bmc migrate-config` be in scope for v1 to ease the config.env → config.toml transition?
- `ec2find` currently uses `selectProfileGroup` (gum-based multi-profile search). In Go this becomes concurrent goroutines per profile — confirm this is the desired behaviour.
- Should `tgselect.sh` be ported to `bmc togglsel` in a future change, or removed?
