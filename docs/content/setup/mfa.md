---
title: "MFA"
weight: 30
---

BMC handles MFA automatically. When you run `bmc profsel` or `bmc console`, BMC checks if your session credentials are still valid and refreshes them if needed — no separate command required.

## Requirements

1. Set `mfa.enabled = true` in `~/.config/bmc/config.json`
2. Add a `[profile-long-term]` section to `~/.aws/credentials`:

```ini
[technative-long-term]
aws_access_key_id     = AKIA...
aws_secret_access_key = ...
aws_mfa_device        = arn:aws:iam::123456789012:mfa/your-username
```

When the session expires, BMC prompts for a 6-digit TOTP code. If `totp_script` is configured, BMC runs it automatically to fetch the code.

## TOTP script

`totp_script` is executed via `sh -c`. It should print a 6-digit code to stdout.

Example with [rbw](https://github.com/doy/rbw) (Bitwarden CLI):
```json
{
  "mfa": {
    "totp_script": "rbw get my-aws-mfa-entry --field totp"
  }
}
```

Interactive TUI tools (e.g. `gum filter`) work because the script runs with the terminal attached.

## Per-profile TOTP scripts

If you manage multiple AWS accounts with different TOTP credentials, use `profile_scripts` to map each source profile to its own script. The global `totp_script` acts as a fallback for profiles not listed.

```json
{
  "mfa": {
    "enabled": true,
    "totp_script": "rbw code \"Technative AWS (new)\"",
    "profile_scripts": {
      "wvandrtoorren": "rbw code \"Personal AWS\""
    }
  }
}
```

When `bmc profsel` resolves the source profile (e.g. `wvandrtoorren`), it uses the matching script from `profile_scripts`. Profiles not listed fall back to `totp_script`. If neither is configured, BMC prompts for manual input.

## Clipboard integration

After fetching the TOTP code, BMC can copy it to the clipboard and simulate a paste keystroke.

`copy_command` receives the code via **stdin**. `paste_command` runs 300ms later and simulates a paste in the focused window.

### Wayland (wl-clipboard)

```json
{
  "mfa": {
    "copy_command": "wl-copy",
    "paste_command": "wl-paste | wtype -"
  }
}
```

### X11 (xclip)

```json
{
  "mfa": {
    "copy_command": "xclip -selection clipboard",
    "paste_command": "xdotool key ctrl+v"
  }
}
```

Both fields are optional. If only `copy_command` is set, the code is copied but not auto-pasted. If neither is set, no clipboard interaction occurs.
