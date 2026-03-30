# BMC Changelog

## NEXT VERSION
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
