## 1. Config: WatcherConfig struct

- [x] 1.1 Add `WatcherConfig` struct to `internal/config/config.go` with `FirefoxDebugPort int` field (JSON: `firefox_debug_port`)
- [x] 1.2 Add `Watcher WatcherConfig` field to the top-level `Config` struct
- [x] 1.3 Set default value `FirefoxDebugPort: 9222` in `Defaults()`

## 2. Firefox profile helpers

- [x] 2.1 Create `internal/watcher/firefox.go`: implement `FindDefaultProfile() (string, error)` that reads `~/.mozilla/firefox/profiles.ini` and returns the path of the default profile (`Default=1`)
- [x] 2.2 Implement `IsDebugPortConfigured(profileDir string) bool` that checks `user.js` and `prefs.js` in the profile dir for the required devtools preferences
- [x] 2.3 Implement `WriteDebugPortConfig(profileDir string, port int) error` that appends the four devtools preferences to `user.js`
- [x] 2.4 Implement `FirefoxIsRunning() bool` using `pgrep firefox` or scanning `/proc` for a firefox process

## 3. CDP client

- [x] 3.1 Create `internal/watcher/cdp.go`: implement `CDPClient` struct with `host string` and `port int`
- [x] 3.2 Implement `(c *CDPClient) IsReachable() bool` — HTTP GET to `http://host:port/json`, returns true on 200
- [x] 3.3 Implement `(c *CDPClient) FindConsoleTab() (wsURL string, err error)` — GET `/json`, parse JSON array, find first tab where `url` matches `https://*.console.aws.amazon.com/*`
- [x] 3.4 Implement `(c *CDPClient) Evaluate(wsURL, expression string) error` — open WebSocket to `wsURL`, send `{"id":1,"method":"Runtime.evaluate","params":{"expression":"..."}}`, wait for response with matching `id`, close WebSocket; enforce 5-second timeout
- [x] 3.5 Implement `(c *CDPClient) RefreshSession(federationURL string) error` — calls `FindConsoleTab`, builds the fetch JS expression, calls `Evaluate`

## 4. Daemon: CDP refresh path

- [x] 4.1 In `internal/watcher/daemon.go`: initialise a `CDPClient` in `RunDaemon()` using `bmcCfg.Watcher.FirefoxDebugPort`; skip if port is `0`
- [x] 4.2 Pass the `CDPClient` (or nil) into `refreshSession`
- [x] 4.3 In `refreshSession`: if CDPClient is non-nil, call `cdp.RefreshSession(signinURL)` first; on success, skip the browser open; on failure, log and fall through to tab-based method
- [x] 4.4 Track CDP availability in `WatcherState`: add `CDPActive bool` field to the state struct
- [x] 4.5 Update `WatcherState.CDPActive` in `RunDaemon()` based on `cdp.IsReachable()` result at startup

## 5. bmc watcher setup subcommand

- [x] 5.1 Add `watcherSetupCmd` to `cmd/watcher.go` and register it under `watcherCmd`
- [x] 5.2 Implement `runWatcherSetup`: call `watcher.FindDefaultProfile()`, handle not-found error with manual instructions
- [x] 5.3 Check `watcher.IsDebugPortConfigured(profileDir)`: if true, print "already configured" and exit
- [x] 5.4 Call `watcher.WriteDebugPortConfig(profileDir, port)` to write preferences
- [x] 5.5 Call `watcher.FirefoxIsRunning()`: if true, print warning to restart Firefox; if false, print "start Firefox to activate"

## 6. bmc watcher status: show CDP mode

- [x] 6.1 In `runWatcherStatus`: read `CDPActive` from `WatcherState` and print "CDP active" or "tab fallback" in the status output

## 7. Tests

- [x] 7.1 Unit test `FindDefaultProfile` with a temporary `profiles.ini` fixture
- [x] 7.2 Unit test `IsDebugPortConfigured` for the cases: pref in user.js, pref in prefs.js, pref absent
- [x] 7.3 Unit test `WriteDebugPortConfig` verifies the four preferences are written correctly
- [ ] 7.4 Manual test: run `bmc watcher setup`, restart Firefox, verify `localhost:9222/json` returns tab list
- [ ] 7.5 Manual test: run `bmc console --watch`, wait for refresh, verify no new tab opens and `bmc watcher status` shows "CDP active"
