---
title: "console"
weight: 20
description: "Open the AWS Management Console for the current or selected profile"
---

`bmc console` opens the AWS Management Console in your browser for the current or selected AWS profile.

## Usage

```bash
bmc console              # Open console for current profile (AWS_PROFILE)
bmc console -p myprofile # Open console for a specific profile
bmc console -p           # Force interactive profile selection
bmc console -s ec2       # Open console at a specific service URL
```

## Browser options

### Firefox containers with Granted (recommended)

The [Granted](https://addons.mozilla.org/en-US/firefox/addon/granted/) Firefox extension enables container tabs — each AWS profile opens in its own isolated container with separate cookies and sessions.

**Setup:**

1. Install [Firefox](https://www.mozilla.org/firefox/)
2. Install the [Granted extension](https://addons.mozilla.org/en-US/firefox/addon/granted/) from the Firefox Add-ons store
3. Enable in `~/.config/bmc/config.json`:

```json
{
  "console": {
    "firefox_containers": true
  }
}
```

**How it works:**

When `firefox_containers = true`, BMC passes the AWS console URL to the Granted extension via a special URL scheme (`granted-containers://`). Granted opens it in a Firefox container tab named after the AWS profile. Each profile gets its own isolated session — you can be logged in to multiple AWS accounts simultaneously in separate tabs without any cross-contamination.

If Firefox is not your default browser, Granted still works as long as Firefox is installed.

### Chrome profile isolation

Opens the console in a dedicated Chrome profile per AWS account. See [Chrome Profiles](/advanced/chrome-profiles/) for setup details.

```json
{
  "console": {
    "chrome_profiles": true,
    "chrome_binary": "google-chrome"
  }
}
```

### Default browser

Without either option, BMC opens the AWS console URL in your system default browser.

## MFA

If MFA is enabled and the session has expired, `bmc console` automatically refreshes it before opening the browser. See [MFA setup](/setup/mfa/).
