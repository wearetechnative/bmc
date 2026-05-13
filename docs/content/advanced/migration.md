---
title: "Migration from bash version"
weight: 30
---

If you were using the previous bash-based version of BMC, here is what changed and how to migrate.

## Steps

1. Install the new binary (same name: `bmc`)
2. Run `bmc install-shell-integration` — replaces the `source bmc profsel` pattern
3. Create `~/.config/bmc/config.json` — replaces `~/.config/bmc/config.env`
4. Run `bmc doctor` to verify setup

## Breaking changes

| Old | New |
|---|---|
| `source bmc profsel` | `eval "$(bmc profsel)"` — handled automatically by the shell wrapper |
| `~/.config/bmc/config.env` | `~/.config/bmc/config.json` |
| `bmc gencompletions` | `bmc completion bash\|zsh` (Cobra built-in) |

## Config migration

Old `config.env`:
```bash
BMC_MFA_ENABLED=true
BMC_MFA_TOTP_SCRIPT="rbw get my-mfa --field totp"
```

New `config.json`:
```json
{
  "mfa": {
    "enabled": true,
    "totp_script": "rbw get my-mfa --field totp"
  }
}
```

See the full [configuration reference](/setup/configuration/) for all available fields.
