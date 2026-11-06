# ec2ls Specification

## Purpose
TBD - created by archiving change enhance-instance-list-scheduler-status. Update Purpose after archive.

## Requirements

### Requirement: List All EC2 Instances
The `bmc ec2ls` command SHALL list all EC2 instances in the current AWS profile with formatted instance details displayed in an interactive table, correctly parsing field values that contain spaces.

#### Scenario: Display instance information columns
- **WHEN** displaying the instance table
- **THEN** the command SHALL show at minimum: instance ID, private IP address, public IP address, state, hibernation status, instance name (from Name tag), and scheduler configuration status
- **AND** the command SHALL correctly parse and display instance information including Name tags that contain spaces
- **AND** the command SHALL use tab character as the field separator for awk processing
- **AND** the command SHALL preserve tabs during variable expansion by quoting variables
- **AND** the command SHALL correctly extract all fields regardless of spaces within field values

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
The `bmc ec2ls` command SHALL present instance information in a bordered table on stdout for easy reading, using the columns from the `ec2.columns` configuration.

#### Scenario: Format output as table
- **WHEN** displaying instance information without `--json` or `--txt`
- **THEN** the command SHALL render a table with border lines, a bold header row and one row per instance
- **AND** each column SHALL be as wide as its widest value

#### Scenario: Cell padding
- **WHEN** the table is rendered
- **THEN** every cell, including header cells, SHALL have one space between the left border and its content
- **AND** one space between its content and the right border

#### Scenario: Table written to stdout
- **WHEN** the table is rendered
- **THEN** the table SHALL be written to stdout so it can be piped to other commands

#### Scenario: Display CSV header
- **WHEN** building the table output
- **THEN** the command SHALL include a header row with the column names from `ec2.columns`
- **AND** the header SHALL NOT be emitted as a comma-separated line
