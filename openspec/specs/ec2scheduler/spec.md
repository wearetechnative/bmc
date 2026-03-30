# ec2scheduler Specification

## Purpose
The `bmc ec2scheduler` command provides convenient management of EC2 instance scheduler overrides using the `Ignore_scheduler` tag. This allows users to temporarily extend an instance's runtime beyond scheduled stop times without completely disabling the scheduler. The override is time-based and automatically expires, returning the instance to its normal schedule.
## Requirements
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

### Requirement: Interactive Instance Selection
The `bmc ec2scheduler` command SHALL use an interactive table interface for users to select which instance to manage.

#### Scenario: Select instance from table
- **WHEN** user runs `bmc ec2scheduler` and multiple instances are available
- **THEN** the command SHALL present a `gum` table showing instance details including instance ID, name, current state, and Ignore_scheduler status

#### Scenario: Extract instance region for console URL
- **WHEN** an instance is selected
- **THEN** the command SHALL determine the AWS region from the instance availability zone
- **AND** the region SHALL be used to construct region-specific console URLs

### Requirement: Interactive Action Menu
The `bmc ec2scheduler` command SHALL present an action menu after instance selection, allowing users to set or remove scheduler override.

#### Scenario: Display action menu
- **WHEN** user selects an instance from the table
- **THEN** the command SHALL display the current Ignore_scheduler status (if set)
- **AND** the command SHALL present a menu with options: "Set ignore until time", "Remove ignore override", "Cancel"

#### Scenario: User cancels from menu
- **WHEN** user selects "Cancel" from the action menu
- **THEN** the command SHALL exit without making any changes
- **AND** the command SHALL display a cancellation message

### Requirement: Set Ignore_scheduler Override
The `bmc ec2scheduler` command SHALL allow users to set an Ignore_scheduler tag with a time and timezone value.

#### Scenario: User sets ignore override time
- **WHEN** user selects "Set ignore until time" from the action menu
- **THEN** the command SHALL prompt for time in HH:MM format
- **AND** the command SHALL show an example: "Example: 22:00"
- **AND** the command SHALL validate the time matches HH:MM format (24-hour)

#### Scenario: User provides timezone
- **WHEN** user has provided a valid time
- **THEN** the command SHALL prompt for timezone
- **AND** the command SHALL show examples: "Examples: Europe/Amsterdam, America/New_York, UTC"
- **AND** the command SHALL accept free-form text input

#### Scenario: Create Ignore_scheduler tag
- **WHEN** user has provided both time and timezone
- **THEN** the command SHALL combine them into format: "HH:MM Timezone"
- **AND** the command SHALL create or update the `Ignore_scheduler` tag on the instance with this value
- **AND** the command SHALL display a success message showing the instance ID and ignore-until value

#### Scenario: Update existing Ignore_scheduler tag
- **WHEN** an instance already has an `Ignore_scheduler` tag
- **AND** user sets a new ignore override time
- **THEN** the command SHALL update the existing tag value with the new time and timezone

### Requirement: Remove Ignore_scheduler Override
The `bmc ec2scheduler` command SHALL allow users to remove an existing Ignore_scheduler tag.

#### Scenario: Remove existing ignore override
- **WHEN** user selects "Remove ignore override" from the action menu
- **AND** the instance has an `Ignore_scheduler` tag
- **THEN** the command SHALL delete the `Ignore_scheduler` tag from the instance
- **AND** the command SHALL display a success message confirming removal

#### Scenario: Attempt to remove non-existent override
- **WHEN** user selects "Remove ignore override" from the action menu
- **AND** the instance does not have an `Ignore_scheduler` tag
- **THEN** the command SHALL display a message: "No ignore override is currently set"
- **AND** the command SHALL exit without making changes

### Requirement: Validate Time Format
The `bmc ec2scheduler` command SHALL validate time input before creating tags.

#### Scenario: Valid time format
- **WHEN** user enters time in HH:MM format (e.g., "22:00", "08:30", "16:45")
- **THEN** the command SHALL accept the input and proceed to timezone prompt

#### Scenario: Invalid time format
- **WHEN** user enters time in incorrect format (e.g., "22", "10:00 PM", "25:00")
- **THEN** the command SHALL reject the input
- **AND** the command SHALL display an error message explaining the correct format
- **AND** the command SHALL re-prompt for time input

### Requirement: Handle Instances Without Scheduler Tags
The `bmc ec2scheduler` command SHALL provide helpful guidance when users select instances without the base InstanceScheduler tag configured.

#### Scenario: User selects instance without InstanceScheduler tag
- **WHEN** user selects an instance that does not have the `InstanceScheduler` tag
- **THEN** the command SHALL display a message indicating the instance does not have scheduler configuration
- **AND** the command SHALL prompt the user if they want to add the InstanceScheduler tag
- **AND** if confirmed, SHALL offer to open AWS Console to the instance details page

#### Scenario: User accepts opening AWS Console for tag configuration
- **WHEN** user confirms they want to add the InstanceScheduler tag via AWS Console
- **THEN** the command SHALL use `assumego` to open the AWS Console directly to the instance details page
- **AND** the console URL SHALL be in format: `https://<region>.console.aws.amazon.com/ec2/home?region=<region>#InstanceDetails:instanceId=<instance-id>`
- **AND** the command SHALL execute with environment variables `GRANTED_ALIAS_CONFIGURED="true"` and `GRANTED_ENABLE_AUTO_REASSUME=true`

#### Scenario: User declines opening AWS Console
- **WHEN** user is shown instructions to add InstanceScheduler tag manually
- **AND** user declines to open the AWS Console
- **THEN** the command SHALL exit gracefully
- **AND** the command SHALL display the instructions for future reference

#### Scenario: User declines adding scheduler tag
- **WHEN** user selects an instance without scheduler tags
- **AND** user declines to add the scheduler tag
- **THEN** the command SHALL exit without making any changes or showing console instructions

### Requirement: Provide Clear User Feedback
The `bmc ec2scheduler` command SHALL provide clear feedback about ignore override operations.

#### Scenario: Confirm successful override creation
- **WHEN** an Ignore_scheduler tag is successfully created or updated
- **THEN** the command SHALL display a message showing the instance ID and the time until which it will ignore scheduled stops

#### Scenario: Confirm successful override removal
- **WHEN** an Ignore_scheduler tag is successfully removed
- **THEN** the command SHALL display a message confirming removal and indicating the instance will resume normal schedule

#### Scenario: Show error on operation failure
- **WHEN** a tag operation fails (e.g., due to AWS API errors or permission issues)
- **THEN** the command SHALL display a clear error message explaining what went wrong
- **AND** the command SHALL exit with a non-zero status code

