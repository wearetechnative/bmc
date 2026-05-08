## MODIFIED Requirements

### Requirement: Console URL opened via configured browser method
The `bmc console` command SHALL open the AWS console URL using the method configured in `~/.config/bmc/config.toml`. Three modes are supported:

1. **Default** (`firefox_containers = false`, `chrome_profiles = false`): use platform default (`xdg-open` on Linux, `open` on macOS)
2. **Firefox containers** (`firefox_containers = true`): pass URL to Firefox via `ext+granted-containers:` scheme (requires Granted extension)
3. **Chrome profiles** (`chrome_profiles = true`, experimental): launch Chromium-based browser with a bmc-managed `--user-data-dir` per AWS profile

#### Scenario: Default browser open
- **WHEN** user runs `bmc console` and neither `firefox_containers` nor `chrome_profiles` is enabled
- **THEN** the command SHALL open the signed console URL using `xdg-open` on Linux or `open` on macOS

#### Scenario: Firefox container open
- **WHEN** user runs `bmc console` and `firefox_containers = true`
- **THEN** the command SHALL invoke `firefox "ext+granted-containers:<signed-console-url>"`
- **AND** the Granted extension SHALL open the URL in a dedicated container tab

#### Scenario: Chrome profile open (experimental)
- **WHEN** user runs `bmc console` and `chrome_profiles = true`
- **THEN** the command SHALL invoke the configured Chrome binary with `--user-data-dir=~/.config/bmc/chrome/profiles/<profile>/` and the signed console URL
