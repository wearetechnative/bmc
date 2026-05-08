---
# bmc-p8w2
title: fish-support
status: todo
type: task
priority: normal
created_at: 2026-05-08T20:00:00Z
updated_at: 2026-05-08T20:00:00Z
github-issue: https://github.com/wearetechnative/bmc/issues/41
---

Add Fish shell support for bmc shell integration.

Currently `bmc install-shell-integration` only supports zsh and bash. Fish shell users need a manual snippet or a dedicated Fish function.

The profsel wrapper needs to work with Fish syntax (functions, not aliases; `eval` works differently in Fish).
