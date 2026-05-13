## MODIFIED Requirements

### Requirement: Clipboard Integration for TOTP Codes
The system SHALL automatically copy generated TOTP codes to the clipboard using the configured `copy_command`, and optionally simulate a paste keystroke using `paste_command`.

#### Scenario: Clipboard copy with configured command
- **WHEN** `mfa.copy_command` is configured in `~/.config/bmc/config.json`
- **AND** `totpScript` successfully generates a TOTP code
- **THEN** the system SHALL copy the TOTP code to clipboard using the configured command
- **AND** the system SHALL display a confirmation message "-- Copied to clipboard"
- **AND** the system SHALL also display the TOTP code for manual reference

#### Scenario: Clipboard command with arguments
- **WHEN** the config file contains `"copy_command": "xclip -selection clipboard"`
- **AND** a TOTP code is generated
- **THEN** the system SHALL execute the copy command with all arguments
- **AND** the TOTP code SHALL be copied to the system clipboard

#### Scenario: Clipboard copy without configured command
- **WHEN** `mfa.copy_command` is not defined in `~/.config/bmc/config.json`
- **AND** `totpScript` successfully generates a TOTP code
- **THEN** the system SHALL display the TOTP code
- **AND** the system SHALL NOT attempt to copy to clipboard
- **AND** the system SHALL NOT display an error about missing clipboard command

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
