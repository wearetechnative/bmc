---
title: "ec2ls"
weight: 31
description: "List EC2 instances as a table, plain text or JSON array"
---

`bmc ec2ls` lists all EC2 instances in the current AWS profile.

## Usage

```bash
bmc ec2ls            # Display as interactive table
bmc ec2ls --txt      # Output as aligned plain text without borders
bmc ec2ls --txt --no-headers   # Same, without the header line
bmc ec2ls --json     # Output as JSON array
```

## Table output

The table columns are configurable via `ec2.columns` in `~/.config/bmc/config.json`. See [Configuration](/setup/configuration/#ec2-columns).

## Plain text output

`--txt` prints the same columns as the table, aligned with spaces but without border lines. This makes the output easy to process with `grep`, `awk` and friends:

```bash
bmc ec2ls --txt | grep running | awk '{print $1}'
```

Add `--no-headers` to omit the header line, so every line is an instance (`--no-headers` requires `--txt`):

```bash
bmc ec2ls --txt --no-headers | grep -w running | awk '{print $1}'
```

Empty values (such as a missing public IP) are written as spaces, so whitespace-splitting tools like `awk` see fewer fields on those lines. Select on content (as above) or use `--json` with `jq` when you need a specific column reliably.

Messages such as an MFA refresh prompt go to stderr, so only the header and instance rows end up in the pipe. `--txt` and `--json` cannot be combined.

## JSON output

`--json` outputs all instances as a JSON array with AWS PascalCase keys, ignoring the `columns` config:

```json
[
  {
    "InstanceId": "i-0abc123",
    "Name": "prod-web-01",
    "PrivateIpAddress": "10.0.1.5",
    "PublicIpAddress": "",
    "State": "running",
    "Hibernate": "no",
    "Scheduler": "yes",
    "Profile": ""
  }
]
```

All fields are always present. Pipe into `jq` for filtering:

```bash
bmc ec2ls --json | jq '.[] | select(.State == "running") | .Name'
```
