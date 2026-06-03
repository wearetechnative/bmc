## 1. State File and Watcher Package

- [x] 1.1 Create `internal/watcher/` package with `state.go`: define `WatcherState` and `Session` structs matching the `watcher.json` schema (PID, started_at, port, sessions with profile, service, container_name, expiry, refresh_at)
- [x] 1.2 Implement `ReadState()` and `WriteState()` in `state.go` for atomic read/write of `~/.config/bmc/watcher.json`
- [x] 1.3 Implement `IsAlive(pid int) bool` helper using `kill -0` (syscall) to check whether the daemon PID is still running
- [x] 1.4 Implement `EnsureWatcher()` in `state.go`: reads state, checks PID liveness, clears stale state if dead, returns whether daemon is already running

## 2. Local HTTP Server

- [x] 2.1 Create `internal/watcher/server.go`: implement `StartServer() (port int, err error)` that binds to a random available localhost port
- [x] 2.2 Add `GET /refresh?t=<token>` handler that serves an HTML page using `fetch(decodedURL, {credentials:'include', mode:'no-cors', redirect:'follow'}).finally(() => window.close())`
- [x] 2.3 Add `GET /health` handler returning `200 OK` (used by status checks)
- [x] 2.4 Store the signed token → federation URL mapping in-memory in the server (map protected by mutex)

## 3. Daemon Loop

- [x] 3.1 Create `internal/watcher/daemon.go`: implement `RunDaemon()` that starts the HTTP server, writes PID and port to `watcher.json`, then enters the poll loop
- [x] 3.2 Implement the 30-second poll loop: read `watcher.json`, find sessions with `refresh_at <= now`, trigger refresh for each
- [x] 3.3 Implement `refreshSession(session Session, port int)`: call `awsops.OpenConsole` equivalent to build a new federation URL, register it in the server's token map, then open the browser via `ext+granted-containers` or Chrome profile pointing to `http://localhost:<port>/refresh?t=<token>`
- [x] 3.4 After refresh: update the session's `expiry` and `refresh_at` (expiry + 1h, refresh_at = expiry - 5min) in `watcher.json`
- [x] 3.5 Implement self-termination: if poll loop finds no sessions with a future `refresh_at`, exit cleanly and clear `watcher.json`
- [x] 3.6 Handle SIGTERM: on signal, clear `watcher.json` and exit

## 4. Expose Federation URL Builder in awsops

- [x] 4.1 Refactor `internal/awsops/console.go`: extract `BuildFederationURL(profile, service string, cfg config.Config) (signinURL string, expiry time.Time, err error)` as an exported function reusable by the watcher
- [x] 4.2 Update `OpenConsole` to call `BuildFederationURL` internally (no behaviour change)

## 5. bmc watcher Command

- [x] 5.1 Create `cmd/watcher.go` with `watcherCmd` (Use: "watcher") and register it in `rootCmd`
- [x] 5.2 Implement `bmc watcher start`: check `BMC_WATCHER_DAEMON` env var — if set, call `watcher.RunDaemon()`; if not set, check for running daemon, fork self with `BMC_WATCHER_DAEMON=1` and `Setsid: true`, print confirmation
- [x] 5.3 Implement `bmc watcher stop`: read PID from `watcher.json`, send SIGTERM, clear state file; handle "not running" case
- [x] 5.4 Implement `bmc watcher status`: read `watcher.json`, print each session (profile, container, countdown to next refresh); handle "not running" case

## 6. bmc console --watch Flag

- [x] 6.1 Add `--watch` / `-w` boolean flag to `consoleCmd` in `cmd/console.go`
- [x] 6.2 After successful `OpenConsole` call: if `--watch` is set, call `watcher.EnsureWatcher()` and append the session to `watcher.json`
- [x] 6.3 If `EnsureWatcher()` reports daemon not running: fork the daemon (same logic as `bmc watcher start`)

## 7. Integration and Testing

- [ ] 7.1 Manual test: `bmc console --watch` for a Firefox container profile — verify session refreshes after ~55 minutes without user interaction
- [ ] 7.2 Manual test: `bmc watcher status` shows active session with correct countdown
- [ ] 7.3 Manual test: `bmc watcher stop` terminates daemon and clears state file
- [ ] 7.4 Manual test: run `bmc console --watch` twice for two different profiles — verify both sessions appear in `bmc watcher status`
- [x] 7.5 Write unit tests for `IsAlive`, `ReadState`/`WriteState`, and the poll loop session expiry logic
