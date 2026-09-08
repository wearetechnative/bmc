## Purpose

Defines how users select an AWS profile — the interactive picker's navigation behavior and the command-line flags (`-p/--profile`, `-P/--pick`) shared by every command that needs a profile.

## Requirements

### Requirement: Back Navigation in Profile Selection
The `selectProfileInteractive` function SHALL allow users to return to the group selection menu when they press ESC at the profile selection stage, without exiting the entire command. Pressing Ctrl+C at any level SHALL cancel the entire selection.

#### Scenario: User presses ESC at profile selection and returns to group menu
- **GIVEN** the user has selected a profile group
- **WHEN** the user presses ESC at the profile selection stage
- **THEN** the system SHALL return to the group selection menu
- **AND** allow the user to select a different group or cancel entirely

#### Scenario: User presses Ctrl+C at profile selection to exit
- **GIVEN** the user has selected a profile group
- **WHEN** the user presses Ctrl+C at the profile selection stage
- **THEN** the system SHALL exit the selection process with no profile selected

#### Scenario: User cancels at group selection to exit
- **GIVEN** the user is at the group selection menu
- **WHEN** the user presses ESC or Ctrl+C at the group selection stage
- **THEN** the system SHALL exit the selection process gracefully with no profile selected

#### Scenario: User successfully selects profile after returning to group menu
- **GIVEN** the user pressed ESC at profile selection and returned to group menu
- **WHEN** the user selects a different group and then selects a valid profile
- **THEN** the system SHALL proceed with the selected profile normally

#### Scenario: Navigation loop with preferred profile flag
- **GIVEN** the user runs profsel with `-p <profile-name>` flag
- **WHEN** the selectProfileInteractive function is called
- **THEN** the system SHALL bypass both group and profile selection menus
- **AND** SHALL use the specified profile directly without entering the navigation loop

### Requirement: Profile Selection Command-Line Flags
Every command that selects an AWS profile SHALL expose two distinct, single-purpose flags:

- `-p` / `--profile <NAME>` — SHALL accept a profile name as its value. Both the space form (`-p NAME`) and the equals form (`-p=NAME`) SHALL be accepted. The flag SHALL NOT consume or interfere with positional arguments such as an instance search term.
- `-P` / `--pick` — a boolean flag that SHALL force the interactive profile picker, ignoring any `AWS_PROFILE` environment variable.

When neither flag is provided, the command SHALL use `AWS_PROFILE` if it is set, and otherwise fall back to the interactive picker.

Bare `-p` with no value SHALL NOT be a valid way to request interactive selection; forcing interactive selection is done exclusively via `-P`/`--pick`.

This contract SHALL apply uniformly to: `ec2`, `ec2connect`, `ec2ls`, `ec2scheduler`, `ec2stopstart`, `ecsconnect`, `console`, and `profsel`.

#### Scenario: Profile name and positional search argument together
- **WHEN** the user runs `bmc ec2connect -p TN-Production compute2`
- **THEN** the system SHALL resolve the profile to `TN-Production`
- **AND** SHALL treat `compute2` as the instance search term
- **AND** SHALL NOT report an "accepts at most 1 arg(s)" error

#### Scenario: Equals form for profile value
- **WHEN** the user runs `bmc ec2connect -p=TN-Production compute2`
- **THEN** the system SHALL resolve the profile to `TN-Production`
- **AND** SHALL treat `compute2` as the instance search term

#### Scenario: Force interactive picker with a positional argument
- **WHEN** the user runs `bmc ec2connect --pick compute2`
- **THEN** the system SHALL open the interactive profile picker, ignoring `AWS_PROFILE`
- **AND** after a profile is chosen SHALL treat `compute2` as the instance search term

#### Scenario: Force interactive picker via short flag
- **WHEN** the user runs a profile-taking command with `-P`
- **THEN** the system SHALL open the interactive profile picker, ignoring `AWS_PROFILE`

#### Scenario: No profile flag falls back to environment then picker
- **GIVEN** `AWS_PROFILE` is set to a valid profile
- **WHEN** the user runs a profile-taking command with neither `-p` nor `-P`
- **THEN** the system SHALL use the `AWS_PROFILE` value
- **AND WHEN** `AWS_PROFILE` is not set
- **THEN** the system SHALL open the interactive picker

#### Scenario: Bare -p no longer triggers interactive selection
- **WHEN** the user runs a profile-taking command with a bare `-p` and no value
- **THEN** the system SHALL report that the flag needs an argument
- **AND** SHALL NOT open the interactive picker

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
