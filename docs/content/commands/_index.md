---
title: "Commands"
weight: 30
description: "All bmc commands with usage and examples"
---

| Command | Description |
|---|---|
| [profsel](/commands/profsel/) | Select an AWS profile and export it to your shell |
| [console](/commands/console/) | Open the AWS console for the current profile |
| [ec2](/commands/ec2/) | Unified EC2 command — pick instance, choose action |
| [ec2ls](/commands/ec2ls/) | List EC2 instances as table or JSON |
| [ec2connect](/commands/ec2connect/) | Connect via SSH or SSM Session Manager |
| [ec2find](/commands/ec2find/) | Find instances across all profiles in a group |
| [ec2stopstart](/commands/ec2stopstart/) | Stop or start an instance |
| [ec2scheduler](/commands/ec2scheduler/) | Toggle the InstanceScheduler tag |
| [ecsconnect](/commands/ecsconnect/) | Interactive shell into an ECS container |

## Other commands

```bash
bmc version                    # Show version
bmc doctor                     # Check dependencies and setup
bmc install-shell-integration  # Install profsel shell wrapper
bmc completion bash|zsh|fish   # Generate shell completions
```
