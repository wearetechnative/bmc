---
title: "ec2find"
weight: 33
description: "Find EC2 instances across all profiles in an account group"
---

`bmc ec2find` searches for EC2 instances across all AWS profiles in an account group simultaneously.

## Usage

```bash
bmc ec2find nginx            # Search across profiles (select group interactively)
bmc ec2find nginx --json     # Same search, output as JSON array
```

## How it works

1. You select an AWS account group from an interactive picker
2. BMC queries all profiles in that group in parallel
3. Results matching your search string are shown

The search is a case-insensitive substring match on: instance ID, name, private IP, public IP, and profile name.

## JSON output

`--json` outputs matching instances as a JSON array. Group selection remains interactive via the terminal, so piping still works:

```bash
bmc ec2find nginx --json | jq '.[].InstanceId'
```

The JSON always includes the `Profile` field identifying which AWS profile each instance belongs to:

```json
[
  {
    "InstanceId": "i-0abc123",
    "Name": "nginx-prod-01",
    "PrivateIpAddress": "10.0.1.5",
    "PublicIpAddress": "3.75.10.20",
    "State": "running",
    "Hibernate": "no",
    "Scheduler": "yes",
    "Profile": "technative-prod"
  }
]
```
