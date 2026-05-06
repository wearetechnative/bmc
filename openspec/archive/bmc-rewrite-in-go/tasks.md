## 1. Project Scaffolding

- [x] 1.1 Initialise Go module (`go mod init github.com/wearetechnative/bmc`)
- [x] 1.2 Add core dependencies: cobra, bubbletea, bubbles, lipgloss, aws-sdk-go-v2, BurntSushi/toml, gopkg.in/ini.v1
- [x] 1.3 Create directory structure: `cmd/`, `internal/awsconfig/`, `internal/awsops/`, `internal/mfa/`, `internal/ui/`
- [x] 1.4 Create `main.go` wiring cobra root command
- [x] 1.5 Embed `VERSION-bmc` via `go:embed` and wire into `cmd/version.go`

## 2. AWS Config Parsing (internal/awsconfig)

- [x] 2.1 Implement `~/.aws/config` ini parser: read profiles with `role_arn`, `source_profile`, `group` fields
- [x] 2.2 Implement `~/.aws/credentials` ini parser: read access keys, session tokens, expiration
- [x] 2.3 Implement source profile resolution logic (role profiles vs credentials-only profiles)
- [x] 2.4 Implement profile group extraction and deduplication
- [x] 2.5 Write unit tests for config parsing with fixture files

## 3. Config Management (internal/config)

- [x] 3.1 Define TOML config struct with `mfa` and `ec2` sections and defaults
- [x] 3.2 Implement `~/.config/bmc/config.toml` reader with graceful absent-file handling
- [x] 3.3 Detect legacy `~/.config/bmc/config.env` and surface warning via doctor

## 4. TUI Components (internal/ui)

- [x] 4.1 Implement filterable list component (bubbles/list) for profile group and profile selection
- [x] 4.2 Implement table component (bubbles/table) for EC2 instance display
- [x] 4.3 Implement spinner component (bubbles/spinner) for async wait states
- [x] 4.4 Implement text input component (bubbles/textinput) for MFA code prompt and custom username
- [x] 4.5 Implement confirm prompt (huh or custom) for stop/start and instance start confirmations
- [x] 4.6 Add non-TTY detection: fall back to plain text prompts when stdin is not a terminal

## 5. MFA Session Management (internal/mfa)

- [x] 5.1 Implement expiration check: read `expiration` from `~/.aws/credentials` under `[sourceProfile]`
- [x] 5.2 Implement MFA device lookup from `[sourceProfile-long-term]` section
- [x] 5.3 Implement TOTP acquisition: run `totp_script` from config or prompt user
- [x] 5.4 Implement clipboard copy via `clipboard_command` (non-fatal if command not found)
- [x] 5.5 Implement `sts.GetSessionToken` call via AWS SDK Go v2
- [x] 5.6 Implement credential write: update `[sourceProfile]` section in `~/.aws/credentials` with file lock
- [x] 5.7 Write unit tests for expiration check and credential write format

## 6. AWS Profile Selection (cmd/profsel.go)

- [x] 6.1 Implement `bmc profsel` command using filterable group → profile TUI flow
- [x] 6.2 Implement `-p` flag for pre-selected profile
- [x] 6.3 Implement `-l` flag for tabular profile listing
- [x] 6.4 Implement `--json` flag outputting `{source_profile, profile_name, profile_arn}`
- [x] 6.5 Ensure stdout output is `export AWS_PROFILE=<name>` for eval consumption
- [x] 6.6 Wire MFA session check into profsel flow after profile selection

## 7. Shell Integration (cmd/install_shell.go)

- [x] 7.1 Implement `bmc install-shell-integration` detecting bash vs zsh from `$SHELL`
- [x] 7.2 Write profsel wrapper function to `~/.zshrc` or `~/.bashrc` if not already present
- [x] 7.3 Handle unsupported shells: print wrapper for manual installation

## 8. AWS Console (cmd/console.go + internal/awsops/console.go)

