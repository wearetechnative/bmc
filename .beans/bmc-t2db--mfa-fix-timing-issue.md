---
# bmc-t2db
title: mfa-fix-timing-issue
status: draft
type: task
priority: normal
created_at: 2026-05-21T05:59:56Z
updated_at: 2026-05-21T06:03:50Z
---

aws-switch checkt mfa geldigheid en vernieuwd indien niet meer geldig is.
Dit lijkt niet helemaal te werken
Om 07:57 deed ik:
 aws-switch -p technative       

en kreeg ik de melding
-- Using AWS source-profile: technative
Current MFA Session Valid, until: 2026-05-21 06:11:25

dat tijdstip klopt dus niet. Wellicht wordt utc gebruikt?
