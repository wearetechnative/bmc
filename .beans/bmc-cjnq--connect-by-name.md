---
# bmc-cjnq
title: connect-by-name
status: completed
openspec-link: openspec/changes/archive/2026-05-13-ec2connect-connect-by-name/proposal.md
type: task
priority: normal
created_at: 2026-05-13T07:19:59Z
updated_at: 2026-05-13T09:30:00Z
---

ec2connect has -i flag to connect to specific instance_id. it should also be possible to give a name/part of name of the instance.
for example: ec2-nixhost-prod-992382728492  should be able to preselect with script option that has nixhost as a keyword

like: bmc e2connect --name nixhost
the function should search for ec2_instances with that name. If more than one is found it should give a list to select from
