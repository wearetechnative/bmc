## MODIFIED Requirements

### Requirement: Console URL opened via configured browser method
The `bmc console` command SHALL open the AWS console URL using the method configured in `~/.config/bmc/config.toml`. When `[console] firefox_containers = true`, the URL SHALL be passed to Firefox via the `ext+granted-containers:` scheme. Otherwise the existing platform default (`xdg-open` / `open`) SHALL be used.

#### Scenario: Default browser open (firefox_containers = false)
- **WHEN** user runs `bmc console` and `firefox_containers` is `false` or unset
- **THEN** the command SHALL open the signed console URL using `xdg-open` on Linux or `open` on macOS

#### Scenario: Firefox container open (firefox_containers = true)
- **WHEN** user runs `bmc console` and `firefox_containers = true`
- **THEN** the command SHALL invoke `firefox "ext+granted-containers:<signed-console-url>"`
- **AND** the Granted extension SHALL open the URL in a dedicated container tab
