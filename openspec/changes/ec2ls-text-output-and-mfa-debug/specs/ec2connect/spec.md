# Spec Delta

## MODIFIED Requirements

### Requirement: Automatic SSH Connection Selection
When SSH-specific command-line flags are provided to the ec2connect command, the system SHALL automatically select SSH as the connection method without prompting the user to choose between SSH and SSM.

#### Scenario: User provides username flag
- **WHEN** the user runs `bmc ec2connect -u <username>`
- **THEN** the system SHALL automatically select SSH connection method
- **AND** the system SHALL NOT prompt for connection type selection

#### Scenario: User provides identity file flag
- **WHEN** the user runs `bmc ec2connect -i <path_to_key>`
- **THEN** the system SHALL automatically select SSH connection method
- **AND** the system SHALL NOT prompt for connection type selection

#### Scenario: User provides both username and identity file flags
- **WHEN** the user runs `bmc ec2connect -u <username> -i <path_to_key>`
- **THEN** the system SHALL automatically select SSH connection method
- **AND** the system SHALL NOT prompt for connection type selection

#### Scenario: User provides no SSH-specific flags
- **WHEN** the user runs `bmc ec2connect` without -u or -i flags
- **THEN** the system SHALL prompt the user to choose between SSH and SSM connection methods
- **AND** SSM SHALL be listed first and be the default selection
- **AND** the system SHALL respect the user's manual selection
