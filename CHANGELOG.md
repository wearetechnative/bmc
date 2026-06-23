# BMC Changelog

## NEXT VERSION

## [0.6.0] - 23 Jun 2026

### Added
- **`bmc watcher`**: Background daemon that automatically refreshes AWS console federation sessions before they expire, keeping your browser logged in without manual re-authentication
  - `bmc console --watch` / `-w`: registers the session with the watcher after opening the console
  - `bmc watcher start`: starts the keep-alive daemon (auto-started by `--watch`)
  - `bmc watcher stop`: stops the daemon
  - `bmc watcher status`: shows active sessions with next-refresh countdown, indicates whether CDP or tab-fallback is used
  - `bmc watcher setup`: configures Firefox for CDP-based invisible refresh (writes `user.js` in the default Firefox profile)
- **CDP-based invisible session refresh**: when Firefox Remote Debugging Protocol is available (port 9222 by default), session refresh executes a `fetch()` call inside an existing AWS console tab — no new tab opened, no focus change
  - Falls back to local refresh page (fetch + auto-close) if CDP is unavailable or no console tab is found
  - Falls back to opening the federation URL directly if the local page fails
- **`watcher.firefox_debug_port` config option**: controls the CDP port (default `9222`); set to `0` to disable CDP entirely
- **`mfa.profile_scripts` config option**: per-source-profile TOTP script overrides — map each AWS source profile name to a different TOTP script command; `totp_script` remains the global fallback for profiles not listed

## [0.5.5] - 28 May 2026

### Fixed
- **`console -s <service>` 400 Bad Request**: URL generation now uses region-aware path-based console URLs (`https://<region>.console.aws.amazon.com/<service>/home`) instead of the broken subdomain format (`https://<service>.console.aws.amazon.com/`) that caused 400 errors for most services
  - Region is resolved from the selected AWS profile
  - The `-s` flag now accepts console sub-paths (e.g., `systems-manager/parameters`) for deep linking directly to a specific page

## [0.5.4] - 21 May 2026

## [0.5.3] - 20 May 2026

### Fixed
- **`profsel -l` inside eval wrapper**: Profile list output now goes to stderr, preventing shell errors when `bmc profsel` is wrapped in `eval`
- **Shell integration wrapper**: Fixed two bugs — `profsel` arg was doubled in the underlying command, and `-l`/`--list` output was incorrectly eval'd. Wrapper now skips eval for list mode and uses `${@:2}` (bash/zsh) / `$argv[2..]` (fish) to avoid arg duplication

### Added
- **`-k`/`--key` flag for `ec2connect`**: Specify an SSH identity file (`.pem` key pair) when connecting via SSH
  - Path is passed directly to `ssh -i` — no validation performed by bmc
  - Providing `-k` automatically selects SSH as the connection method (no method picker shown)
  - Restores functionality from the original bash-era `ec2connect.sh`

## [0.5.2] - 18 May 2026

### Added
- **`-p`/`--profile` flag on AWS commands**: Pass `-p <profile>` directly to `ec2connect`, `ec2ls`, `ec2`, `ec2scheduler`, `ec2stopstart`, and `ecsconnect` to specify an AWS profile without setting `AWS_PROFILE` in the environment
  - Flag takes priority over `AWS_PROFILE` env var; MFA check still runs
  - `console` and `profsel` retain their existing `-p` behaviour unchanged

## [0.5.1] - 13 May 2026

