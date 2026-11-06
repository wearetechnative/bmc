# Spec Delta

## MODIFIED Requirements

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
