# config-format-json Specification

## Purpose
Defines requirements for bmc's configuration file format using JSON.

## Requirements

### Requirement: JSON Config File Location
The system SHALL read configuration from `~/.config/bmc/config.json`.

#### Scenario: Config file present
- **WHEN** `~/.config/bmc/config.json` exists
- **THEN** the system SHALL parse it as JSON
- **AND** apply the values to override defaults

#### Scenario: Config file absent
- **WHEN** `~/.config/bmc/config.json` does not exist
- **THEN** the system SHALL use default values without error

#### Scenario: Config file malformed
- **WHEN** `~/.config/bmc/config.json` exists but contains invalid JSON
- **THEN** the system SHALL return an error identifying the file path

### Requirement: Migration Hint
The system SHALL detect the presence of a legacy `config.toml` and guide the user.

#### Scenario: Legacy config detected
- **WHEN** `~/.config/bmc/config.json` does not exist
- **AND** `~/.config/bmc/config.toml` exists
- **THEN** the system SHALL print a hint to stderr explaining the config format has changed to JSON
- **AND** the hint SHALL include a brief JSON example of the equivalent config
- **AND** the system SHALL continue with default values (not error out)
