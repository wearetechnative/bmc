# Spec Delta

## MODIFIED Requirements

### Requirement: Clear User-Facing Messages
The system SHALL display clear, user-friendly messages during MFA authentication that focus on user actions rather than technical implementation details. Messages that only report state (which source profile is used, that the current session is still valid, the TOTP code produced by a script) are diagnostic and SHALL only be displayed when `--debug` is given. Messages that inform about or require user action SHALL always be displayed. All MFA messages SHALL be written to stderr.

#### Scenario: Source profile identification
- **WHEN** MFA authentication begins and `--debug` is given
- **THEN** the system SHALL display the message "-- Using AWS source-profile: {profile_name}"
- **AND** the system SHALL NOT display technical variable names like "sourceProfile"
- **AND** the message SHALL use consistent formatting with `--` prefix

#### Scenario: Source profile hidden by default
- **WHEN** MFA authentication begins without `--debug`
- **THEN** the system SHALL NOT display the source-profile message

#### Scenario: Valid session reported only in debug
- **WHEN** the MFA session is still valid
- **THEN** the system SHALL display "Current MFA Session Valid, until: {expiration}" only when `--debug` is given
- **AND** without `--debug` the system SHALL display no MFA message at all

#### Scenario: MFA session refresh indication
- **WHEN** an MFA session needs to be refreshed
- **THEN** the system SHALL display "-- Refreshing MFA session for {profile_name}..." regardless of `--debug`
- **AND** the system SHALL NOT display the raw aws-mfa command with ARNs and flags
- **AND** the message SHALL indicate an action in progress

#### Scenario: TOTP script execution feedback
- **WHEN** a TOTP script is about to execute
- **THEN** the system SHALL display "-- Executing TOTP script..." regardless of `--debug`
- **AND** the message SHALL appear before script execution begins
- **AND** users SHALL have feedback while waiting for password managers or TOTP tools

#### Scenario: TOTP code echo only in debug
- **WHEN** a TOTP script has produced a code
- **THEN** the system SHALL display the code only when `--debug` is given

#### Scenario: Manual TOTP entry
- **WHEN** no TOTP script is configured
- **THEN** the system SHALL display "-- No TOTP script configured. Please enter MFA code manually." regardless of `--debug`
- **AND** the message SHALL clearly guide the user on what to do next

#### Scenario: Refresh confirmation
- **WHEN** an MFA session has been refreshed successfully
- **THEN** the system SHALL display "-- MFA session refreshed, valid until: {expiration}" regardless of `--debug`
