---
title: "ec2scheduler"
weight: 35
description: "Toggle the InstanceScheduler tag on an EC2 instance"
---

`bmc ec2scheduler` enables or disables the AWS Instance Scheduler tag on an EC2 instance.

## Usage

```bash
bmc ec2scheduler     # Pick instance → enable or disable scheduler
```

## What it does

The command adds or removes the `InstanceScheduler` tag on the selected instance. AWS Instance Scheduler uses this tag to automatically start and stop instances on a configured schedule — useful for saving costs on non-production environments.

| `Scheduler` column | Tag present |
|---|---|
| `yes` | `InstanceScheduler` tag is set |
| `no` | Tag is absent |

## See also

The [ec2](/commands/ec2/) unified command includes scheduler toggling in its action menu.
