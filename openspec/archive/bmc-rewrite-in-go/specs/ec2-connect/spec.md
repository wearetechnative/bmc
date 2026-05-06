## ADDED Requirements

### Requirement: Select EC2 instance via TUI
The system SHALL present a filterable table of EC2 instances for selection before connecting.

#### Scenario: Instance selected interactively
- **WHEN** user runs `bmc ec2connect` with no flags
- **THEN** system shows a filterable table and user selects an instance

#### Scenario: Instance specified with -i flag
- **WHEN** user runs `bmc ec2connect -i <instance-id>`
- **THEN** system skips the selection table and uses the specified instance

#### Scenario: Filter committed before selection in connection method list
- **WHEN** user types to filter the connection method list and presses Enter while in filter mode
- **THEN** system commits the filter without immediately selecting; a second Enter selects the highlighted item

### Requirement: Handle stopped instances before connecting
The system SHALL handle stopped instances according to the `ec2.auto_start_stopped` config value.

#### Scenario: Config is `always`
- **WHEN** selected instance is stopped and config is `always`
- **THEN** system starts the instance automatically without prompting

#### Scenario: Config is `never`
- **WHEN** selected instance is stopped and config is `never`
- **THEN** system exits with an error message

#### Scenario: Config is `prompt` (default)
- **WHEN** selected instance is stopped and config is `prompt`
- **THEN** system asks user to confirm before starting the instance

### Requirement: SSH connection via existing .ssh/config ProxyCommand
The system SHALL connect via SSH by executing `ssh <user>@<instance-id>`. The ProxyCommand in `~/.ssh/config` handles `ec2-instance-connect` key upload and SSM tunnel — bmc does not re-implement this.

#### Scenario: SSH connection with user selection
- **WHEN** user selects SSH method and no `-u` flag is provided
- **THEN** system presents a list of common users (root, ubuntu, ec2_user, other) and prompts for custom input if `other` is selected

#### Scenario: SSH connection with -u flag
- **WHEN** user provides `-u <username>`
- **THEN** system skips user selection and connects as specified user

#### Scenario: SSH binary not found
- **WHEN** `ssh` binary is not found at connection time
- **THEN** system displays prerequisite error with install instructions (including 3 Nix options) and exits

#### Scenario: SSH exec handoff
- **WHEN** all checks pass
- **THEN** system uses `syscall.Exec` to replace the bmc process with `ssh <user>@<instance-id>`

### Requirement: SSM session connection
The system SHALL connect via SSM by executing `aws ssm start-session --target <instance-id>`.

#### Scenario: aws CLI or session-manager-plugin not found
- **WHEN** `aws` binary or `session-manager-plugin` is not found at SSM selection time
- **THEN** system displays prerequisite error with actionable install instructions (apt, brew, nix-env, nix profile, NixOS config) and exits

#### Scenario: SSM exec handoff
- **WHEN** all checks pass
- **THEN** system uses `syscall.Exec` to replace bmc process with `aws ssm start-session --target <instance-id>`
