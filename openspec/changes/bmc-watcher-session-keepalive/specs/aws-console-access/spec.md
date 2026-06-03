## ADDED Requirements

### Requirement: Keep-alive flag for console sessions
The `bmc console` command SHALL support a `--watch` flag that registers the opened console session with the watcher daemon. If the watcher daemon is not running, it SHALL be started automatically.

#### Scenario: Open console with watch flag
- **WHEN** user runs `bmc console --watch`
- **THEN** the console SHALL open as normal (same profile selection and browser behaviour)
- **AND** the session SHALL be registered with the watcher daemon
- **AND** the watcher daemon SHALL be started if not already running

#### Scenario: Open console without watch flag
- **WHEN** user runs `bmc console` without `--watch`
- **THEN** the console SHALL open as normal
- **AND** no watcher session SHALL be registered
- **AND** the watcher daemon SHALL NOT be started or affected
