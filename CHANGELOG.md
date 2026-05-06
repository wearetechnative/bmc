# BMC Changelog

## [0.2.13] - 06 May 2026

### Added
- **Release automation**: `release.sh` handles version bump (semver, dropping legacy 4th component), CHANGELOG finalization, Nix vendorHash auto-update, local build verification, git tag creation, and optional push
- **GitHub Actions release workflow**: Triggers GoReleaser on `v*` tag push — builds multi-platform binaries and updates the Homebrew tap
- **Homebrew tap**: Setup documented in `docs/homebrew-tap-setup.md`; install via `brew tap wearetechnative/tap && brew install bmc`
- **Version ldflags injection**: GoReleaser and Nix builds now stamp the binary with the correct version via `-X cmd.Version`; local `go build` falls back to the embedded `VERSION-bmc`
- **Back navigation in TUI menus**: Press ESC in any list to go back to the previous menu level; press Ctrl+C to cancel entirely. Works in `bmc profsel` (group → profile) and `bmc ecsconnect` (cluster → service → task → container)
- **Configurable EC2 columns**: Set `columns` in `[ec2]` section of `~/.config/bmc/config.toml` to control which columns appear in EC2 tables and in what order. Default order now puts `Name` second. All EC2 commands (`ec2ls`, `ec2connect`, `ec2stopstart`, `ec2scheduler`, `ec2find`) use the same column list
- **ec2ls formatted output**: `bmc ec2ls` now renders a bordered table with bold headers using lipgloss — pipeable and scrollable via terminal scrollback
- **Selection table improvements**: Interactive EC2 selection tables (ec2connect, ec2stopstart, ec2scheduler) now show a footer with row count and key hints, and adapt height to the terminal

### Fixed
- **TUI colors in eval context**: Running `eval "$(bmc profsel)"` now correctly shows colors and styling. Previously lipgloss detected piped stdout as a non-TTY and fell back to plain output despite the TUI rendering to `/dev/tty`
- **TUI list sizing**: List height is now clamped to the actual terminal height and updates on resize; spurious pagination dots no longer appear on short lists; eval-context and sub-group lists render correctly via `/dev/tty`
- **TUI blank rows**: Short lists (≤ 4 items, e.g. connection method and SSH user pickers) no longer show blank rows between the last item and the bottom of the list
- **profsel wrapper hint**: Running `bmc profsel` directly in a terminal now prints a tip on stderr pointing to `bmc install-shell-integration` when the shell wrapper is not in use
- **install-shell-integration on managed dotfiles**: When `~/.zshrc` or `~/.bashrc` is not writable (e.g. home-manager on NixOS), the command now prints actionable manual snippets for home-manager zsh/bash and Fish shell instead of a generic permission error

### Added
- **Go rewrite**: bmc is now a single self-contained Go binary — no runtime dependencies for core operations (gum, jq, awk, jsonify-aws-dotfiles, assumego, aws-mfa all eliminated)
- **`bmc doctor`**: New system health check command — verifies AWS config, credentials, optional tools (ssh, aws CLI v2, session-manager-plugin), MFA setup, and shell integration
- **`bmc install-shell-integration`**: Installs the profsel shell wrapper into `~/.zshrc` or `~/.bashrc`
- **Bubbletea TUI**: All interactive prompts now use bubbletea + bubbles (filterable lists, tables, spinner, text input)
- **MFA inlined**: STS GetSessionToken and `~/.aws/credentials` write logic built in — no external aws-mfa tool needed
- **AWS console native**: Federation URL flow built in — no external assumego/Granted tool needed
- **GoReleaser distribution**: Binary releases for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
- **Homebrew tap**: `brew install wearetechnative/tap/bmc`
- **Nix flake**: `buildGoModule`-based flake with nix-env, nix profile, and NixOS configuration support

### Changed
- **BREAKING**: Config file moved from `~/.config/bmc/config.env` (bash) to `~/.config/bmc/config.toml` (TOML)
- **BREAKING**: `source bmc profsel` replaced by shell wrapper + `eval "$(bmc profsel)"` — run `bmc install-shell-integration` once to set up
- Shell completions now via cobra built-in: `bmc completion bash/zsh` (replaces `bmc gencompletions`)
- `ec2find` now queries profiles concurrently using goroutines

### Removed
- Bash scripts removed: `_bmclib.sh`, `ec2connect.sh`, `ec2scheduler.sh`, `ecsconnect.sh`, `profsel.sh`, `aws-profile-select.sh`, `find-profile.awk`, `wouter.sh`

## 0.2.12.0 - 30 Mar 2026
- Fix: `bmc ec2ls` Name column now displays complete names when Name tag contains spaces
- Fix: `bmc profsel -p <profile>` now correctly recognizes source profiles (credentials-based profiles without role_arn)
- Fix: Removed confusing usage menu display when profile is not found - error message is now shown alone

