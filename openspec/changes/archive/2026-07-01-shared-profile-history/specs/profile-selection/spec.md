## ADDED Requirements

### Requirement: Recently used groups appear at the top of group selection
The interactive profile selector SHALL surface recently used account groups at the top of the group list, each marked with description "recent", when the user has previous profile selection history.

#### Scenario: Group list with history
- **WHEN** the user enters interactive profile selection and has previously selected profiles
- **THEN** groups containing recently used profiles SHALL appear at the top of the group list
- **AND** each recent group SHALL be shown with description "recent"
- **AND** remaining groups SHALL follow in their normal order without the "recent" label

#### Scenario: Group list without history
- **WHEN** the user enters interactive profile selection and has no previous history
- **THEN** all groups SHALL appear in their normal order without any "recent" labels

### Requirement: Recently used profiles appear at the top of profile selection
Within a selected group, the interactive profile selector SHALL surface recently used profiles at the top of the profile list, each marked with description "recent".

#### Scenario: Profile list with recent profiles in group
- **WHEN** the user selects an account group and has previously selected profiles within that group
- **THEN** recently used profiles within that group SHALL appear at the top of the profile list
- **AND** each recent profile SHALL be shown with description "recent"
- **AND** remaining profiles SHALL follow with their normal account ID / role name description

#### Scenario: Profile list with no recent profiles in group
- **WHEN** the user selects an account group that contains no recently used profiles
- **THEN** all profiles SHALL appear in their normal order with account ID / role name descriptions

### Requirement: Shared history across all interactive profile selectors
All commands that trigger interactive profile selection (`bmc console`, `bmc profsel`, and commands using `ensureAWSProfile()`) SHALL read from and write to the same history store (`~/.local/share/bmc/profile-history.json`), so that a profile selected in any command surfaces as recent in all others.

#### Scenario: Profile selected in console appears recent in profsel
- **WHEN** the user interactively selects a profile via `bmc console`
- **THEN** that profile SHALL appear as recent the next time the user runs `bmc profsel` interactively

#### Scenario: Profile selected in profsel appears recent in ec2connect
- **WHEN** the user interactively selects a profile via `bmc profsel`
- **THEN** that profile SHALL appear as recent the next time `bmc ec2connect` triggers interactive selection

#### Scenario: Non-interactive selection does not write history
- **WHEN** a profile is resolved via the `-p` flag or `AWS_PROFILE` environment variable (no interactive picker shown)
- **THEN** that selection SHALL NOT be written to history