- [x] 8.1 Implement STS AssumeRole using AWS SDK Go v2 with 3600s duration
- [x] 8.2 Implement federation sign-in token request to `https://signin.aws.amazon.com/federation`
- [x] 8.3 Build final console sign-in URL with optional service destination
- [x] 8.4 Implement browser opener: `xdg-open` on Linux, `open` on macOS
- [x] 8.5 Wire `bmc console` with profile selection, `-p` flag, `-s` flag, and existing `AWS_PROFILE` support

## 9. EC2 Operations (cmd/ec2*.go + internal/awsops/ec2.go)

- [x] 9.1 Implement `ec2ls`: call `DescribeInstances`, normalise Hibernate/Scheduler columns, render bubbles/table
- [x] 9.2 Implement `ec2stopstart`: instance TUI selection, state check, start/stop/hibernate with spinner wait
- [x] 9.3 Implement `ec2scheduler`: instance TUI selection, read current tag, toggle InstanceScheduler tag
- [x] 9.4 Implement `ec2find`: profile group selection, concurrent `DescribeInstances` across profiles, merged result table with Profile column
- [x] 9.5 Implement `useOrSelectAWSProfile` equivalent: use `AWS_PROFILE` if set, else invoke profsel TUI

## 10. EC2 Connect (cmd/ec2connect.go)

- [x] 10.1 Implement instance table selection with `-h` flag override
- [x] 10.2 Implement stopped instance handling per `ec2.auto_start_stopped` config
- [x] 10.3 Implement connection method selection (SSH / SSM) via filterable list
- [x] 10.4 Implement SSH path: user selection TUI, prerequisite check for `ssh`, `syscall.Exec` to `ssh <user>@<id>`
- [x] 10.5 Implement SSM path: prerequisite check for `aws` v2 + `session-manager-plugin`, `syscall.Exec` to `aws ssm start-session --target <id>`

## 11. ECS Connect (cmd/ecsconnect.go + internal/awsops/ecs.go)

- [x] 11.1 Implement cluster listing and filterable selection with breadcrumb display
- [x] 11.2 Implement service listing and selection for chosen cluster
- [x] 11.3 Implement running task listing and selection for chosen service
- [x] 11.4 Implement container listing and selection for chosen task
- [x] 11.5 Implement prerequisite check for `aws` v2 + `session-manager-plugin`
- [x] 11.6 Implement `syscall.Exec` handoff to `aws ecs execute-command`

## 12. Prerequisites & Doctor (cmd/doctor.go + internal/prereqs)

- [x] 12.1 Implement prerequisite checker: binary lookup with version detection for `ssh`, `aws` (v1 vs v2), `session-manager-plugin`
- [x] 12.2 Implement lazy check helpers called from ec2connect and ecsconnect with formatted error output (all install methods including 3 Nix options)
- [x] 12.3 Implement `bmc doctor`: run all checks across Core, Optional, MFA, Shell integration, Legacy config categories
- [x] 12.4 Exit code: 0 if all checks pass, non-zero if any fail

## 13. Distribution

- [x] 13.1 Create `.goreleaser.yml` with build matrix (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64)
- [x] 13.2 Configure GoReleaser Homebrew tap section for `wearetechnative/homebrew-tap`
- [x] 13.3 Update `flake.nix` to use `buildGoModule` with correct `vendorHash`
- [x] 13.4 Update `package.nix` for Go build
- [x] 13.5 Add `go mod vendor` to release process (for Nix `vendorHash` reproducibility)
- [x] 13.6 Update README with new install methods (brew, nix-env, nix profile, NixOS) and migration guide

## 14. Cleanup & Migration

- [x] 14.1 Remove bash scripts from repo: `_bmclib.sh`, `bmc` (bash), `ec2connect.sh`, `ec2scheduler.sh`, `ecsconnect.sh`, `profsel.sh`, `aws-profile-select.sh`, `find-profile.awk`, `wouter.sh`
- [x] 14.2 Update CHANGELOG with breaking changes (config.env → config.toml, profsel shell wrapper)
- [x] 14.3 Verify `release.sh` version bump still works with new binary
