# EC2 Connection with Automatic Instance Start

## ADDED Requirements

### Requirement: Prompt to Start Stopped Instances
When a user selects a stopped EC2 instance in the ec2connect command, the system SHALL offer to start the instance before attempting connection, unless configured otherwise.

#### Scenario: User confirms starting stopped instance
- **WHEN** the user selects a stopped EC2 instance
- **AND** the config option `BMC_AUTO_START_STOPPED_INSTANCES` is unset or set to "prompt"
- **AND** the user confirms the prompt to start the instance
- **THEN** the system SHALL start the instance using `aws ec2 start-instances`
- **AND** the system SHALL wait until the instance reaches "running" state
- **AND** the system SHALL proceed with the normal connection flow

#### Scenario: User declines starting stopped instance
- **WHEN** the user selects a stopped EC2 instance
- **AND** the config option `BMC_AUTO_START_STOPPED_INSTANCES` is unset or set to "prompt"
- **AND** the user declines the prompt to start the instance
- **THEN** the system SHALL exit gracefully without starting the instance
- **AND** the system SHALL NOT attempt to connect to the instance

#### Scenario: Auto-start mode enabled
- **WHEN** the user selects a stopped EC2 instance
- **AND** the config option `BMC_AUTO_START_STOPPED_INSTANCES` is set to "always"
- **THEN** the system SHALL automatically start the instance without prompting
- **AND** the system SHALL wait until the instance reaches "running" state
- **AND** the system SHALL proceed with the normal connection flow

#### Scenario: Auto-start mode disabled
- **WHEN** the user selects a stopped EC2 instance
- **AND** the config option `BMC_AUTO_START_STOPPED_INSTANCES` is set to "never"
- **THEN** the system SHALL display an error message about the stopped state
- **AND** the system SHALL exit without starting the instance
- **AND** the system SHALL NOT prompt the user

#### Scenario: Running instance unchanged
- **WHEN** the user selects a running EC2 instance
- **THEN** the system SHALL proceed with connection immediately
- **AND** the system SHALL NOT prompt to start the instance
- **AND** behavior SHALL be identical to pre-change functionality

### Requirement: Configuration Option for Auto-Start Behavior
The system SHALL support a configuration option `BMC_AUTO_START_STOPPED_INSTANCES` in the BMC config file to control automatic instance starting behavior.

#### Scenario: Config file with auto-start always
- **WHEN** the config file `~/.config/bmc/config.env` contains `BMC_AUTO_START_STOPPED_INSTANCES="always"`
- **THEN** the system SHALL automatically start stopped instances without prompting
- **AND** the system SHALL wait for instances to reach running state before connecting

#### Scenario: Config file with auto-start never
- **WHEN** the config file `~/.config/bmc/config.env` contains `BMC_AUTO_START_STOPPED_INSTANCES="never"`
- **THEN** the system SHALL never start stopped instances
- **AND** the system SHALL exit with an error message when a stopped instance is selected

#### Scenario: Config file with prompt mode
- **WHEN** the config file `~/.config/bmc/config.env` contains `BMC_AUTO_START_STOPPED_INSTANCES="prompt"`
- **THEN** the system SHALL prompt the user before starting stopped instances
- **AND** the system SHALL respect the user's response

#### Scenario: Config option unset or missing
- **WHEN** the config option `BMC_AUTO_START_STOPPED_INSTANCES` is not present in the config file
- **THEN** the system SHALL default to "prompt" behavior
- **AND** the system SHALL ask the user before starting stopped instances

### Requirement: Non-Stopped State Handling
The system SHALL maintain existing error handling for EC2 instances in states other than "stopped" or "running".

#### Scenario: Instance in pending state
- **WHEN** the user selects an EC2 instance in "pending" state
- **THEN** the system SHALL display an error message about the current state
- **AND** the system SHALL exit without attempting to start or connect

#### Scenario: Instance in stopping state
- **WHEN** the user selects an EC2 instance in "stopping" state
- **THEN** the system SHALL display an error message about the current state
- **AND** the system SHALL exit without attempting to start or connect

#### Scenario: Instance in terminated state
- **WHEN** the user attempts to select an EC2 instance in "terminated" state
- **THEN** the system SHALL NOT display the instance in the selection list
- **AND** behavior SHALL be identical to pre-change functionality

### Requirement: Improved Error Messages
The system SHALL provide clear, concise error messages without redundant connection attempt information.

#### Scenario: Error message clarity
- **WHEN** the system displays an error for a non-running instance
- **THEN** the error message SHALL state the instance ID and current state
- **AND** the error message SHALL NOT include "Not executing the SSH-command"
- **AND** the message format SHALL be: "!!! Instance chosen is not running. Current state is : {state}."
