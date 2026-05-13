## MODIFIED Requirements

### Requirement: install-shell-integration handles unwritable rc file
The system SHALL detect when the target rc file cannot be opened for writing due to a permission error. Instead of surfacing a generic error, the system SHALL print an explanation and manual wrapper snippets for all supported shell configurations.

#### Scenario: rc file is not writable
- **WHEN** `bmc install-shell-integration` is run and the rc file returns a permission denied error
- **THEN** system prints an explanation that the file is managed externally (e.g. home-manager) and outputs wrapper snippets for: home-manager zsh, home-manager bash, manual zsh/bash, and Fish shell

#### Scenario: rc file is writable
- **WHEN** `bmc install-shell-integration` is run and the rc file is writable
- **THEN** system appends the wrapper and reports success (existing behaviour unchanged)

#### Scenario: Fish shell on NixOS
- **WHEN** `bmc install-shell-integration` is run with Fish shell on NixOS
- **THEN** system prints NixOS-specific manual instructions with a `programs.fish.functions` home-manager snippet
- **AND** system does NOT attempt to write any files
