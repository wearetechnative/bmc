## ADDED Requirements

### Requirement: bmc manages isolated Chrome profiles per AWS profile (experimental)
When `[console] chrome_profiles = true` is set in `~/.config/bmc/config.toml`, `bmc console` SHALL launch a Chromium-based browser with a dedicated `--user-data-dir` at `~/.config/bmc/chrome/profiles/<sanitized-aws-profile-name>/`. This feature is experimental.

#### Scenario: Chrome profile opened for known profile
- **WHEN** `chrome_profiles = true` and user runs `bmc console` with profile `TN-Production`
- **THEN** bmc SHALL launch the configured Chrome binary with `--user-data-dir=~/.config/bmc/chrome/profiles/TN-Production/`
- **AND** the signed AWS console URL SHALL be passed as the startup URL

#### Scenario: Separate profiles are isolated
- **WHEN** user opens `bmc console` with `TN-Production` and then with `TN-NonProduction`
- **THEN** each SHALL open in its own Chrome window with a separate user-data-dir
- **AND** the two sessions SHALL NOT share cookies or login state

### Requirement: New Chrome profile is seeded from default Chrome profile on first use
When a bmc Chrome profile directory does not yet exist, bmc SHALL attempt to seed it by copying extensions and preferences from the user's default Chrome profile. Cookies, login data, and history SHALL NOT be copied.

#### Scenario: First use with existing default Chrome profile
- **WHEN** `chrome_profiles = true` and the bmc profile directory does not exist
- **AND** a default Chrome profile is found at the expected OS path
- **THEN** bmc SHALL copy `Extensions/`, `Local Extension Settings/`, and `Preferences` to the new profile directory before launching Chrome

#### Scenario: First use without default Chrome profile
- **WHEN** `chrome_profiles = true` and the bmc profile directory does not exist
- **AND** no default Chrome profile is found
- **THEN** bmc SHALL create a fresh empty profile directory and launch Chrome without seeding

#### Scenario: Subsequent use reuses existing profile
- **WHEN** `chrome_profiles = true` and the bmc profile directory already exists
- **THEN** bmc SHALL launch Chrome directly without copying any files
- **AND** the existing session data (cookies, etc.) SHALL be preserved

### Requirement: Chrome binary is configurable
The `[console] chrome_binary` config option SHALL specify which Chromium-based binary to use (default: `google-chrome`).

#### Scenario: Default binary used when not configured
- **WHEN** `chrome_binary` is not set in config
- **THEN** bmc SHALL use `google-chrome` as the binary name

#### Scenario: Custom binary configured
- **WHEN** `chrome_binary = "brave-browser"` is set
- **THEN** bmc SHALL invoke `brave-browser` instead of `google-chrome`

#### Scenario: Binary not found in PATH
- **WHEN** the configured Chrome binary is not found in PATH
- **THEN** bmc SHALL return a clear error message naming the missing binary

### Requirement: AWS profile name is sanitized for use as directory name
Characters that are invalid or problematic in directory names SHALL be replaced with `-` when constructing the profile directory path.

#### Scenario: Profile name with slash
- **WHEN** the AWS profile name is `org/TN-Production`
- **THEN** the profile directory SHALL be `~/.config/bmc/chrome/profiles/org-TN-Production/`
