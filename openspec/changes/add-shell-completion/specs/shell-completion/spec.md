# Shell Completion Specification

## ADDED Requirements

### Requirement: Bash Completion Support
The system SHALL provide bash completion for the `bmc` command and all registered subcommands.

#### Scenario: Complete main command
- **WHEN** user types `bmc <TAB>`
- **THEN** system SHALL display all available subcommands

#### Scenario: Complete subcommand
- **WHEN** user types `bmc ec<TAB>`
- **THEN** system SHALL complete or suggest matching subcommands (ec2ls, ec2connect, ec2stopstart, ec2find, ec2scheduler, ecsconnect)

#### Scenario: No matches
- **WHEN** user types `bmc xyz<TAB>` with no matching commands
- **THEN** system SHALL display no suggestions

### Requirement: Zsh Completion Support
The system SHALL provide zsh completion for the `bmc` command and all registered subcommands.

#### Scenario: Complete main command in zsh
- **WHEN** user types `bmc <TAB>` in zsh
- **THEN** system SHALL display all available subcommands with descriptions

#### Scenario: Complete subcommand in zsh
- **WHEN** user types `bmc ec<TAB>` in zsh
- **THEN** system SHALL complete or suggest matching subcommands with descriptions

#### Scenario: Navigate suggestions
- **WHEN** user presses TAB multiple times in zsh
- **THEN** system SHALL allow cycling through completion options

### Requirement: Completion Generation Command
The system SHALL provide a `gencompletions` command to generate shell completion scripts.

#### Scenario: Generate bash completion script
- **WHEN** user runs `bmc gencompletions bash`
- **THEN** system SHALL output bash completion script to stdout

#### Scenario: Generate zsh completion script
- **WHEN** user runs `bmc gencompletions zsh`
- **THEN** system SHALL output zsh completion script to stdout

#### Scenario: Invalid shell argument
- **WHEN** user runs `bmc gencompletions invalid`
- **THEN** system SHALL display error message and show supported shells

### Requirement: Installation Instructions
The system SHALL provide brief installation suggestions along with the generated completion script.

#### Scenario: Bash installation suggestion
- **WHEN** user runs `bmc gencompletions bash`
- **THEN** system SHALL output completion script followed by brief installation instructions for bash

#### Scenario: Zsh installation suggestion
- **WHEN** user runs `bmc gencompletions zsh`
- **THEN** system SHALL output completion script followed by brief installation instructions for zsh

#### Scenario: Instructions format
- **WHEN** installation suggestions are displayed
- **THEN** suggestions SHALL include common installation methods (e.g., sourcing in rc file, copying to completion directory)

### Requirement: Dynamic Command Discovery
The completion system SHALL dynamically discover available commands from the `bmc` command registration.

#### Scenario: New command appears in completion
- **WHEN** a new command is registered in `bmc`
- **THEN** completion system SHALL automatically include the new command without modification

#### Scenario: Command description in completion
- **WHEN** completion suggestions are displayed
- **THEN** system SHALL include command descriptions where supported by the shell

### Requirement: Documentation
The system SHALL provide clear documentation for enabling and using shell completion.

#### Scenario: README includes completion setup
- **WHEN** user reads README.md
- **THEN** document SHALL include instructions for enabling bash completion

#### Scenario: README includes zsh setup
- **WHEN** user reads README.md
- **THEN** document SHALL include instructions for enabling zsh completion

#### Scenario: Common installation patterns
- **WHEN** user follows documentation
- **THEN** instructions SHALL cover common installation locations (e.g., ~/.bashrc, ~/.zshrc, completion directories)
