## ADDED Requirements

### Requirement: Configurable EC2 instance columns
The system SHALL allow users to configure which columns appear in EC2 instance tables and in what order via the `[ec2]` section of `~/.config/bmc/config.toml`. The `columns` field SHALL accept a list of PascalCase column name strings. Available columns are: `InstanceId`, `Name`, `PrivateIP`, `PublicIP`, `State`, `Hibernate`, `Scheduler`, `Profile`.

#### Scenario: User configures custom column list
- **WHEN** `~/.config/bmc/config.toml` contains `[ec2]` with `columns = ["InstanceId", "Name", "State"]`
- **THEN** all EC2 instance table commands (ec2ls, ec2connect, ec2stopstart, ec2scheduler) SHALL display only those three columns in that order

#### Scenario: Unknown column name in config
- **WHEN** the configured column list contains a name not matching any known field (e.g., `"Foo"`)
- **THEN** the column SHALL appear in the table with `"n/a"` as the value for every row, without error

#### Scenario: No config file present
- **WHEN** no `~/.config/bmc/config.toml` exists
- **THEN** the default column list SHALL be used: `["InstanceId", "Name", "PrivateIP", "PublicIP", "State", "Hibernate", "Scheduler"]`

### Requirement: Default column order puts Name second
The default column list SHALL be `["InstanceId", "Name", "PrivateIP", "PublicIP", "State", "Hibernate", "Scheduler"]`, placing the `Name` column immediately after `InstanceId`.

#### Scenario: Default display without config
- **WHEN** no `columns` setting is configured
- **THEN** `ec2ls` output SHALL show columns in the order: InstanceId, Name, PrivateIP, PublicIP, State, Hibernate, Scheduler

### Requirement: Shared row builder for all EC2 table commands
A single shared function SHALL build instance table rows for all EC2 commands. The `ec2ls` display path and the `selectInstanceID` interactive selection path SHALL use the same column list and the same row-building logic.

#### Scenario: Column list consistency across commands
- **WHEN** a user runs `ec2ls` and then `ec2connect` with the same config
- **THEN** both SHALL display the same columns in the same order

### Requirement: ec2find always includes Profile column
The `ec2find` command SHALL always display a `Profile` column regardless of the configured column list. If `Profile` is not present in the configured columns, it SHALL be appended.

#### Scenario: Profile column appended in ec2find
- **WHEN** the configured columns do not include `"Profile"`
- **THEN** `ec2find` output SHALL include all configured columns plus `"Profile"` as the last column

#### Scenario: Profile column not duplicated in ec2find
- **WHEN** the configured columns already include `"Profile"`
- **THEN** `ec2find` output SHALL not show `"Profile"` twice
