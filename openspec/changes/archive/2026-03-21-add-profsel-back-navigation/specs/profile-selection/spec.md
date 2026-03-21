# profile-selection Spec Delta

## ADDED Requirements

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

## MODIFIED Requirements
None - this change adds new functionality without modifying existing requirements.

## REMOVED Requirements
None - this change extends existing functionality without removing capabilities.
