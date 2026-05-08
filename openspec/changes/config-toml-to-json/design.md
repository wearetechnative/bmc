## Context

bmc currently uses `github.com/BurntSushi/toml` to parse `~/.config/bmc/config.toml`. The TOML library is only used in `internal/config/config.go` — nowhere else in the codebase. Switching to JSON removes this external dependency in favour of Go's stdlib `encoding/json`.

## Goals / Non-Goals

**Goals:**
- Replace TOML config with JSON (`~/.config/bmc/config.json`)
- Remove `BurntSushi/toml` dependency from `go.mod`
- Provide a clear migration hint when the old `config.toml` is found

**Non-Goals:**
- Supporting both formats simultaneously
- Auto-converting `config.toml` to `config.json` (hint only, user converts manually)
- Changing any config fields or defaults

## Decisions

**Use `encoding/json` (stdlib)**
No external library needed. The config struct is simple (flat nested structs). JSON tags replace TOML tags — a mechanical rename.

**Migration hint, not auto-migration**
If `config.json` is absent but `config.toml` exists, bmc prints a one-time hint to stderr explaining what to do. Auto-migration risks silently writing a file the user didn't ask for, and the manual conversion is trivial.

**Config field names stay the same**
JSON keys match the existing TOML keys (snake_case). Users only need to reformat their file, not rename keys.

## Risks / Trade-offs

- **Breaking change for all existing users** → Mitigated by migration hint and clear CHANGELOG entry
- **JSON is less readable than TOML for configs** (no comments, stricter syntax) → Accepted per product decision
- If `BurntSushi/toml` is added back later for other reasons, it just gets re-added to `go.mod`

## Migration Plan

Users with an existing `config.toml`:

```toml
# config.toml
[mfa]
enabled = true
totp_script = "/usr/bin/rbw get my-aws-mfa-entry --field totp"

[console]
firefox_containers = true
```

Convert to:

```json
{
  "mfa": {
    "enabled": true,
    "totp_script": "/usr/bin/rbw get my-aws-mfa-entry --field totp"
  },
  "console": {
    "firefox_containers": true
  }
}
```

bmc will print the hint on startup if `config.toml` exists and `config.json` does not.
