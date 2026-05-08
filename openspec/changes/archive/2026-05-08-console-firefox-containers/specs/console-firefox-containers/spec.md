## ADDED Requirements

### Requirement: Firefox container support is configurable
The `bmc` configuration SHALL support a `[console] firefox_containers` boolean option (default: `false`). When set to `true`, `bmc console` SHALL open the AWS console in a Firefox container tab via the Granted extension.

#### Scenario: firefox_containers enabled in config
- **WHEN** `~/.config/bmc/config.toml` contains `[console]\nfirefox_containers = true`
- **AND** the user opens the console interactively or via `-p`
- **THEN** the console URL SHALL be opened with `firefox "ext+granted-containers:<url>"`

#### Scenario: firefox_containers disabled (default)
- **WHEN** `firefox_containers` is `false` or not set
- **THEN** the console URL SHALL be opened via `xdg-open` (Linux) or `open` (macOS) as before

#### Scenario: firefox not found in PATH with containers enabled
- **WHEN** `firefox_containers = true`
- **AND** `firefox` is not found in PATH
- **THEN** bmc SHALL return a clear error indicating Firefox is required for container support
