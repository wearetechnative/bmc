## MODIFIED Requirements

### Requirement: profsel hints user when wrapper is absent
The system SHALL detect whether stdout is a terminal after printing the `export AWS_PROFILE=` line. If stdout is a terminal (output is not being captured by eval), the system SHALL print a single-line hint to stderr directing the user to install the shell wrapper.

#### Scenario: profsel run directly in terminal
- **WHEN** user runs `bmc profsel` directly (stdout is a TTY)
- **THEN** system prints `export AWS_PROFILE=<name>` to stdout AND prints a hint to stderr: `Tip: run 'bmc install-shell-integration' to set AWS_PROFILE automatically`

#### Scenario: profsel run via eval
- **WHEN** user runs `eval "$(bmc profsel)"` (stdout is a pipe)
- **THEN** system prints only `export AWS_PROFILE=<name>` to stdout; no hint is printed

### Requirement: install-shell-integration handles unwritable rc file
The system SHALL detect when the target rc file cannot be opened for writing due to a permission error. Instead of surfacing a generic error, the system SHALL print an explanation and manual wrapper snippets for all supported shell configurations.

#### Scenario: rc file is not writable
- **WHEN** `bmc install-shell-integration` is run and the rc file returns a permission denied error
- **THEN** system prints an explanation that the file is managed externally (e.g. home-manager) and outputs wrapper snippets for: home-manager zsh, home-manager bash, manual zsh/bash, and Fish shell

#### Scenario: rc file is writable
- **WHEN** `bmc install-shell-integration` is run and the rc file is writable
- **THEN** system appends the wrapper and reports success (existing behaviour unchanged)
