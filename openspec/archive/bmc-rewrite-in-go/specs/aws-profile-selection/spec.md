## ADDED Requirements

### Requirement: Profile selection with group filtering
The system SHALL present AWS profiles grouped by their `group` field from `~/.aws/config`, allowing the user to first filter by group then select a profile using a fuzzy-filterable list (bubbles/list).

#### Scenario: Interactive group and profile selection
- **WHEN** user runs `bmc profsel` with no flags and no `AWS_PROFILE` set
- **THEN** system shows a filterable list of profile groups, then a filterable list of profiles in the selected group, then outputs `export AWS_PROFILE=<selected-profile>`

#### Scenario: Pre-select profile with -p flag
- **WHEN** user runs `bmc profsel -p <profile-name>`
- **THEN** system skips the interactive selection and outputs `export AWS_PROFILE=<profile-name>`

#### Scenario: List profiles with -l flag
- **WHEN** user runs `bmc profsel -l`
- **THEN** system prints all profiles in tabular format (Group, Name, ARN number) to stdout

#### Scenario: JSON output with --json flag
- **WHEN** user runs `bmc profsel --json`
- **THEN** system outputs `{"source_profile": "...", "profile_name": "...", "profile_arn": "..."}` as JSON to stdout

#### Scenario: Cancelled selection
- **WHEN** user presses Escape or Ctrl+C during profile selection
- **THEN** system exits cleanly with no output and exit code 0

#### Scenario: Filter committed before selection
- **WHEN** user types characters to filter the list and presses Enter while the list is in filter mode
- **THEN** system commits the filter (narrows the list) without immediately selecting an item; a second Enter then selects the highlighted item

#### Scenario: Items rendered correctly
- **WHEN** the filterable list is displayed
- **THEN** all items are visible with their title text; none appear as blank rows

### Requirement: AWS config parsed natively
The system SHALL parse `~/.aws/config` and `~/.aws/credentials` directly in Go without any external tools (no jq, awk, jsonify-aws-dotfiles).

#### Scenario: Config file present
- **WHEN** `~/.aws/config` exists with profile entries
- **THEN** system reads all profiles including `role_arn`, `source_profile`, and `group` fields

#### Scenario: Config file absent
- **WHEN** `~/.aws/config` does not exist
- **THEN** system exits with an actionable error message

### Requirement: Source profile resolution
The system SHALL resolve the `source_profile` for a given profile, handling both config-based role profiles and credentials-based direct profiles.

#### Scenario: Role profile with source_profile
- **WHEN** selected profile has a `role_arn` and `source_profile` in config
- **THEN** `sourceProfile` is set to the value of `source_profile`

#### Scenario: Credentials-only profile
- **WHEN** selected profile exists in `~/.aws/credentials` but not as a role in config
- **THEN** `sourceProfile` is set to the profile name itself
