## 1. Config

- [x] 1.1 Add `ConsoleConfig` struct to `internal/config/config.go` with `FirefoxContainers bool` field (`toml:"firefox_containers"`)
- [x] 1.2 Add `Console ConsoleConfig` field to `Config` struct
- [x] 1.3 Update `Defaults()` to set `FirefoxContainers: false`

## 2. Browser open logic

- [x] 2.1 Update `OpenConsole()` signature in `internal/awsops/console.go` to accept `cfg config.Config`
- [x] 2.2 Update `openBrowser()` to accept a `firefoxContainers bool` parameter
- [x] 2.3 When `firefoxContainers = true`: look up `firefox` in PATH, return clear error if not found
- [x] 2.4 When `firefoxContainers = true`: call `firefox "ext+granted-containers:<url>"` instead of `xdg-open`

## 3. Wire up in console command

- [x] 3.1 In `cmd/console.go`, pass `cfg` to `awsops.OpenConsole()` (cfg is already loaded for MFA)

## 4. Documentation

- [x] 4.1 Add `[console]` section to the configuration reference in `README.md` with `firefox_containers` option

## 5. Verification

- [ ] 5.1 With `firefox_containers = false`: console opens via `xdg-open` as before
- [ ] 5.2 With `firefox_containers = true`: `firefox "ext+granted-containers:..."` is called and opens in a container tab
- [ ] 5.3 With `firefox_containers = true` and firefox not in PATH: clear error message shown

