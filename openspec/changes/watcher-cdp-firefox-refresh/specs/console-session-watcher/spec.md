## MODIFIED Requirements

### Requirement: Automatic session refresh
The watcher daemon SHALL refresh each registered session approximately 5 minutes before its federation session expires, without user interaction. The refresh SHALL use CDP if available, falling back to the tab-based method otherwise.

#### Scenario: Refresh via CDP (invisible)
- **WHEN** a session's `refresh_at` time is reached
- **AND** CDP is reachable on the configured `firefox_debug_port`
- **AND** an AWS console tab is open in the browser
- **THEN** the daemon SHALL execute the federation fetch inside the existing tab via CDP
- **AND** no new tab SHALL be opened
- **AND** browser focus SHALL NOT change

#### Scenario: Refresh via invisible tab (fetch approach)
- **WHEN** a session's `refresh_at` time is reached
- **AND** CDP is not available or no console tab is found
- **AND** the watcher's local HTTP server is accessible
- **THEN** the daemon SHALL open `http://localhost:<port>/refresh?t=<token>` in the browser container or profile
- **AND** the page SHALL fetch the federation URL with `credentials: include, mode: no-cors, keepalive: true`
- **AND** the tab SHALL close itself after the fetch is initiated

#### Scenario: Refresh fallback (visible tab)
- **WHEN** both CDP and the invisible tab method fail
- **THEN** the daemon SHALL fall back to opening the federation URL directly in the container or profile
- **AND** the session SHALL be updated with the new expiry time

## ADDED Requirements

### Requirement: Watcher status shows refresh mode
The `bmc watcher status` command SHALL indicate whether CDP is active or the tab-based fallback is in use.

#### Scenario: CDP active
- **WHEN** user runs `bmc watcher status`
- **AND** the daemon successfully reached the CDP endpoint at startup
- **THEN** the status output SHALL indicate "CDP active" alongside session information

#### Scenario: CDP not available
- **WHEN** user runs `bmc watcher status`
- **AND** CDP is not reachable
- **THEN** the status output SHALL indicate "tab fallback" alongside session information
