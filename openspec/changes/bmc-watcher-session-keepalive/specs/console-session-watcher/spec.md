## ADDED Requirements

### Requirement: Watcher daemon start
The `bmc watcher start` command SHALL start a background daemon process that monitors and refreshes active AWS console sessions. If a daemon is already running (PID in `~/.config/bmc/watcher.json` is alive), the command SHALL print a message indicating it is already running and exit without spawning a second daemon.

#### Scenario: No daemon running
- **WHEN** user runs `bmc watcher start`
- **AND** no watcher daemon is currently running
- **THEN** the command SHALL spawn a detached background daemon process
- **AND** print a confirmation message with the daemon PID

#### Scenario: Daemon already running
- **WHEN** user runs `bmc watcher start`
- **AND** a watcher daemon is already running
- **THEN** the command SHALL print "watcher already running (PID <n>)" and exit without spawning a second daemon

### Requirement: Watcher daemon stop
The `bmc watcher stop` command SHALL terminate the running daemon and clear the state file `~/.config/bmc/watcher.json`.

#### Scenario: Stop running daemon
- **WHEN** user runs `bmc watcher stop`
- **AND** a watcher daemon is running
- **THEN** the command SHALL send SIGTERM to the daemon process
- **AND** remove or clear `~/.config/bmc/watcher.json`

#### Scenario: Stop when no daemon running
- **WHEN** user runs `bmc watcher stop`
- **AND** no watcher daemon is running
- **THEN** the command SHALL print "watcher is not running" and exit cleanly

### Requirement: Watcher daemon status
The `bmc watcher status` command SHALL display all active sessions being monitored, including the profile name, browser mode (container/profile name), and time until the next refresh.

#### Scenario: Active sessions exist
- **WHEN** user runs `bmc watcher status`
- **AND** the daemon is running with one or more sessions
- **THEN** the command SHALL print each session with its profile name and a human-readable countdown to the next refresh (e.g., "refreshes in 47m")

#### Scenario: No active sessions or daemon not running
- **WHEN** user runs `bmc watcher status`
- **AND** no daemon is running or no sessions are registered
- **THEN** the command SHALL print "watcher is not running" or "no active sessions"

### Requirement: Session registration via console --watch
When `bmc console --watch` is used, the session SHALL be registered with the watcher daemon. If no daemon is running, one SHALL be started automatically.

#### Scenario: Register new session, daemon not running
- **WHEN** user runs `bmc console --watch`
- **AND** no watcher daemon is running
- **THEN** the console SHALL open normally
- **AND** a watcher daemon SHALL be started
- **AND** the session SHALL be registered in `~/.config/bmc/watcher.json`

#### Scenario: Register new session, daemon already running
- **WHEN** user runs `bmc console --watch`
- **AND** a watcher daemon is already running
- **THEN** the console SHALL open normally
- **AND** the session SHALL be appended to `~/.config/bmc/watcher.json`
- **AND** no second daemon SHALL be spawned

### Requirement: Automatic session refresh
The watcher daemon SHALL refresh each registered session approximately 5 minutes before its federation session expires, without user interaction.

#### Scenario: Refresh via invisible tab (fetch approach)
- **WHEN** a session's `refresh_at` time is reached
- **AND** the watcher's local HTTP server is accessible
- **THEN** the daemon SHALL open `http://localhost:<port>/refresh?t=<token>` in the browser container or profile
- **AND** the page SHALL fetch the federation URL with `credentials: include, mode: no-cors`
- **AND** the tab SHALL close itself after the fetch resolves

#### Scenario: Refresh fallback (visible tab)
- **WHEN** the fetch-based refresh fails (e.g., SameSite cookie restriction)
- **THEN** the daemon SHALL fall back to opening the federation URL directly in the container or profile
- **AND** the session SHALL be updated with the new expiry time

### Requirement: Daemon self-termination
The watcher daemon SHALL exit automatically when no active sessions remain.

#### Scenario: All sessions expired or removed
- **WHEN** all registered sessions have passed their expiry time
- **AND** no new sessions have been added
- **THEN** the daemon SHALL exit cleanly
- **AND** `~/.config/bmc/watcher.json` SHALL be cleared or removed

### Requirement: State file cleared on fresh daemon start
When the watcher daemon starts fresh (no living previous daemon), it SHALL clear `~/.config/bmc/watcher.json` before writing its own state.

#### Scenario: Stale state from crashed daemon
- **WHEN** `~/.config/bmc/watcher.json` contains a PID that is no longer alive
- **AND** `bmc watcher start` or `bmc console --watch` is run
- **THEN** the stale state file SHALL be cleared
- **AND** a new daemon SHALL be started with a clean state
