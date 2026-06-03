## Why

AWS console federation sessions expire after 1 hour — a hard AWS limit for role-assumed credentials. Users who switch browser tabs or work across multiple AWS accounts find their console sessions silently expired, requiring them to manually re-run `bmc console` each time.

## What Changes

- Add `bmc watcher` command with `start`, `stop`, and `status` subcommands
- Add `--watch` flag to `bmc console` that opens the console and ensures the watcher daemon is running
- Watcher daemon runs as a detached background process (no systemd or launchd dependency)
- Watcher tracks active sessions in `~/.config/bmc/watcher.json` and refreshes them ~5 minutes before expiry
- Watcher runs a local HTTP mini-server that serves a refresh page; the page fetches the new federation URL invisibly and closes itself
- Watcher self-terminates when no active sessions remain

## Capabilities

### New Capabilities

- `console-session-watcher`: Background daemon that monitors and refreshes AWS console federation sessions before they expire

### Modified Capabilities

- `aws-console-access`: `bmc console` gains a `--watch` flag that registers the session with the watcher daemon

## Impact

- New file: `cmd/watcher.go`
- New package: `internal/watcher/` (daemon loop, state file, HTTP server)
- Modified: `cmd/console.go` (add `--watch` flag, register session)
- Modified: `internal/awsops/console.go` (expose `BuildFederationURL` for reuse by watcher)
- New state file: `~/.config/bmc/watcher.json`
- No new external dependencies
