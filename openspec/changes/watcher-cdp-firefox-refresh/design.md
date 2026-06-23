## Context

The watcher daemon refreshes AWS console sessions by opening a temporary browser tab that runs a federation fetch and closes itself. Despite using `keepalive: true` and `window.close()`, Firefox shifts focus to an adjacent tab on close — not the tab the user was actively working in. This is because programmatic `window.close()` does not trigger Firefox's "return to previously active tab" history the way a user-initiated Ctrl+W does.

Firefox exposes a Remote Debugging Protocol (a subset of Chrome DevTools Protocol, CDP) over a local HTTP/WebSocket endpoint. This allows external processes to enumerate open tabs and execute JavaScript inside them, enabling a truly invisible refresh with zero browser disruption.

## Goals / Non-Goals

**Goals:**
- Refresh the AWS console session with no visible browser change (no new tab, no focus shift)
- Provide a one-command Firefox setup (`bmc watcher setup`)
- Degrade gracefully to the existing tab-based method when CDP is unavailable

**Non-Goals:**
- Supporting browsers other than Firefox (Chrome has its own debug protocol variant but different container model)
- Full CDP implementation — only the minimum needed: tab listing and `Runtime.evaluate`
- Replacing the tab-based fallback entirely

## Decisions

### Decision: Minimal CDP client in stdlib, no external dependency

The CDP over WebSocket requires:
1. `GET http://localhost:PORT/json` — list tabs (plain HTTP)
2. WebSocket upgrade to the tab's `webSocketDebuggerUrl`
3. Send/receive JSON-RPC messages (`Runtime.evaluate`)

Go's standard library handles HTTP natively. WebSocket frames can be implemented with `golang.org/x/net/websocket` (already a transitive dependency of the project) or with raw `net` + `bufio`. Using `golang.org/x/net/websocket` is acceptable since it is already in the module graph via other dependencies.

**Alternative considered**: Add `github.com/gorilla/websocket`. Rejected — adds a direct dependency for minimal usage.

### Decision: Find the AWS console tab by URL prefix, not container ID

CDP's `/json` endpoint returns each tab's `url` and `type`. AWS console tabs have URLs matching `https://*.console.aws.amazon.com/*`. Matching by URL prefix is simpler and more robust than mapping container name → `userContextId` → tab.

If multiple AWS console tabs are open (different accounts), pick the one whose URL matches the session's profile region, or — if ambiguous — pick any (all tabs in the same container share the same cookie jar for that container, so refreshing in any one of them refreshes the session for all).

**Alternative considered**: Map container name to `userContextId` via `containers.json`, then match tab's `userContextId` attribute in CDP response. Firefox's CDP tab listing does NOT expose `userContextId` directly, making this approach unreliable.

### Decision: `bmc watcher setup` writes to `user.js`, not `prefs.js`

`prefs.js` is managed by Firefox and overwritten on every startup. `user.js` is the correct override mechanism — Firefox reads it at startup and merges into `prefs.js`. Changes to `user.js` survive Firefox updates.

### Decision: `firefox_debug_port` defaults to 9222, disabled by 0

Port 9222 is the conventional CDP port (used by Chrome, Chromium, and Firefox). Defaulting to 9222 means users who already have CDP enabled (e.g., for other devtools usage) get CDP refresh automatically. Setting to `0` opts out completely and forces the tab fallback.

### Decision: CDP refresh timeout of 5 seconds

If the WebSocket connection or `Runtime.evaluate` does not complete within 5 seconds, the daemon falls back to tab-based refresh. This prevents a hung CDP connection from blocking the refresh indefinitely.

## Risks / Trade-offs

- **Firefox not running with debug port** → fallback to tab method; `bmc watcher status` reports which mode is active
- **Security**: CDP on localhost with `prompt-connection: false` allows any local process to execute JS in Firefox tabs. Acceptable for a developer tool on a personal machine. Not appropriate for shared machines.
- **Firefox update changes `user.js` path or preference names** → `bmc watcher setup` would need updating; low probability
- **Multiple Firefox instances** → CDP port conflict; only one instance can bind port 9222. Edge case for developer machines.

## Open Questions

- Should `bmc watcher setup` also check if Firefox is currently running and warn the user to restart? (Recommendation: yes, detect via `pidof firefox` or `pgrep firefox`)
