## Requirements

### Requirement: ec2ls JSON output
The `bmc ec2ls` command SHALL support a `--json` flag that outputs all EC2 instances as a JSON array to stdout.

#### Scenario: JSON flag produces array output
- **WHEN** user runs `bmc ec2ls --json`
- **THEN** the command SHALL write a JSON array to stdout where each element represents one EC2 instance
- **AND** the output SHALL be valid JSON

#### Scenario: JSON output contains all fields
- **WHEN** user runs `bmc ec2ls --json`
- **THEN** each JSON object SHALL contain exactly the fields: `InstanceId`, `Name`, `PrivateIpAddress`, `PublicIpAddress`, `State`, `Hibernate`, `Scheduler`, `Profile`
- **AND** the `ec2.columns` configuration SHALL be ignored

#### Scenario: JSON key naming follows AWS CLI convention
- **WHEN** user runs `bmc ec2ls --json`
- **THEN** all JSON keys SHALL use PascalCase matching AWS CLI naming (e.g., `InstanceId` not `instance_id`)

#### Scenario: Empty list produces empty array
- **WHEN** user runs `bmc ec2ls --json` and there are no instances
- **THEN** the command SHALL output `[]`

#### Scenario: Table output unchanged without flag
- **WHEN** user runs `bmc ec2ls` without `--json`
- **THEN** the command SHALL behave exactly as before, rendering a formatted table

### Requirement: ec2find JSON output
The `bmc ec2find` command SHALL support a `--json` flag that outputs matching EC2 instances as a JSON array to stdout.

#### Scenario: JSON flag produces array output
- **WHEN** user runs `bmc ec2find <search> --json`
- **THEN** the command SHALL write a JSON array of matching instances to stdout
- **AND** interactive group selection SHALL still be presented via the terminal (bubbletea TUI on /dev/tty)

#### Scenario: JSON output contains all fields including Profile
- **WHEN** user runs `bmc ec2find <search> --json`
- **THEN** each JSON object SHALL contain exactly the fields: `InstanceId`, `Name`, `PrivateIpAddress`, `PublicIpAddress`, `State`, `Hibernate`, `Scheduler`, `Profile`
- **AND** the `Profile` field SHALL contain the name of the AWS profile the instance belongs to

#### Scenario: JSON flag is pipeable alongside TUI group selection
- **WHEN** user runs `bmc ec2find web --json | jq '.[].InstanceId'`
- **THEN** the bubbletea group picker SHALL render on the terminal
- **AND** only the JSON array SHALL be written to stdout

#### Scenario: No matches with JSON flag
- **WHEN** user runs `bmc ec2find <search> --json` and no instances match
- **THEN** the command SHALL output `[]`

#### Scenario: Table output unchanged without flag
- **WHEN** user runs `bmc ec2find <search>` without `--json`
- **THEN** the command SHALL behave exactly as before, rendering a formatted table
