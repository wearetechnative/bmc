## Why

The watcher's current tab-based refresh opens a temporary Firefox container tab that shifts focus to an adjacent tab instead of returning to the user's previously active tab — a consistent disruption when working across multiple AWS accounts. Using Firefox's built-in Remote Debugging Protocol (CDP), the refresh can be executed invisibly inside an existing console tab with no new tab and no focus change.

## What Changes

- Add `bmc watcher setup` subcommand: detects the default Firefox profile, checks current devtools configuration, and writes the required preferences to `user.js`
- Add `watcher` section to `config.json` with a `firefox_debug_port` field (default: `9222`, set to `0` to disable)
- Watcher daemon refresh logic: attempt CDP-based refresh first; fall back to existing tab-based method if CDP is unavailable or fails
- `bmc watcher status` shows whether CDP is active or falling back to tab mode

## Capabilities

### New Capabilities

- `watcher-cdp-refresh`: Invisible session refresh via Firefox Remote Debugging Protocol — executes the federation fetch inside an existing container tab without opening a new tab

### Modified Capabilities

- `console-session-watcher`: Refresh strategy gains CDP as primary path with tab-based fallback; adds `bmc watcher setup` subcommand and `firefox_debug_port` config field

## Impact

- Modified: `internal/config/config.go` — add `WatcherConfig` struct with `FirefoxDebugPort int`
- Modified: `internal/watcher/daemon.go` — CDP refresh path in `refreshSession`
- New: `internal/watcher/cdp.go` — CDP client (WebSocket, tab discovery, `Runtime.evaluate`)
- Modified: `cmd/watcher.go` — add `setup` subcommand
- New: `internal/watcher/firefox.go` — Firefox profile detection, `user.js` writing, `containers.json` parsing
- No new external dependencies (uses Go standard library `net/http` and `golang.org/x/net/websocket`, or stdlib `net` for WebSocket frames)
