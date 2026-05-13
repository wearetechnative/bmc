---
title: "Chrome Profiles"
weight: 10
---

BMC can open the AWS console in a dedicated, isolated Chrome profile per AWS account. Each profile has its own cookies and session — no cross-account contamination.

## Enable

```json
{
  "console": {
    "chrome_profiles": true,
    "chrome_binary": "google-chrome"
  }
}
```

For Brave or Chromium:

```json
{
  "console": {
    "chrome_profiles": true,
    "chrome_binary": "brave-browser"
  }
}
```

## How it works

On first use for a profile, BMC:
1. Creates a new Chrome profile directory at `~/.config/bmc/chrome/profiles/<profile-name>/`
2. Copies extensions and preferences from your default Chrome profile (without cookies or login data)
3. Opens Chrome with that profile pointing at the AWS console URL

Subsequent opens reuse the same profile directory, so your session persists between uses.

## Reset a profile

Delete the profile directory to reset it:

```bash
rm -rf ~/.config/bmc/chrome/profiles/<profile-name>/
```

BMC will recreate it fresh on next use.

## Notes

- Profile directories are not managed by BMC after creation — deleting them is safe
- This feature is marked **Experimental** in the configuration reference
- Firefox users: see [Firefox Containers](/commands/console/#firefox-containers-recommended) instead
