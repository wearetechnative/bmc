# ec2scheduler Specification

## Purpose
TBD - created by archiving change add-ec2-scheduler-toggle. Update Purpose after archive.
## Requirements
### Requirement: List All EC2 Instances
The `bmc ec2scheduler` command SHALL list all EC2 instances in an interactive table, showing their scheduler status.

#### Scenario: List instances with InstanceScheduler tag
- **WHEN** user runs `bmc ec2scheduler` and instances exist with the `InstanceScheduler` tag
- **THEN** the command SHALL display those instances with their current scheduler status showing as "enabled" and the schedule value

#### Scenario: List instances with InstanceScheduler_DISABLED tag
- **WHEN** user runs `bmc ec2scheduler` and instances exist with the `InstanceScheduler_DISABLED` tag
- **THEN** the command SHALL display those instances with their current scheduler status showing as "disabled" and the schedule value

#### Scenario: List instances without scheduler tags
- **WHEN** user runs `bmc ec2scheduler` and instances exist without either `InstanceScheduler` or `InstanceScheduler_DISABLED` tags
- **THEN** the command SHALL display those instances with their scheduler status showing as "none" or "not configured"
- **AND** the scheduler value column SHALL be empty or show "N/A"

#### Scenario: List all instances regardless of scheduler tag presence
- **WHEN** user runs `bmc ec2scheduler`
- **THEN** the command SHALL display all non-terminated EC2 instances in the current AWS profile
- **AND** the command SHALL show scheduler status for each instance (enabled, disabled, or none)

#### Scenario: No instances exist
- **WHEN** user runs `bmc ec2scheduler` and no EC2 instances exist in the current AWS profile
- **THEN** the command SHALL display a message indicating no instances were found and exit gracefully

### Requirement: Interactive Instance Selection
The `bmc ec2scheduler` command SHALL use an interactive table interface for users to select which instance's scheduler tag to toggle.

#### Scenario: Select instance from table
- **WHEN** user runs `bmc ec2scheduler` and multiple instances are available
- **THEN** the command SHALL present a `gum` table showing instance details including instance ID, name, current state, and scheduler status

#### Scenario: Display relevant instance information
- **WHEN** displaying the instance table
- **THEN** the command SHALL show at minimum: instance ID, instance name (from Name tag), instance state, current scheduler status (enabled/disabled/none), and the scheduler value (schedule definition)

#### Scenario: Extract instance region for console URL
- **WHEN** an instance is selected
- **THEN** the command SHALL determine the AWS region from the instance availability zone
- **AND** the region SHALL be used to construct region-specific console URLs

### Requirement: Confirmation Before Toggle
The `bmc ec2scheduler` command SHALL request user confirmation before toggling the scheduler tag, displaying current tag details.

#### Scenario: Display current tag information
- **WHEN** user selects an instance to toggle
- **THEN** the command SHALL display the instance ID, current tag name, and current tag value
- **AND** the command SHALL prompt for confirmation before making changes

#### Scenario: User confirms toggle
- **WHEN** confirmation prompt is displayed
- **AND** user confirms the action
- **THEN** the command SHALL proceed with the tag toggle operation

#### Scenario: User cancels toggle
- **WHEN** confirmation prompt is displayed
- **AND** user cancels the action
- **THEN** the command SHALL exit without making any changes
- **AND** the command SHALL display a cancellation message

### Requirement: Toggle Scheduler Tag Name
The `bmc ec2scheduler` command SHALL toggle the scheduler tag name between `InstanceScheduler` and `InstanceScheduler_DISABLED` for the selected instance after confirmation.

#### Scenario: Toggle from enabled to disabled
- **WHEN** user selects an instance with the `InstanceScheduler` tag
- **THEN** the command SHALL rename the tag to `InstanceScheduler_DISABLED` while preserving the original tag value

#### Scenario: Toggle from disabled to enabled
- **WHEN** user selects an instance with the `InstanceScheduler_DISABLED` tag
- **THEN** the command SHALL rename the tag to `InstanceScheduler` while preserving the original tag value

#### Scenario: Preserve tag value during toggle
- **WHEN** toggling the scheduler tag name
- **THEN** the command SHALL preserve the original tag value exactly as it was before the toggle

### Requirement: Handle Instances Without Scheduler Tags
The `bmc ec2scheduler` command SHALL provide helpful guidance when users select instances without scheduler tags.

#### Scenario: User selects instance without scheduler tag
- **WHEN** user selects an instance that has neither `InstanceScheduler` nor `InstanceScheduler_DISABLED` tags
- **THEN** the command SHALL display a message indicating the instance does not have a scheduler tag configured
- **AND** the command SHALL prompt the user if they want to add the scheduler tag

#### Scenario: User confirms adding scheduler tag
- **WHEN** user selects an instance without scheduler tags
- **AND** user confirms they want to add the scheduler tag
- **THEN** the command SHALL display instructions that the user should add the tag manually in the AWS Console
- **AND** the command SHALL provide the exact tag name to add (InstanceScheduler)
- **AND** the command SHALL offer to open the AWS Console using `bmc console`

#### Scenario: User accepts opening AWS Console
- **WHEN** user is shown instructions to add scheduler tag manually
- **AND** user confirms they want to open the AWS Console
- **THEN** the command SHALL use `assumego` to open the AWS Console directly to the selected instance details page
- **AND** the command SHALL use the current AWS_PROFILE value
- **AND** the command SHALL construct the console URL with the instance ID in the format: `https://<region>.console.aws.amazon.com/ec2/home?region=<region>#InstanceDetails:instanceId=<instance-id>`
- **AND** the command SHALL execute with environment variables `GRANTED_ALIAS_CONFIGURED="true"` and `GRANTED_ENABLE_AUTO_REASSUME=true`
- **AND** the AWS Console SHALL open in the browser showing the instance details page where tags can be added

#### Scenario: User declines opening AWS Console
- **WHEN** user is shown instructions to add scheduler tag manually
- **AND** user declines to open the AWS Console
- **THEN** the command SHALL exit gracefully
- **AND** the command SHALL display the instructions for future reference

#### Scenario: User declines adding scheduler tag
- **WHEN** user selects an instance without scheduler tags
- **AND** user declines to add the scheduler tag
- **THEN** the command SHALL exit without making any changes or showing console instructions

### Requirement: Restrict Tag Name Modifications
The `bmc ec2scheduler` command SHALL only work with the specific tag names `InstanceScheduler` and `InstanceScheduler_DISABLED`.

#### Scenario: Only handle defined tag names
- **WHEN** the command operates on instance tags
- **THEN** it SHALL only recognize and modify tags named exactly `InstanceScheduler` or `InstanceScheduler_DISABLED`

### Requirement: Provide Clear User Feedback
The `bmc ec2scheduler` command SHALL provide clear feedback about toggle operations.

#### Scenario: Confirm successful toggle
- **WHEN** a scheduler tag is successfully toggled
- **THEN** the command SHALL display a confirmation message indicating the instance ID, the old tag name, the new tag name, and the preserved tag value

#### Scenario: Show error on toggle failure
- **WHEN** a tag toggle operation fails (e.g., due to AWS API errors or permission issues)
- **THEN** the command SHALL display a clear error message explaining what went wrong and exit with a non-zero status code

#### Scenario: Show current status before toggle
- **WHEN** displaying the instance selection table
- **THEN** each instance SHALL clearly show whether scheduling is currently "enabled" or "disabled"

