## ADDED Requirements

### Requirement: Customer group selection
The system SHALL present an interactive list of all available customer profile groups using `gum filter`, allowing the user to select one.

#### Scenario: User selects a customer group
- **WHEN** the user runs `bmc accountls`
- **THEN** the system SHALL display all unique profile groups from `jsonify-aws-dotfiles` output in a `gum filter` prompt

#### Scenario: User cancels group selection
- **WHEN** the user presses Escape or Ctrl-C during group selection
- **THEN** the command SHALL exit with a non-zero exit code without printing any output

### Requirement: Account listing display
After the user selects a customer group, the system SHALL display a table of all AWS profiles belonging to that group.

#### Scenario: Customer has multiple profiles
- **WHEN** the user selects a customer group that contains profiles
- **THEN** the system SHALL display a table with columns: Profile Name, Account ID, Role Name
- **AND** the table SHALL be formatted using `gum table`
- **AND** the command SHALL exit after displaying the table

#### Scenario: Customer has no profiles
- **WHEN** the user selects a customer group that has no matching profiles
- **THEN** the system SHALL display a message indicating no accounts were found
- **AND** the command SHALL exit with a zero exit code

### Requirement: Command registration
The `accountls` subcommand SHALL be registered in the `bmc` dispatcher with a description.

#### Scenario: Command appears in usage
- **WHEN** the user runs `bmc` or `bmc usage`
- **THEN** `accountls` SHALL appear in the list of available commands with a description
