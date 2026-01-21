# MFA Authentication User Experience

## ADDED Requirements

### Requirement: Clear User-Facing Messages
The system SHALL display clear, user-friendly messages during MFA authentication that focus on user actions rather than technical implementation details.

#### Scenario: Source profile identification
- **WHEN** MFA authentication begins
- **THEN** the system SHALL display the message "-- Using AWS source-profile: {profile_name}"
- **AND** the system SHALL NOT display technical variable names like "sourceProfile"
- **AND** the message SHALL use consistent formatting with `--` prefix

#### Scenario: MFA session refresh indication
- **WHEN** an MFA session needs to be refreshed
- **THEN** the system SHALL display "-- Refreshing MFA session for {profile_name}..."
- **AND** the system SHALL NOT display the raw aws-mfa command with ARNs and flags
- **AND** the message SHALL indicate an action in progress

#### Scenario: TOTP script execution feedback
- **WHEN** a TOTP script is about to execute
- **THEN** the system SHALL display "-- Executing TOTP script..."
- **AND** the message SHALL appear before script execution begins
- **AND** users SHALL have feedback while waiting for password managers or TOTP tools

#### Scenario: Successful clipboard copy
- **WHEN** TOTP code is generated
- **AND** clipboard copy command succeeds
- **THEN** the system SHALL display the TOTP code
- **AND** the system SHALL display "-- Copied to clipboard"
- **AND** the success message SHALL only appear when copy actually succeeded

#### Scenario: Failed clipboard copy
- **WHEN** TOTP code is generated
- **AND** clipboard copy command fails or is not available
- **THEN** the system SHALL display the TOTP code
- **AND** the system SHALL display "-- Note: Clipboard copy failed (command not found or error)"
- **AND** the system SHALL suppress error output from the failed clipboard command

#### Scenario: Manual TOTP entry
- **WHEN** no TOTP script is configured
- **THEN** the system SHALL display "-- No TOTP script configured. Please enter MFA code manually."
- **AND** the message SHALL clearly guide the user on what to do next

### Requirement: No Debug Output in Production
The system SHALL NOT display debug information, internal state, or raw commands during normal operation.

#### Scenario: Script invocation
- **WHEN** BMC is executed or sourced
- **THEN** the system SHALL NOT display the script name or invocation path
- **AND** no `$0` or similar debug variables SHALL be echoed

#### Scenario: Internal state
- **WHEN** MFA configuration is being processed
- **THEN** the system SHALL NOT display boolean flags like "MFA: true"
- **AND** the system SHALL NOT display internal variable states
- **AND** configuration values SHALL only be shown when relevant to user actions

#### Scenario: Command execution
- **WHEN** AWS commands are executed
- **THEN** the system SHALL NOT echo the full command with all flags and ARNs
- **AND** only user-relevant status messages SHALL be displayed
- **AND** technical details SHALL be hidden unless errors occur

### Requirement: Consistent Message Formatting
The system SHALL use consistent formatting for all status messages to maintain a professional appearance.

#### Scenario: Status message prefix
- **WHEN** displaying informational messages
- **THEN** messages SHALL use the `--` prefix format
- **AND** messages SHALL clearly indicate the action or status
- **AND** formatting SHALL be consistent across all MFA operations

#### Scenario: Error and warning messages
- **WHEN** displaying errors
- **THEN** error messages SHALL use the `!!` prefix
- **AND** errors SHALL clearly explain what went wrong and what the user should do

#### Scenario: Note messages
- **WHEN** displaying non-critical information
- **THEN** note messages SHALL use "-- Note:" prefix
- **AND** notes SHALL provide helpful context without requiring user action
