# Spec Delta

## Purpose

Provides borderless, column-aligned plain-text output for `bmc ec2ls` so instance lists can be processed with standard shell tools such as `grep`, `awk` and `cut`.

## ADDED Requirements

### Requirement: ec2ls plain-text output
The `bmc ec2ls` command SHALL support a `--txt` flag that writes instances to stdout as plain text: one header line followed by one line per instance, with columns aligned by spaces and no border or separator lines. The columns SHALL follow the `ec2.columns` configuration.

#### Scenario: Text flag produces aligned output without borders
- **WHEN** user runs `bmc ec2ls --txt`
- **THEN** the command SHALL write a header line and one line per instance to stdout
- **AND** the output SHALL NOT contain box-drawing border characters
- **AND** each column SHALL start at the same character position on every line
- **AND** columns SHALL be separated by at least two spaces

#### Scenario: Columns follow configuration
- **WHEN** `ec2.columns` is configured and user runs `bmc ec2ls --txt`
- **THEN** the columns and their order SHALL match the configured `ec2.columns`

#### Scenario: No trailing whitespace
- **WHEN** user runs `bmc ec2ls --txt`
- **THEN** no output line SHALL end with whitespace

#### Scenario: Text and JSON are mutually exclusive
- **WHEN** user runs `bmc ec2ls --txt --json`
- **THEN** the command SHALL exit with an error stating that the flags cannot be combined

### Requirement: ec2ls plain-text output without headers
The `bmc ec2ls` command SHALL support a `--no-headers` flag that omits the header line from `--txt` output. `--no-headers` SHALL only be accepted together with `--txt`.

#### Scenario: Header line omitted
- **WHEN** user runs `bmc ec2ls --txt --no-headers`
- **THEN** every line written to stdout SHALL describe one instance
- **AND** column alignment SHALL be based on the instance values only

#### Scenario: No-headers without txt is rejected
- **WHEN** user runs `bmc ec2ls --no-headers` without `--txt`
- **THEN** the command SHALL exit with the error "--no-headers requires --txt"
- **AND** the command SHALL NOT contact AWS
