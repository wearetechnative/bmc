---
title: "ec2stopstart"
weight: 34
description: "Stop or start an EC2 instance"
---

`bmc ec2stopstart` stops or starts an EC2 instance interactively.

## Usage

```bash
bmc ec2stopstart     # Pick instance → stop or start
```

BMC shows the current state of the selected instance and offers the appropriate action (Stop or Start). For hibernate-enabled instances, Stop also offers a hibernate option.

## Stopping with hibernate

If an instance has hibernation configured (`Hibernate: yes` in `ec2ls`), BMC will ask whether to hibernate instead of a cold stop. Hibernation saves RAM to disk and resumes faster.

## See also

The [ec2](/commands/ec2/) unified command includes stop/start in its action menu.
