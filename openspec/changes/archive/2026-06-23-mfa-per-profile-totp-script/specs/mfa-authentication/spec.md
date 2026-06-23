## MODIFIED Requirements

### Requirement: TOTP Script Execution with Arguments
The system SHALL resolve the TOTP script to execute by first checking `mfa.profile_scripts[sourceProfile]`, then falling back to `mfa.totp_script`, before executing the resolved command. The script SHALL be executed as a shell command using `sh -c`.

#### Scenario: TOTP script resolved from profile_scripts
- **WHEN** `mfa.profile_scripts` contains an entry for the active source profile
- **AND** MFA session renewal is triggered
- **THEN** the system SHALL execute the profile-specific script with `sh -c`
- **AND** the system SHALL capture the TOTP code from the script output
- **AND** the TOTP code SHALL be displayed to the user

#### Scenario: TOTP script resolved from global totp_script fallback
- **WHEN** `mfa.profile_scripts` does not contain an entry for the active source profile
- **AND** `mfa.totp_script` is configured
- **AND** MFA session renewal is triggered
- **THEN** the system SHALL execute the global `totp_script` with `sh -c`
- **AND** the system SHALL capture the TOTP code from the script output
- **AND** the TOTP code SHALL be displayed to the user

#### Scenario: TOTP script with path containing spaces
- **WHEN** the resolved totp script contains a path with spaces
- **AND** MFA session renewal is triggered
- **THEN** the system SHALL correctly handle the path with spaces
- **AND** the system SHALL execute the script successfully

#### Scenario: Simple TOTP script without arguments
- **WHEN** the resolved totp script is a simple path with no arguments
- **AND** MFA session renewal is triggered
- **THEN** the system SHALL execute the script
- **AND** the system SHALL capture and display the TOTP code
