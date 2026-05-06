## MODIFIED Requirements

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
