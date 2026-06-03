## Context

`bmc console` generates a federation signin URL and opens it in the browser (Firefox container or Chrome profile). AWS caps federation sessions at 1 hour when using role-assumed credentials. There is currently no mechanism to refresh these sessions. Users must manually re-run `bmc console` when the session expires.

The watcher must work cross-platform (Linux, macOS) without system service managers (systemd, launchd), because bmc targets developer machines of varying configurations.

## Goals / Non-Goals

**Goals:**
- Automatically refresh AWS console sessions before they expire
- Support multiple concurrent sessions (different accounts/profiles open simultaneously)
- Work with existing Firefox containers (Granted) and Chrome profile browser modes
- No external daemon managers or persistent services
- Self-terminate cleanly when no sessions remain

**Non-Goals:**
- Refreshing sessions that require interactive MFA (TOTP script handles this if needed, but watcher does not prompt)
- Managing non-console bmc sessions (EC2, MFA tokens)
- Browser tab management beyond opening the refresh URL
- Windows support (bmc does not target Windows)

## Decisions

### Decision: Fork detection via environment variable

The same binary and subcommand (`bmc watcher start`) serves both the user-facing invocation and the background daemon. A `BMC_WATCHER_DAEMON=1` environment variable distinguishes the forked child from the parent.

- **Parent path**: forks itself with the env var set, `Setsid: true` (new process group, detached from terminal), then prints confirmation and exits.
- **Daemon path**: runs the watcher loop until no sessions remain.

**Alternative considered**: A hidden `--daemon` flag. Rejected because it creates a visible subcommand that appears in `--help` and must be documented, while the env-var approach is a transparent implementation detail.

### Decision: Singleton watcher process

Only one watcher process runs at a time. The state file stores the daemon PID. Before spawning, `bmc console --watch` and `bmc watcher start` check whether the stored PID is alive (`kill -0`). If alive, they append to the state file and return. If dead, they clear the state file, write the new session, and spawn a fresh daemon.

**Alternative considered**: Per-session daemon (one process per profile). Rejected because it complicates process management and wastes resources for users with many accounts.

### Decision: State file as sole IPC mechanism

`~/.config/bmc/watcher.json` is the communication channel between the CLI and the daemon. The daemon polls it every 30 seconds. New sessions written by `bmc console --watch` are picked up on the next poll cycle.

**Alternative considered**: Unix domain socket or named pipe. Rejected as over-engineering for the polling interval required (refresh every ~55 minutes).

### Decision: Local HTTP server for invisible refresh

Instead of opening the federation URL directly in a new browser tab (which leaves a visible tab pointing at the AWS console), the watcher opens `http://localhost:<port>/refresh?t=<token>` in the container. This local page:

1. Calls `fetch(federationURL, { credentials: 'include', mode: 'no-cors', redirect: 'follow' })` — this follows the federation redirect chain and sets the AWS session cookies within the container's cookie jar.
2. Calls `window.close()` after the fetch resolves, closing the tab.

The result: session is refreshed invisibly. The local HTTP server is started by the watcher daemon and its port is written to `watcher.json`.

**Risk**: AWS may use `SameSite=Lax` or `SameSite=Strict` cookies, in which case the `fetch` approach will not set the session cookie. Fallback behaviour: the tab navigates directly to the federation URL (visible tab, user closes manually). The watcher logs this outcome.

**Alternative considered**: Opening the federation URL directly in the container tab. Works, but leaves a visible tab requiring manual cleanup.

### Decision: State file cleared on fresh daemon start

When a new daemon spawns, it clears `watcher.json` and writes its own PID and the triggering session. This ensures stale entries from a crashed daemon do not persist.

**Consequence**: If the daemon crashes while two sessions are active, the second session is lost from the state file. The user must re-run `bmc console --watch` for that profile. This is acceptable given the low crash probability and the easy recovery path.

## Risks / Trade-offs

- **fetch SameSite failure** → Fallback to visible tab; log warning to stderr at daemon start
- **Port collision on localhost** → Watcher picks a random available port at startup; risk is negligible
- **Daemon zombie if bmc binary is updated** → Old daemon continues running with old binary until sessions expire; harmless
- **watcher.json race condition** → Two concurrent `bmc console --watch` invocations could both read a dead PID and both try to spawn. The second spawn will overwrite the state file, losing the first session. Mitigation: acceptable for the expected usage pattern (sequential console opens)

## Open Questions

- Should `bmc watcher status` show a human-readable countdown, or raw timestamps? (Recommendation: countdown, e.g. "refreshes in 47m")
- Should the watcher log to `~/.config/bmc/watcher.log` for debugging, or stay silent? (Recommendation: optional, enabled by a `--verbose` flag on `bmc watcher start`)
