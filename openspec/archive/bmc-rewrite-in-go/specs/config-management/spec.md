## ADDED Requirements

### Requirement: Read TOML config file
The system SHALL read `~/.config/bmc/config.toml` at startup of any command that requires configuration. If the file is absent, all settings SHALL use defaults.

#### Scenario: Config file present and valid
- **WHEN** `~/.config/bmc/config.toml` exists with valid TOML
- **THEN** system loads all settings from the file

#### Scenario: Config file absent
- **WHEN** `~/.config/bmc/config.toml` does not exist
- **THEN** system uses default values: `mfa.enabled=false`, `ec2.auto_start_stopped=prompt`

#### Scenario: Config file malformed
- **WHEN** `~/.config/bmc/config.toml` contains invalid TOML
- **THEN** system exits with a clear parse error showing the file path and line number

### Requirement: Config schema
The system SHALL support the following TOML structure:

```toml
[mfa]
enabled = true                        # bool, default: false
totp_script = "..."                   # string, default: ""
clipboard_command = "..."             # string, default: ""

[ec2]
auto_start_stopped = "prompt"         # string: always|never|prompt, default: prompt
```

#### Scenario: Unknown keys ignored
- **WHEN** config.toml contains keys not in the schema
- **THEN** system ignores them without error (forward compatibility)

### Requirement: Legacy config.env detection
The system SHALL detect the presence of `~/.config/bmc/config.env` and warn the user to migrate.

#### Scenario: config.env exists, config.toml absent
- **WHEN** `~/.config/bmc/config.env` exists and `~/.config/bmc/config.toml` does not
- **THEN** `bmc doctor` reports a warning: "Legacy config.env found. Please migrate to config.toml." with reference to documentation
