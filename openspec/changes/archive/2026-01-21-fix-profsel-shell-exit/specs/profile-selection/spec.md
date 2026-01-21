# Profile Selection

## ADDED Requirements

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
