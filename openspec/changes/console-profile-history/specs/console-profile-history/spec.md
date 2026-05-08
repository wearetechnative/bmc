## ADDED Requirements

### Requirement: Console profile history is persisted after use
After a successful `bmc console` invocation with interactive profile selection, the selected profile name SHALL be saved to the history file at `$XDG_DATA_HOME/bmc/console-history.json` (defaulting to `~/.local/share/bmc/console-history.json`).

#### Scenario: Profile saved after successful open
- **WHEN** the user interactively selects a profile in `bmc console`
- **AND** the console opens successfully
- **THEN** the selected profile SHALL be prepended to the history file
- **AND** the history file SHALL contain at most 10 entries
- **AND** duplicate entries SHALL be removed (profile appears only once, at the top)

#### Scenario: Profile not saved on cancel
- **WHEN** the user opens `bmc console` but cancels without selecting a profile
- **THEN** the history file SHALL NOT be modified

#### Scenario: Profile not saved on non-interactive invocation
- **WHEN** `bmc console -p <profile>` is used
- **OR** `AWS_PROFILE` is set in the environment
- **THEN** the history file SHALL NOT be modified

#### Scenario: History file does not exist yet
- **WHEN** the history file does not exist
- **AND** the user successfully opens the console
- **THEN** the history file SHALL be created with the selected profile as the only entry

### Requirement: Recent profiles shown at top of console selector
When `bmc console` triggers interactive profile selection and the history file contains entries, the selector SHALL show recent profiles at the top of the list, visually distinguished from the full profile list.

#### Scenario: Recent profiles displayed at top
- **WHEN** the history file contains one or more entries
- **AND** interactive profile selection is triggered
- **THEN** recent profiles SHALL appear at the top of the list
- **AND** each recent profile SHALL show a "recent" label beside its name
- **AND** each profile SHALL appear only once in the list

#### Scenario: Empty history shows normal list
- **WHEN** the history file is empty or does not exist
- **THEN** the profile selector SHALL show the full profile list without a recent section

#### Scenario: Stale history entry for deleted profile
- **WHEN** a profile in the history no longer exists in `~/.aws/config`
- **THEN** it SHALL still appear in the recent section
- **AND** selecting it SHALL result in a clear error message from the profile lookup
