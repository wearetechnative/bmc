## ADDED Requirements

### Requirement: CDP-based invisible session refresh
When `firefox_debug_port` is configured and reachable, the watcher daemon SHALL refresh AWS console sessions by executing the federation fetch inside an existing browser tab via CDP, without opening a new tab or changing browser focus.

#### Scenario: CDP available, console tab found
- **WHEN** a session's `refresh_at` time is reached
- **AND** the CDP endpoint at `http://localhost:<firefox_debug_port>/json` is reachable
- **AND** at least one tab with a URL matching `https://*.console.aws.amazon.com/*` is open
- **THEN** the daemon SHALL execute `fetch(federationURL, {credentials:'include', mode:'no-cors', redirect:'follow'})` in that tab via `Runtime.evaluate`
- **AND** no new browser tab SHALL be opened
- **AND** browser focus SHALL NOT change

#### Scenario: CDP available, no console tab found
- **WHEN** a session's `refresh_at` time is reached
- **AND** the CDP endpoint is reachable
- **AND** no tab with an AWS console URL is open
- **THEN** the daemon SHALL fall back to the tab-based refresh method

#### Scenario: CDP unavailable
- **WHEN** a session's `refresh_at` time is reached
- **AND** the CDP endpoint is not reachable (connection refused, timeout)
- **THEN** the daemon SHALL fall back to the tab-based refresh method without error

#### Scenario: CDP refresh times out
- **WHEN** the CDP WebSocket connection or `Runtime.evaluate` call does not complete within 5 seconds
- **THEN** the daemon SHALL cancel the CDP attempt and fall back to the tab-based refresh method

### Requirement: Firefox profile setup command
The `bmc watcher setup` command SHALL configure the default Firefox profile to expose the CDP debug port on localhost, enabling invisible session refresh.

#### Scenario: Firefox not yet configured
- **WHEN** user runs `bmc watcher setup`
- **AND** the default Firefox profile does not have CDP debug preferences set
- **THEN** the command SHALL write the required preferences to `user.js` in the profile directory
- **AND** print instructions to restart Firefox

#### Scenario: Firefox already configured
- **WHEN** user runs `bmc watcher setup`
- **AND** the required preferences are already present in `user.js` or `prefs.js`
- **THEN** the command SHALL print a confirmation that setup is already complete

#### Scenario: Firefox profile not found
- **WHEN** user runs `bmc watcher setup`
- **AND** no Firefox profile directory can be found via `~/.mozilla/firefox/profiles.ini`
- **THEN** the command SHALL print an error with manual setup instructions

### Requirement: firefox_debug_port configuration
The `config.json` SHALL support a `watcher.firefox_debug_port` field to control CDP integration.

#### Scenario: Default port used
- **WHEN** `config.json` does not specify `watcher.firefox_debug_port`
- **THEN** the daemon SHALL attempt CDP on port `9222`

#### Scenario: CDP disabled via zero port
- **WHEN** `config.json` sets `watcher.firefox_debug_port` to `0`
- **THEN** the daemon SHALL skip CDP entirely and always use tab-based refresh
