---
title: "ec2"
weight: 30
description: "Unified EC2 command — select an instance and act on it"
---

`bmc ec2` is the unified entry point for EC2 operations. Pick an instance once and choose an action from a menu — no need to repeat instance selection across separate commands.

## Usage

```bash
bmc ec2              # Interactive instance picker → action menu
bmc ec2 nginx        # Filter instances by name/ID/IP first
bmc ec2 i-0abc123    # Single match skips the picker
```

The optional search argument filters instances by a case-insensitive substring match on instance name, ID, private IP, or public IP.

## Action menu

After selecting an instance, a menu appears:

| Action | Description |
|---|---|
| **Connect SSH** | SSH into the instance (same flow as `ec2connect`) |
| **Connect SSM** | SSM Session Manager shell (same flow as `ec2connect`) |
| **Start instance** / **Stop instance** | Label adapts to current state |
| **Toggle scheduler** | Enable/disable the `InstanceScheduler` tag |

## Related commands

- [ec2ls](/commands/ec2ls/) — list instances as a table or JSON
- [ec2connect](/commands/ec2connect/) — connect directly without the action menu
- [ec2find](/commands/ec2find/) — search across all profiles in a group
- [ec2stopstart](/commands/ec2stopstart/) — stop or start without the full menu
- [ec2scheduler](/commands/ec2scheduler/) — toggle the scheduler tag directly
