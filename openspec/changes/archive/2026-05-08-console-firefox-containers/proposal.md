## Why

`bmc console` currently opens the AWS console via `xdg-open`, which ignores Firefox Multi-Account Containers. Users who rely on the [Granted](https://addons.mozilla.org/en-US/firefox/addon/granted/) Firefox extension to isolate AWS accounts in separate container tabs get a plain browser tab instead, losing session isolation between accounts.

Relates to: https://github.com/wearetechnative/bmc/issues/39

## What Changes

- A new `[console] firefox_containers` option in `~/.config/bmc/config.toml` (default: `false`)
- When enabled, `bmc console` opens the AWS console URL via the `ext+granted-containers:` URL scheme, which the Granted Firefox extension intercepts to open the URL in a dedicated container tab
- When disabled, existing behavior is unchanged (`xdg-open` with raw signin URL)

## Capabilities

### New Capabilities

- `console-firefox-containers`: Controls whether `bmc console` opens URLs in Firefox container tabs via the Granted extension

### Modified Capabilities

- `aws-console-access`: When `firefox_containers = true`, the browser open step uses `ext+granted-containers:` scheme instead of `xdg-open`

## Impact

- `internal/config/config.go`: add `ConsoleConfig` struct with `FirefoxContainers bool` field
- `internal/awsops/console.go`: pass `FirefoxContainers` flag to `openBrowser()`; format URL as `ext+granted-containers:<url>` when enabled, and call `firefox` directly instead of `xdg-open`
- `cmd/console.go`: pass config to `OpenConsole()`
- `README.md`: document new `[console]` config section