## 0.2.11.0 - 25 mar 2026
- Feature: `bmc profsel --json` flag for machine-readable JSON output of profile selection
- Feature: File descriptor 3 support for JSON output, allowing progress messages to remain visible during interactive selection
- Feature: JSON output format: `{"source_profile": "...", "profile_name": "...", "profile_arn": "..."}`
- Feature: Error handling with JSON error messages: `{"error": "..."}`
- Enhancement: Intelligent output routing - JSON to fd 3 (when available) or stdout (fallback), progress to stdout or stderr
- Enhancement: Support for both interactive (`bmc profsel --json`) and non-interactive (`bmc profsel -p profile --json`) modes
- Enhancement: Backward compatible - works with and without fd 3 redirection
- Enhancement: Usage examples for scripting integration: `PROFILE=$(bmc profsel --json 3>&1 >/dev/null)`

## 0.2.10.0 - 23 mar 2026
- Enhancement: `bmc ec2scheduler` now displays Scheduler column showing if InstanceScheduler tag is configured (yes/no)
- Enhancement: `bmc ec2ls` now displays Hibernate values as "yes/no" instead of "true/false/None" for better readability
- Enhancement: `bmc ec2ls` now displays Scheduler column showing if InstanceScheduler tag is configured (yes/no)
- Enhancement: `bmc ec2find` now includes Scheduler column in search results
- Note: Automated tools parsing CSV output from `bmc ec2ls` or `bmc ec2find` may need updates due to new column order

## 0.2.9.0 - 21 mar 2026
- Feature: New `bmc gencompletions` command to generate shell completion scripts for bash and zsh
- Feature: Tab-completion support for all bmc commands in bash
- Feature: Tab-completion with command descriptions for all bmc commands in zsh
- Enhancement: Dynamic command discovery in completion scripts - automatically includes new commands
- Enhancement: Multiple installation options for shell completion (direct sourcing, file-based, system-wide)
- **Breaking**: `bmc ec2scheduler` now manages Ignore_scheduler tags instead of toggling InstanceScheduler/InstanceScheduler_DISABLED tags
- Feature: Set time-based scheduler overrides with Ignore_scheduler tag (e.g., "22:00 Europe/Amsterdam")
- Feature: Interactive menu to set or remove scheduler overrides
- Feature: Time format validation (HH:MM 24-hour format)
- Feature: Free-form timezone input with helpful examples
- Enhancement: Temporary overrides that automatically expire and return instance to normal schedule
- Enhancement: Table display shows Ignore_scheduler status and ignore-until time
- Enhancement: Guide users to add InstanceScheduler tags via AWS Console for untagged instances (unchanged)
- Enhancement: add back navigation in profile selection - users can now return to group menu by canceling profile selection instead of restarting command

## 0.2.8.0 - 21 jan 2026
- Feature: `bmc console` respects AWS_PROFILE environment variable when set
- Feature: `bmc console -p` (without value) forces profile selection even when AWS_PROFILE is set
- Enhancement: Reduced friction when AWS_PROFILE is already configured
- Fix: `bmc profsel` no longer exits the shell when profile selection is cancelled (Ctrl-C) or no profile is chosen
- Feature: `bmc ec2connect` automatically selects SSH connection when `-u` (username) or `-i` (identity file) flags are provided, eliminating unnecessary connection type prompt
- Feature: `bmc ec2connect` now prompts to start stopped EC2 instances before connecting, streamlining the workflow
- Feature: New config option `BMC_AUTO_START_STOPPED_INSTANCES` to control stopped instance behavior (values: "always", "never", "prompt")
- Enhancement: Improved error messages in `bmc ec2connect` - removed redundant "Not executing the SSH-command" text
- Fix: TOTP script now properly executes with command-line arguments using correct array expansion
- Fix: Clipboard copy now uses correct variable name `clipboardCopyCommand` instead of undefined `clipboardCommand`
- Enhancement: Clear feedback message when TOTP script is not configured instead of displaying undefined variable
- Fix: Clipboard copy now properly validates command exists before showing success message
- Enhancement: Added informative message before executing TOTP script to improve user awareness
- Enhancement: Improved MFA session messages to be more user-friendly and less debug-like

## 2.7.0 - 18 sept 2025
- open profile selection when AWS_PROFILE is not set
- use filter in stead of table/choose
- cleanups

## 0.2.6.7
- cleanups

## 0.2.6.6
- Add -s flag to console option. user bmc console -s <service-name> to directly open the console with the prefered service.\

## 0.2.6.5
- Add -p flag to console option. user bmc console -p <profile-name> to directly open the console with the profile.\

## 0.2.6.4
- Set session duration to 3600s

## 0.2.6.3
- Fix e2stopstart function, rewriting function call

## 0.2.6.2
- Fix spinner ec2stopstart function

## 0.2.6.0
- Refactor ec2ls function, integrated in main library
- ec2find option. Search for keyword in selected aws-profile

## 0.2.5.3
- Fix: table height to fix items not being visible

## 0.2.5.2
- Fix: table height to fix items not being visible

## 0.2.5.1
- Feature: Search option for ec2ls. Now you can search through the output list for a string
- Fix bugs

## 0.2.4
- fix renaming function error ec2connect

## 0.2.3
- rename ec2ssh option to ec2connect
- add options ssm and ssh connection method connecting to ec2

## 0.2.2
- fix usage
- Feature add usage in help
- make VERSION unique for flake distribution

## 0.2.1
- another fix usage

## 0.2.0
- Breaking: renamed cli command
- Feature: add more commands to bmc
- Feature: more refactoring
- Fix: new sourcing fix

## 0.1.1
- Fix: sourcing, detect being sourced 

## 0.1.0

- new official project