### Added
- **Documentation site**: Full documentation available at [bmc.technative.cloud](https://bmc.technative.cloud) — Hugo site with TechNative branding, multi-page navigation covering installation, setup, commands, and advanced topics
  - Automatically deployed via GitHub Actions on every push to `main`
  - Hosted on GitHub Pages at `bmc.technative.cloud`
- **Shortened README**: README now links to the docs site as the primary reference

## [0.5.0] - 13 May 2026

### Added
- **JSON output for `ec2ls` and `ec2find`**: Pass `--json` to get a machine-readable JSON array instead of a table — useful for scripting and piping into `jq`
  - `bmc ec2ls --json` outputs all instances; all fields always included regardless of `ec2.columns` config
  - `bmc ec2find <search> --json` outputs matching instances including the `Profile` field; interactive group selection still works via the terminal
  - JSON keys follow AWS CLI PascalCase convention (`InstanceId`, `PrivateIpAddress`, etc.)
- **`bmc ec2` unified command**: Select an EC2 instance and immediately act on it — no need to repeat instance selection across separate commands
  - Optional positional search argument: `bmc ec2 nginx` filters by name, ID, private IP, or public IP (case-insensitive); single match skips the picker entirely
  - Action menu after selection: **Connect SSH**, **Connect SSM**, **Start instance** / **Stop instance** (label adapts to current state), **Toggle scheduler**
  - All actions reuse the same logic as `ec2connect`, `ec2stopstart`, and `ec2scheduler` — existing commands are unchanged

## [0.4.0] - 13 May 2026

### Breaking Changes
- **`mfa.clipboard_command` renamed to `mfa.copy_command`**: Update your `~/.config/bmc/config.json` — the old field name is silently ignored

### Added
- **MFA paste command**: New `mfa.paste_command` config field for keystroke simulation after clipboard copy (e.g. `xdotool key ctrl+v`, `wtype -k ctrl+v`). Runs 300ms after a successful copy, no action if copy fails.
- **totp_script via sh -c**: Quoted arguments with spaces now work correctly (e.g. `rbw code "My Entry (new)"`). Interactive TUI selection tools (e.g. `gum filter`) render on the terminal and accept keyboard input.
- **Fish shell support**: `bmc install-shell-integration` now supports Fish shell
  - Writes a wrapper function to `~/.config/fish/functions/bmc.fish` (auto-loaded, no restart needed)
  - On NixOS: skips auto-install and prints a `programs.fish.functions` home-manager snippet instead
  - Already-installed check via file existence
- **ec2connect filter by name**: `bmc ec2connect <search>` filters instances before selection using a case-insensitive substring match on instance name, ID, private IP, and public IP
  - Single match: connects immediately without showing a picker
  - Multiple matches: shows the interactive picker with only the matching instances
  - No matches: returns a clear error
  - `-i <id>` flag still takes precedence (with a warning if a search arg is also given)

## [0.3.0] - 08 May 2026

### Breaking Changes
- **Config file format changed**: `~/.config/bmc/config.toml` → `~/.config/bmc/config.json`. Field names are unchanged — only the file format and filename differ. Run `bmc` once with the old `config.toml` still in place to get a migration hint with the exact JSON equivalent.

  Before (`config.toml`):
  ```toml
  [mfa]
  enabled = true
  totp_script = "/usr/bin/rbw get my-aws-mfa-entry --field totp"

  [console]
  firefox_containers = true
  ```
  After (`config.json`):
  ```json
  {
    "mfa": {
      "enabled": true,
      "totp_script": "/usr/bin/rbw get my-aws-mfa-entry --field totp"
    },
    "console": {
      "firefox_containers": true
    }
  }
  ```

### Added
- **Firefox container support**: Set `firefox_containers: true` in the `console` section of `~/.config/bmc/config.json` to open the AWS console in a dedicated Firefox container tab via the [Granted](https://addons.mozilla.org/en-US/firefox/addon/granted/) extension. The container is named after the AWS profile.
- **Chrome profile isolation (experimental)**: Set `chrome_profiles: true` in the `console` section of `~/.config/bmc/config.json` to open the AWS console in a bmc-managed isolated Chrome profile per AWS profile. On first use, extensions and preferences are seeded from your default Chrome profile. Works with Chrome, Brave, Chromium, and Edge via the optional `chrome_binary` setting.
- **Console profile history**: `bmc console` now remembers recently used profiles. The last 10 profiles are shown at the top of the interactive selector with a "recent" label. History is stored in `~/.local/share/bmc/console-history.json` and is only updated on successful interactive opens.
- **`bmc console -p` forces profile selection**: Running `bmc console -p` without a profile name opens the interactive profile selector, even when `AWS_PROFILE` is already set in the environment.

### Fixed
- **MFA credentials write**: `aws_security_token` (legacy alias used by aws-mfa and older AWS SDKs) is now kept in sync with `aws_session_token` when refreshing a session. Previously the stale alias caused authentication failures for tools reading the old key name.
- **MFA credentials file safety**: Credentials file parser now uses consistent options (`IgnoreInlineComment`) in both read and write paths, and returns an error instead of silently discarding all credentials when the file cannot be parsed.
- **MFA check on pre-set profile**: When `AWS_PROFILE` is already set in the environment, `bmc` now validates and refreshes the MFA session before executing any AWS operation. Previously this was skipped, causing `InvalidClientTokenId` errors when the session expired.

### Changed
- **Documentation**: Rewrote `docs/aws-profile-select.md` to reflect the current Go CLI, including shell integration via `eval "$(bmc profsel)"` and `/dev/tty` color rendering behavior
- **OpenSpec specs**: Translated `tui-color-rendering` spec from Dutch to English
- **Repository cleanup**: Removed obsolete bash-era files (`_bmclib.sh`, `_get_var_file.sh`, `tgselect.sh`) and untracked the accidentally committed `bmc-go` build artifact

## [0.2.15] - 06 May 2026

### Fixed
- **Homebrew Formula**: switch back to `brews` for proper Formula syntax (`brew install bmc` without `--cask`)

## [0.2.14] - 06 May 2026

### Fixed
- **Homebrew formula directory**: Formula now published to `Formula/bmc.rb` (was root `bmc.rb`)

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
