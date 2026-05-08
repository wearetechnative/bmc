## Why

`bmc console` supports Firefox container tabs (via Granted) for per-account session isolation, but users on Chrome, Brave, or Chromium have no equivalent. By managing isolated Chrome user-data directories per AWS profile, bmc can provide the same per-account isolation on Chromium-based browsers without requiring any browser extension.

Related bean: [bmc-c9y0](/.beans/bmc-c9y0--containerized-tabs-other-browsers.md)

## What Changes

- A new `[console] chrome_profiles = true` option in `~/.config/bmc/config.toml` (default: `false`)
- When enabled, `bmc console` launches Chrome (or a configured Chromium-based binary) with `--user-data-dir=~/.config/bmc/chrome/profiles/<aws-profile-name>/`
- On first use for a profile, bmc seeds the directory by copying extensions and preferences from the user's default Chrome profile — but not cookies, login data, or history
- Subsequent opens reuse the existing profile (session data persists between invocations)
- A new `[console] chrome_binary` option selects the browser binary (default: `google-chrome`; supports `chromium`, `brave-browser`, `microsoft-edge`, etc.)
- This feature is **experimental** — documented as such in config and README
- Existing `firefox_containers` and default browser behavior are unchanged

## Capabilities

### New Capabilities

- `console-chrome-profiles`: bmc-managed isolated Chrome user-data-dirs per AWS profile, seeded from the user's default Chrome profile on first use

### Modified Capabilities

- `aws-console-access`: The browser open step gains a third mode — `chrome_profiles` — alongside the existing default and `firefox_containers` modes

## Impact

- `internal/config/config.go`: add `ChromeProfiles bool` and `ChromeBinary string` fields to `ConsoleConfig`
- `internal/awsops/console.go`: add `openChromeProfile()` function; extend `openBrowser()` to handle the new mode
- `cmd/console.go`: no changes needed (config already passed through)
- `README.md`: document `[console]` chrome options, marked experimental
- New helper logic: detect default Chrome profile path per OS, selective copy on first use
