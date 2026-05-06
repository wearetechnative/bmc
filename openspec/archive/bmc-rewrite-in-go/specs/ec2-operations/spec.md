## ADDED Requirements

### Requirement: List EC2 instances
The system SHALL list all EC2 instances in the current AWS account/region using AWS SDK Go v2 and display them in a table (bubbles/table) with columns: InstanceId, PrivateIpAddress, PublicIpAddress, State, Hibernate, Name, Scheduler.

#### Scenario: Instances present
- **WHEN** user runs `bmc ec2ls` and instances exist
- **THEN** system displays a paginated table of instances with all required columns

#### Scenario: No instances
- **WHEN** no EC2 instances exist in the region
- **THEN** system displays an empty table with headers

#### Scenario: Hibernate column normalisation
- **WHEN** an instance has `HibernationOptions.Configured = true`
- **THEN** the Hibernate column shows `yes`; otherwise `no`

#### Scenario: Scheduler column normalisation
- **WHEN** an instance has an `InstanceScheduler` tag with any value
- **THEN** the Scheduler column shows `yes`; otherwise `no`

### Requirement: Find EC2 instances across profile groups
The system SHALL search for EC2 instances matching a search string across all profiles in a user-selected profile group, running queries concurrently.

#### Scenario: Search with results
- **WHEN** user runs `bmc ec2find <search-string>` and matching instances are found
- **THEN** system displays matching instances with an additional Profile column

#### Scenario: No search string provided
- **WHEN** user runs `bmc ec2find` without a search string
- **THEN** system displays usage error and exits with non-zero code

### Requirement: Stop and start EC2 instances
The system SHALL allow stopping or starting a selected EC2 instance, with support for hibernate where enabled. A spinner (bubbles/spinner) SHALL indicate waiting for state transition.

#### Scenario: Start a stopped instance
- **WHEN** user selects a stopped instance via `bmc ec2stopstart`
- **THEN** system starts the instance and shows a spinner until state is `running`

#### Scenario: Stop a running instance
- **WHEN** user selects a running instance
- **THEN** system offers stop options: `stop` always; `hibernate` if hibernation is enabled

#### Scenario: Instance in non-actionable state
- **WHEN** selected instance is in `pending`, `stopping`, or other transitional state
- **THEN** system displays current state and exits without action

### Requirement: Toggle EC2 InstanceScheduler tag
The system SHALL toggle the `InstanceScheduler` tag on a selected EC2 instance to enable or disable scheduling.

#### Scenario: Enable scheduling
- **WHEN** instance has no `InstanceScheduler` tag and user confirms enable
- **THEN** system adds the `InstanceScheduler` tag

#### Scenario: Disable scheduling
- **WHEN** instance has an `InstanceScheduler` tag and user confirms disable
- **THEN** system removes the `InstanceScheduler` tag
