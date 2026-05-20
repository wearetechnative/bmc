## ADDED Requirements

### Requirement: SSH identity file flag
The `ec2connect` command SHALL accept a `-k`/`--key` flag that specifies a path to an SSH identity file. When provided, the value SHALL be passed to the `ssh` binary as `-i <path>`. No file existence validation SHALL be performed by bmc.

#### Scenario: Key flag passes identity file to ssh
- **WHEN** the user runs `bmc ec2connect -k /path/to/key.pem`
- **THEN** the resulting ssh invocation includes `-i /path/to/key.pem`

#### Scenario: Key flag with invalid path
- **WHEN** the user runs `bmc ec2connect -k /nonexistent/key.pem`
- **THEN** bmc passes the path to `ssh` without validation and `ssh` reports the error

### Requirement: SSH auto-selected when key flag is provided
When the `-k`/`--key` flag is set, `ec2connect` SHALL automatically use SSH as the connection method without showing the connection method picker.

#### Scenario: Method picker skipped when key is provided
- **WHEN** the user runs `bmc ec2connect -k /path/to/key.pem` without `-u`
- **THEN** the connection method picker is not shown and SSH is used

#### Scenario: Key and user flags together skip all pickers
- **WHEN** the user runs `bmc ec2connect -k /path/to/key.pem -u ubuntu`
- **THEN** neither the method picker nor the user picker is shown
