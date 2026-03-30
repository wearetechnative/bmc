# profile-selection Specification

## Purpose
TBD - created by archiving change fix-profsel-shell-exit. Update Purpose after archive.
## Requirements
### Requirement: Graceful Cancellation
The `bmc profsel` command SHALL handle user cancellation gracefully without exiting the current shell when sourced.

#### Scenario: User cancels profile selection with Ctrl-C
- **WHEN** user runs `. bmc profsel` (sourced) and presses Ctrl-C during profile selection
- **THEN** the command SHALL return to the shell prompt without exiting the shell
- **AND** no AWS_PROFILE environment variable SHALL be set

#### Scenario: User cancels profile selection by not choosing
- **WHEN** user runs `. bmc profsel` (sourced) and closes the selection menu without choosing a profile
- **THEN** the command SHALL return to the shell prompt without exiting the shell
- **AND** no AWS_PROFILE environment variable SHALL be set

#### Scenario: Selection cancelled in executed mode
- **WHEN** user runs `bmc profsel` (not sourced) and cancels the selection
- **THEN** the command SHALL exit normally with appropriate status code
- **AND** display a message indicating no profile was selected

### Requirement: Early Return on Empty Selection
The `bmc profsel` command SHALL check if a profile was selected before proceeding with MFA setup.

#### Scenario: No profile selected, MFA not called
- **WHEN** profile selection returns empty result (due to cancellation or no selection)
- **THEN** the command SHALL not call MFA setup functions
- **AND** SHALL return immediately with appropriate feedback

### Requirement: Safe Error Handling in Sourced Context
Functions called by sourced scripts SHALL use `return` instead of `exit` to prevent closing the user's shell.

#### Scenario: Error in sourced profsel execution
- **WHEN** an error occurs during sourced `bmc profsel` execution
- **THEN** the error handling SHALL use `return` statements
- **AND** SHALL not use `exit` statements that would close the parent shell

### Requirement: Back Navigation in Profile Selection
The `selectAWSProfile` function SHALL allow users to return to the group selection menu when they cancel profile selection, without exiting the entire command.

#### Scenario: User cancels profile selection and returns to group menu
- **GIVEN** user has selected a profile group
- **WHEN** user cancels at the profile selection stage (presses Ctrl-C or ESC)
- **THEN** the system SHALL return to the group selection menu
- **AND** allow the user to select a different group or cancel entirely

#### Scenario: User cancels at group selection to exit
- **GIVEN** user is at the group selection menu
- **WHEN** user cancels at the group selection stage (presses Ctrl-C or ESC)
- **THEN** the system SHALL exit the selection process gracefully
- **AND** unset selectedProfileName variable
- **AND** return control to the shell without closing it (when sourced)

#### Scenario: User successfully selects profile after returning to group menu
- **GIVEN** user has cancelled profile selection and returned to group menu
- **WHEN** user selects a different group
- **AND** selects a valid profile from that group
- **THEN** the system SHALL proceed with the selected profile normally
- **AND** set all profile variables as expected (selectedProfileName, selectedProfileARN, sourceProfile)

#### Scenario: Navigation loop with preferred profile flag
- **GIVEN** user runs profsel with `-p <profile-name>` flag
- **WHEN** the selectAWSProfile function is called
- **THEN** the system SHALL bypass both group and profile selection menus
- **AND** use the specified profile directly
- **AND** not enter the navigation loop

### Requirement: User Feedback During Navigation
The system SHALL provide clear feedback when navigating between menus to avoid user confusion.

#### Scenario: Feedback when returning to group selection
- **GIVEN** user cancels at profile selection stage
- **WHEN** the system returns to group selection
- **THEN** the system MAY display a brief message indicating return to group selection
- **OR** proceed directly to group selection without message (acceptable UX)

### Requirement: JSON Output Mode
The `bmc profsel` command SHALL support a `--json` flag that outputs profile selection results as valid JSON. The output destination SHALL be file descriptor 3 when available, falling back to stdout when fd 3 is not available.

