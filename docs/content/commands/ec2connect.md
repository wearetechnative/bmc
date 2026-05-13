---
title: "ec2connect"
weight: 32
description: "Connect to an EC2 instance via SSH or SSM"
---

`bmc ec2connect` connects to a running EC2 instance via SSH or AWS SSM Session Manager.

## Usage

```bash
bmc ec2connect               # Interactive instance picker
bmc ec2connect nginx         # Filter instances by name/ID/IP first
bmc ec2connect -i i-0abc123  # Connect to a specific instance ID
bmc ec2connect -u ubuntu     # SSH as a specific user (skips method picker)
```

## Connection methods

### SSH

Requires `ssh` on your PATH. BMC prompts for the SSH user (root, ubuntu, ec2-user, or custom) unless `-u` is specified.

### SSM Session Manager

Requires:
- `aws` CLI v2
- `session-manager-plugin`

Check with `bmc doctor`. No SSH keys or open ports needed.

## Stopped instances

If the selected instance is stopped, BMC behaviour depends on `ec2.auto_start_stopped` in config:

| Value | Behaviour |
|---|---|
| `prompt` (default) | Ask whether to start the instance |
| `always` | Start automatically without prompting |
| `never` | Exit with an error |

## Prerequisites

```bash
bmc doctor    # Check ssh, aws CLI, and session-manager-plugin
```
