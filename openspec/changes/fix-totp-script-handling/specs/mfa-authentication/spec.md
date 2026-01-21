# MFA Authentication with TOTP Script Integration

## MODIFIED Requirements

### Requirement: TOTP Script Execution with Arguments
The system SHALL properly execute TOTP scripts configured as arrays with command-line arguments in the config file.

#### Scenario: TOTP script with multiple arguments
- **WHEN** the config file contains `totpScript=("/path/to/script" "-t" "code" "-q" "new")`
- **AND** MFA session renewal is triggered
- **THEN** the system SHALL execute the script with all provided arguments
- **AND** the system SHALL capture the TOTP code from the script output
- **AND** the TOTP code SHALL be displayed to the user

#### Scenario: TOTP script with path containing spaces
- **WHEN** the config file contains `totpScript=("/path with spaces/script.sh" "--flag")`
- **AND** MFA session renewal is triggered
- **THEN** the system SHALL correctly handle the path with spaces
- **AND** the system SHALL execute the script successfully

#### Scenario: Simple TOTP script without arguments
- **WHEN** the config file contains `totpScript=("/path/to/simple-totp.sh")`
- **AND** MFA session renewal is triggered
- **THEN** the system SHALL execute the script
- **AND** the system SHALL capture and display the TOTP code

### Requirement: Clipboard Integration for TOTP Codes
The system SHALL automatically copy generated TOTP codes to the clipboard using the configured clipboard command.

#### Scenario: Clipboard copy with configured command
- **WHEN** `clipboardCopyCommand` is configured in config.env
- **AND** `totpScript` successfully generates a TOTP code
- **THEN** the system SHALL copy the TOTP code to clipboard using the configured command
- **AND** the system SHALL display a confirmation message "-- Copied to clipboard"
- **AND** the system SHALL also display the TOTP code for manual reference

#### Scenario: Clipboard command with arguments
- **WHEN** the config file contains `clipboardCopyCommand=("/usr/bin/xclip" "-selection" "clipboard")`
- **AND** a TOTP code is generated
- **THEN** the system SHALL execute the clipboard command with all arguments
- **AND** the TOTP code SHALL be copied to the system clipboard

#### Scenario: Clipboard copy without configured command
- **WHEN** `clipboardCopyCommand` is not defined in config.env
- **AND** `totpScript` successfully generates a TOTP code
- **THEN** the system SHALL display the TOTP code
- **AND** the system SHALL NOT attempt to copy to clipboard
- **AND** the system SHALL NOT display an error about missing clipboard command

### Requirement: Clear User Feedback for TOTP Configuration
The system SHALL provide clear feedback to users based on their TOTP script configuration status.

#### Scenario: TOTP script not configured
- **WHEN** `totpScript` is not defined or empty in config.env
- **AND** MFA authentication is required
- **THEN** the system SHALL display a helpful message indicating manual MFA entry is needed
- **AND** the message SHALL be: "-- No TOTP script configured. Please enter MFA code manually."
- **AND** the system SHALL NOT display undefined or empty variables

#### Scenario: TOTP script configured and executed successfully
- **WHEN** `totpScript` is configured
- **AND** the script executes successfully and returns a TOTP code
- **THEN** the system SHALL display the generated code
- **AND** if clipboard is configured, SHALL copy the code and confirm
- **AND** the system SHALL proceed with MFA authentication using the code

#### Scenario: TOTP script execution fails
- **WHEN** `totpScript` is configured
- **AND** the script fails to execute or returns an error
- **THEN** the system SHALL display an error message
- **AND** the system SHALL allow the user to enter MFA code manually
- **AND** the MFA authentication flow SHALL continue

### Requirement: Configuration Documentation
The system documentation SHALL provide clear examples of TOTP script and clipboard command configuration using array syntax.

#### Scenario: Configuration file examples
- **WHEN** users consult BMC documentation
- **THEN** documentation SHALL include example configurations for totpScript using array syntax
- **AND** documentation SHALL include example configurations for clipboardCopyCommand
- **AND** documentation SHALL explain how to configure scripts with arguments
- **AND** documentation SHALL provide examples for common TOTP tools

#### Scenario: Backward compatibility
- **WHEN** existing users have array-based totpScript configurations
- **THEN** the system SHALL continue to work with their existing configurations
- **AND** no configuration migration SHALL be required
- **AND** behavior SHALL be identical to intended functionality
