## Context

`openBrowser()` in `internal/awsops/console.go` currently calls `xdg-open <url>` on Linux and `open <url>` on macOS. This opens the URL in the default browser without any awareness of Firefox Multi-Account Containers.

The [Granted Firefox extension](https://addons.mozilla.org/en-US/firefox/addon/granted/) registers itself as the handler for the `ext+granted-containers:` protocol. When Firefox receives a URL with this scheme, the extension intercepts it and opens the destination in a named container tab, isolating cookies and sessions per AWS account.

The protocol was confirmed via the Firefox handlers.json on a system with Granted installed:
```json
"ext+granted-containers": {
  "action": 2,
  "handlers": [{ "name": "Open links in Granted Containers",
                  "uriTemplate": "moz-extension://.../opener.html#%s" }]
}
```

So to open a console URL in a Granted container, the call becomes:
```
firefox "ext+granted-containers:https://signin.aws.amazon.com/federation?..."
```

## Goals / Non-Goals

**Goals:**
- Add opt-in config option `[console] firefox_containers = true`
- When enabled, wrap the signin URL with `ext+granted-containers:` and invoke `firefox` directly
- Zero behavior change when option is disabled (default)

**Non-Goals:**
- Support for other container extensions (Sessionbox, honsiorovskyi's ext+container)
- macOS support for containers (Firefox + Granted works on macOS too, but `open -a Firefox` behavior is not tested)
- Auto-detecting whether the Granted extension is installed

## Decisions

### Decision: invoke `firefox` directly, not `xdg-open`

`xdg-open "ext+granted-containers:..."` only works if Firefox is the default browser. Calling `firefox` directly is more reliable and matches how `granted console --firefox` works internally. The `CustomBrowserPath` in Granted's own config (`~/.granted/config`) confirms this pattern.

**Alternative considered**: use `xdg-open` with the `ext+granted-containers:` scheme. Rejected — fails silently when Firefox is not the default browser.

### Decision: `firefox_containers` as a boolean in `[console]` config section

A simple boolean is sufficient for the current use case. Adding a new `[console]` section in `config.go` is clean and leaves room for future console-specific options (e.g., `console-profile-history`).

```toml
[console]
firefox_containers = true
```

**Alternative considered**: a `browser_command` template string. More flexible but overengineered for one use case.

### Decision: pass `cfg` to `OpenConsole()` instead of just the flag

`cmd/console.go` already loads `config.Load()` for MFA. Passing the full `config.Config` struct to `OpenConsole()` is consistent with how other commands work and avoids adding a new parameter type.

## Risks / Trade-offs

- **Firefox not in PATH**: if `firefox` is not found, the open fails. Mitigation: return a clear error message pointing to the `firefox_containers` config option.
- **Granted extension not installed**: the `ext+granted-containers:` URL opens in Firefox but the handler is missing — Firefox shows a "no handler" error. Acceptable — this is an opt-in feature; users enabling it are expected to have Granted installed.
- **macOS**: `firefox "ext+granted-containers:..."` should work on macOS if Firefox is in PATH. Not explicitly tested, but no special handling needed.

## Migration Plan

No migration needed. `firefox_containers` defaults to `false`; existing behavior is preserved.
