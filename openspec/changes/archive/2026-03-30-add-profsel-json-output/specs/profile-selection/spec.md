# profile-selection Spec Delta

## ADDED Requirements

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

## MODIFIED Requirements

None. This change is purely additive and does not modify existing requirements.

## REMOVED Requirements

None. All existing functionality remains intact.
