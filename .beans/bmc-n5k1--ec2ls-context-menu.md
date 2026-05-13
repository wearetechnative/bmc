---
# bmc-n5k1
title: ec2ls-context-menu
status: todo
type: task
priority: normal
created_at: 2026-05-08T20:00:00Z
updated_at: 2026-05-13T08:11:54Z
---

Add an action context menu to `bmc ec2ls` so users can act on a selected EC2 instance directly from the list.

Requested actions:
- Start instance
- Stop / reboot instance
- SSH into instance
- enable/disable schedule

Currently `ec2ls` is display-only. The interactive selection tables in `ec2connect`, `ec2stopstart`, etc. are separate commands. This would unify them under a single entry point.
