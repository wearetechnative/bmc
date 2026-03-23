# ec2scheduler Spec Delta

## MODIFIED Requirements

### Requirement: List All EC2 Instances
The `bmc ec2scheduler` command SHALL list all EC2 instances in an interactive table, showing their scheduler configuration status and Ignore_scheduler status.

#### Scenario: Display scheduler configuration status
- **WHEN** user runs `bmc ec2scheduler`
- **THEN** the command SHALL display a "Scheduler" column showing whether the InstanceScheduler tag is configured
- **AND** the Scheduler column SHALL display "yes" if the InstanceScheduler tag exists with any value
- **AND** the Scheduler column SHALL display "no" if the InstanceScheduler tag is missing

#### Scenario: Display ignore override information
- **WHEN** displaying the instance table
- **THEN** the command SHALL show at minimum: instance ID, instance name (from Name tag), instance state, scheduler configuration status (yes/no), and the Ignore_scheduler value (time until which instance will ignore scheduled stops)

#### Scenario: List instances with scheduler tag
- **WHEN** user runs `bmc ec2scheduler` and instances exist with the `InstanceScheduler` tag
- **THEN** the Scheduler column SHALL display "yes" for those instances

#### Scenario: List instances without scheduler tag
- **WHEN** user runs `bmc ec2scheduler` and instances exist without the `InstanceScheduler` tag
- **THEN** the Scheduler column SHALL display "no" for those instances
