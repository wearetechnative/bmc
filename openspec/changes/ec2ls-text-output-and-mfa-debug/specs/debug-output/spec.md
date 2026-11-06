# Spec Delta

## Purpose

Defines a global `--debug` flag that reveals diagnostic output which is hidden during normal use, keeping the default output focused on results and required user actions.

## ADDED Requirements

### Requirement: Global debug flag
The `bmc` root command SHALL accept a persistent `--debug` flag that is available on every subcommand. When set, commands SHALL write diagnostic messages to stderr in addition to their normal output. Without the flag, diagnostic messages SHALL NOT be shown.

#### Scenario: Debug flag accepted by subcommands
- **WHEN** user runs any subcommand with `--debug` (e.g. `bmc ec2ls --debug` or `bmc --debug profsel`)
- **THEN** the command SHALL accept the flag without error

#### Scenario: Diagnostic output only with debug
- **WHEN** user runs a command that has diagnostic messages without `--debug`
- **THEN** those diagnostic messages SHALL NOT be written

#### Scenario: Debug output never pollutes stdout
- **WHEN** user runs a command with `--debug`
- **THEN** all diagnostic messages SHALL be written to stderr
- **AND** stdout SHALL contain the same content as without `--debug`
