## Context

`bmc console` currently has two browser open modes: default (`xdg-open`/`open`) and Firefox containers (via `ext+granted-containers:` URL scheme + Granted extension). The Firefox approach relies on a browser-native isolation primitive (Multi-Account Containers) intercepted by an extension. Chrome has no equivalent primitive, but supports fully isolated browser instances via `--user-data-dir`. bmc can manage these directories itself, giving per-AWS-profile Chrome isolation without requiring any browser extension.

## Goals / Non-Goals

**Goals:**
- Per-AWS-profile Chrome isolation via bmc-managed `--user-data-dir` directories
- Seed new profiles from the user's existing default Chrome profile (extensions + preferences only)
- Support any Chromium-based binary (Chrome, Chromium, Brave, Edge) via config
- Mark the feature as experimental in docs and config

**Non-Goals:**
- Supporting Chrome's built-in profile system (`--profile-directory`) — too fragile (internal names like "Profile 2")
- Syncing Chrome profile changes back to the default profile
- Managing extensions within bmc-created profiles (user can install extensions manually after first launch)
- Non-Chromium browsers beyond Firefox (already handled separately)

## Decisions

### user-data-dir over profile-directory
Chrome's `--profile-directory` flag selects a named sub-profile within an existing user-data-dir. The internal directory names (`Profile 1`, `Profile 2`) don't match display names and are fragile across Chrome reinstalls. Using `--user-data-dir` gives each AWS profile its own fully independent Chrome instance with a predictable path derived from the AWS profile name.

### Profile directory location: `~/.config/bmc/chrome/profiles/<aws-profile-name>/`
Keeps all bmc-managed state under `~/.config/bmc/`. Predictable, user-inspectable, easy to clean up per profile.

### Seed from default Chrome profile (selective copy)
Copying only `Extensions/`, `Preferences`, and `Local Extension Settings/` from the user's default Chrome profile gives the new profile the user's extensions without importing personal session data.

Files copied:
- `Extensions/` — installed extension binaries
- `Local Extension Settings/` — extension config (non-sensitive)
- `Preferences` — Chrome settings and toolbar layout

Files explicitly excluded:
- `Cookies` — personal sessions
- `Login Data` — saved passwords
- `History`, `Web Data`, `Visited Links` — browsing history

If the default Chrome profile cannot be found (path doesn't exist, Chrome not installed), bmc creates a fresh empty profile without seeding — no error.

### Chrome launch flags
```
--user-data-dir=<path>
--no-first-run          ← skip "Welcome to Chrome" wizard
--no-default-browser-check  ← suppress "make Chrome default" prompt
```
Extensions are allowed (not disabled) so users can install tools like 1Password into bmc profiles.

### Default Chrome binary: `google-chrome`
Falls back gracefully: if binary not found in PATH, return a clear error message. User can override with `chrome_binary = "chromium"` or `"brave-browser"`.

### Experimental flag in config
The feature is gated behind `chrome_profiles = true` (default: false). README documents it under an "Experimental" section. No behavioral change for users who don't opt in.

## Risks / Trade-offs

- [Profile directory size] Each bmc Chrome profile directory can grow to 100MB+ over time → Mitigation: document that `~/.config/bmc/chrome/profiles/` can be safely deleted to reset; no automatic cleanup.
- [Default Chrome profile path varies by OS/distro] Linux: `~/.config/google-chrome/Default/`, macOS: `~/Library/Application Support/Google/Chrome/Default/`, Brave: different path → Mitigation: detect by checking multiple known paths; skip seeding gracefully if not found.
- [Chrome running during seed copy] If Chrome is open, some files may be locked → Mitigation: copy is best-effort; skip locked files silently. Extensions directory is rarely locked.
- [Profile name → directory name mapping] AWS profile names may contain characters invalid in directory names (`/`, `:`) → Mitigation: sanitize profile name (replace invalid chars with `-`) before use as directory name.

## Open Questions

- Should bmc offer a `bmc console --reset-chrome-profile` flag to delete and re-seed a specific profile? (Not in scope for this change, candidate for follow-up.)
