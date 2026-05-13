---
title: "Configuration"
weight: 20
---

BMC is configured via `~/.config/bmc/config.json`.

## Example

```json
{
  "mfa": {
    "enabled": true,
    "totp_script": "/usr/bin/rbw get my-aws-mfa-entry --field totp",
    "copy_command": "wl-copy",
    "paste_command": "wl-paste | wtype -"
  },
  "ec2": {
    "auto_start_stopped": "prompt",
    "columns": ["InstanceId", "Name", "PrivateIP", "PublicIP", "State", "Hibernate", "Scheduler"]
  },
  "console": {
    "firefox_containers": true
  }
}
```

## Reference

| Key | Type | Default | Description |
|---|---|---|---|
| `mfa.enabled` | bool | `false` | Enable MFA session management |
| `mfa.totp_script` | string | `""` | Shell command to generate a 6-digit TOTP code |
| `mfa.copy_command` | string | `""` | Command to copy TOTP to clipboard (receives code via stdin) |
| `mfa.paste_command` | string | `""` | Command to simulate paste keystroke 300ms after copy |
| `ec2.auto_start_stopped` | string | `"prompt"` | `always` / `never` / `prompt` — what to do when connecting to a stopped instance |
| `ec2.columns` | []string | all columns | Columns to show in EC2 tables, in order |
| `console.firefox_containers` | bool | `false` | Open AWS console in Firefox container tabs via the [Granted](https://addons.mozilla.org/en-US/firefox/addon/granted/) extension |
| `console.chrome_profiles` | bool | `false` | Open AWS console in isolated Chrome profiles per AWS account |
| `console.chrome_binary` | string | `"google-chrome"` | Chrome binary to use (`"chromium"`, `"brave-browser"`, etc.) |

## EC2 columns

The `ec2.columns` field controls which columns appear in EC2 instance tables and in what order. Available values:

| Column | Description |
|---|---|
| `InstanceId` | EC2 instance ID |
| `Name` | Value of the `Name` tag |
| `PrivateIP` | Private IPv4 address |
| `PublicIP` | Public IPv4 address |
| `State` | Instance state (`running`, `stopped`, etc.) |
| `Hibernate` | Whether hibernation is enabled (`yes`/`no`) |
| `Scheduler` | Whether the `InstanceScheduler` tag is set (`yes`/`no`) |
| `Profile` | AWS profile name (always shown in `ec2find`) |

Unknown column names render as `n/a`.
