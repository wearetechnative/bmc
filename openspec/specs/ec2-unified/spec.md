### Requirement: ec2 command exists
The system SHALL provide a top-level `bmc ec2` command that combines instance selection with an action menu.

#### Scenario: Command is registered
- **WHEN** user runs `bmc ec2`
- **THEN** the command is listed in `bmc --help` and executes without error

### Requirement: Optional search argument filters instances
The command SHALL accept an optional positional argument that filters the instance list using a case-insensitive substring match on InstanceID, Name, PrivateIP, and PublicIP.

#### Scenario: No search argument shows full list
- **WHEN** user runs `bmc ec2` with no argument
- **THEN** the interactive instance picker shows all instances

#### Scenario: Matching search narrows picker
- **WHEN** user runs `bmc ec2 nginx` and multiple instances match
- **THEN** the picker shows only matching instances

#### Scenario: Single match skips picker
- **WHEN** user runs `bmc ec2 i-0abc123` and exactly one instance matches
- **THEN** the picker is skipped and the action menu is shown immediately

#### Scenario: No match returns error
- **WHEN** user runs `bmc ec2 xyz` and no instances match
- **THEN** the command exits with a clear error message containing the search term

### Requirement: Action menu presented after instance selection
After an instance is selected, the system SHALL show an action menu with the following items:
- **Connect SSH** — connects via SSH (same flow as `ec2connect` SSH)
- **Connect SSM** — connects via SSM Session Manager (same flow as `ec2connect` SSM)
- **Start instance** or **Stop instance** — label adapts to current instance state; uses same logic as `ec2stopstart`
- **Toggle scheduler** — enables or disables the InstanceScheduler tag (same logic as `ec2scheduler`)

#### Scenario: Action menu shown after selection
- **WHEN** user selects an instance
- **THEN** the action menu is displayed with all applicable actions

#### Scenario: Start/Stop label reflects current state
- **WHEN** selected instance is running
- **THEN** the action menu shows "Stop instance"

#### Scenario: Start/Stop label reflects current state (stopped)
- **WHEN** selected instance is stopped
- **THEN** the action menu shows "Start instance"

#### Scenario: Connect SSH delegates to existing SSH flow
- **WHEN** user selects "Connect SSH"
- **THEN** the SSH user selection and connection behave identically to `bmc ec2connect` SSH

#### Scenario: Connect SSM delegates to existing SSM flow
- **WHEN** user selects "Connect SSM"
- **THEN** the SSM session behaves identically to `bmc ec2connect` SSM

### Requirement: Existing EC2 commands unchanged
The system SHALL not modify the behaviour or interface of `ec2ls`, `ec2connect`, `ec2stopstart`, `ec2scheduler`, or `ec2find`.

#### Scenario: ec2connect still works independently
- **WHEN** user runs `bmc ec2connect`
- **THEN** the command behaves exactly as before
