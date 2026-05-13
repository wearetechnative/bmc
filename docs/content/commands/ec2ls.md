---
title: "ec2ls"
weight: 31
description: "List EC2 instances as a table or JSON array"
---

`bmc ec2ls` lists all EC2 instances in the current AWS profile.

## Usage

```bash
bmc ec2ls            # Display as interactive table
bmc ec2ls --json     # Output as JSON array
```

## Table output

The table columns are configurable via `ec2.columns` in `~/.config/bmc/config.json`. See [Configuration](/setup/configuration/#ec2-columns).

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