#### Scenario: Successfully output JSON for selected profile (non-interactive)
- **GIVEN** user has a valid AWS profile named "my-dev-profile"
- **AND** the profile has source_profile "my-org" and role ARN "arn:aws:iam::123456789012:role/DevRole"
- **WHEN** user runs `bmc profsel -p my-dev-profile --json`
- **THEN** the command SHALL output valid JSON to stdout
- **AND** the JSON SHALL contain key "source_profile" with value "my-org"
- **AND** the JSON SHALL contain key "profile_name" with value "my-dev-profile"
- **AND** the JSON SHALL contain key "profile_arn" with value "arn:aws:iam::123456789012:role/DevRole"
- **AND** no other text SHALL be written to stdout (only JSON)
- **AND** the command SHALL exit with status code 0

#### Scenario: JSON output with interactive selection
- **WHEN** user runs `bmc profsel --json` without -p flag
- **THEN** the command SHALL display interactive gum menus for profile selection
- **AND** user can select profile group and profile using gum filters
- **WHEN** user completes the selection
- **THEN** the command SHALL output valid JSON to stdout with the selected profile information
- **AND** the command SHALL exit with status code 0

#### Scenario: JSON output for non-existent profile
- **WHEN** user runs `bmc profsel -p nonexistent-profile --json`
- **THEN** the command SHALL output valid JSON to stdout
- **AND** the JSON SHALL contain key "error" indicating profile not found
- **AND** the command SHALL exit with non-zero status code

#### Scenario: JSON output when selection is cancelled
- **WHEN** user runs `bmc profsel --json` and cancels the interactive selection
- **THEN** the command SHALL output valid JSON to stdout
- **AND** the JSON SHALL contain key "error" with value "no profile selected"
- **AND** the command SHALL exit with non-zero status code

#### Scenario: Progress messages with fd 3 available
- **WHEN** user runs `bmc profsel --json 3>&1 >/dev/null` (fd 3 redirected)
- **THEN** all MFA progress messages SHALL be written to stdout
- **AND** all "Using AWS source-profile" messages SHALL be written to stdout
- **AND** the final JSON result SHALL be written to fd 3
- **AND** progress messages remain visible during interactive selection

#### Scenario: Progress messages without fd 3 (backward compatible)
- **WHEN** user runs `bmc profsel --json 2>/dev/null` (fd 3 not used)
- **THEN** all MFA progress messages SHALL be written to stderr
- **AND** all "Using AWS source-profile" messages SHALL be written to stderr
- **AND** the final JSON result SHALL be written to stdout
- **AND** stderr can be suppressed with 2>/dev/null for clean JSON output

#### Scenario: JSON mode does not set environment variable
- **WHEN** user runs `. bmc profsel -p <profile> --json` (sourced)
- **THEN** the AWS_PROFILE environment variable SHALL NOT be set
- **AND** JSON output SHALL be written to stdout
- **AND** no "Source this script..." message SHALL be displayed
- **AND** the command SHALL return without modifying the environment

#### Scenario: JSON output is parseable by jq
- **WHEN** user runs `bmc profsel -p <profile> --json | jq .`
- **THEN** jq SHALL successfully parse the output
- **AND** SHALL not produce any parsing errors

#### Scenario: Clean JSON output to variable with visible progress
- **WHEN** user runs `PROFILE=$(bmc profsel --json 3>&1 >/dev/null)`
- **THEN** the variable SHALL contain only valid JSON
- **AND** progress messages SHALL be visible during execution
- **AND** the JSON SHALL be captured in the variable

#### Scenario: Clean JSON output without progress
- **WHEN** user runs `bmc profsel -p <profile> --json 2>/dev/null`
- **THEN** only valid JSON SHALL be output
- **AND** the output SHALL be a single JSON object with no additional text

### Requirement: JSON Flag Compatibility
The `--json` flag SHALL be compatible with the `-p` flag for non-interactive use and SHALL work without `-p` for interactive use.

#### Scenario: JSON combined with profile listing flag
- **WHEN** user runs `bmc profsel -l --json`
- **THEN** the command SHALL either ignore --json and show profile list
- **OR** show an error indicating incompatible flags
- **AND** SHALL NOT output JSON format

#### Scenario: JSON combined with preferred profile flag
- **WHEN** user runs `bmc profsel -p my-profile --json`
- **THEN** the command SHALL select the specified profile non-interactively
- **AND** output the result in JSON format
- **AND** skip interactive gum menus

#### Scenario: JSON without preferred profile flag
- **WHEN** user runs `bmc profsel --json`
- **THEN** the command SHALL display interactive profile selection
- **AND** output the result in JSON format after selection

