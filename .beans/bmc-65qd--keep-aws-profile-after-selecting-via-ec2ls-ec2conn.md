---
# bmc-65qd
title: Keep AWS profile after selecting via ec2ls, ec2connect etc
status: draft
type: task
priority: normal
created_at: 2026-09-08T09:49:52Z
updated_at: 2026-09-08T10:06:46Z
---

When selecting a intance or looking for an instance withouth profile, you select one. But after the command it does not save the selected profile to terminal variables. I propose we add a extra config option called `keep_aws_profile` with a true or false, when true it keeps the selected profile and false the opposite
