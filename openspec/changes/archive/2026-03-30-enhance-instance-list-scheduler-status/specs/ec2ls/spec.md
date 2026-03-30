# ec2ls Spec

## ADDED Requirements

### Requirement: List All EC2 Instances
The `bmc ec2ls` command SHALL list all EC2 instances in the current AWS profile with formatted instance details displayed in an interactive table.

#### Scenario: List all instances
- **WHEN** user runs `bmc ec2ls`
- **THEN** the command SHALL query AWS EC2 for all instances in the current AWS profile
- **AND** the command SHALL display the results in a formatted table

#### Scenario: Display instance information columns
- **WHEN** displaying the instance table
- **THEN** the command SHALL show at minimum: instance ID, private IP address, public IP address, state, hibernation status, instance name (from Name tag), and scheduler configuration status

#### Scenario: No AWS profile set
- **WHEN** user runs `bmc ec2ls` without AWS_PROFILE set
- **THEN** the command SHALL invoke profile selection flow
- **AND** after profile selection, SHALL proceed with instance listing

### Requirement: Display Hibernation Status
The `bmc ec2ls` command SHALL display hibernation configuration in a user-friendly format.

#### Scenario: Instance with hibernation enabled
- **WHEN** an instance has hibernation enabled (HibernationOptions.Configured = true)
- **THEN** the Hibernate column SHALL display "yes"

#### Scenario: Instance with hibernation disabled
- **WHEN** an instance has hibernation disabled (HibernationOptions.Configured = false)
- **THEN** the Hibernate column SHALL display "no"

#### Scenario: Instance with missing hibernation configuration
- **WHEN** an instance has no hibernation configuration data (None or null)
- **THEN** the Hibernate column SHALL display "no"

### Requirement: Display Scheduler Configuration Status
The `bmc ec2ls` command SHALL display whether instances are managed by the EC2 scheduler.

#### Scenario: Instance with scheduler tag
- **WHEN** an instance has the `InstanceScheduler` tag (regardless of value)
- **THEN** the Scheduler column SHALL display "yes"

#### Scenario: Instance without scheduler tag
- **WHEN** an instance does not have the `InstanceScheduler` tag
- **THEN** the Scheduler column SHALL display "no"

### Requirement: Interactive Table Display
The `bmc ec2ls` command SHALL present instance information in an interactive table for easy reading.

#### Scenario: Format output as table
- **WHEN** displaying instance information
- **THEN** the command SHALL use `gum table` to render a formatted table
- **AND** the command SHALL use proportional column widths for readability

#### Scenario: Display CSV header
- **WHEN** building the table output
- **THEN** the command SHALL include a CSV header row with column names
- **AND** the header SHALL be: "InstanceId,PrivateIpAddress,PublicIpAddress,State,Hibernate,Name,Scheduler"
