---
# bmc-p6m8
title: instance-table-filter-search
status: completed
type: feature
priority: normal
created_at: 2026-09-08T05:22:16Z
updated_at: 2026-09-08T07:00:00Z
openspec-link: openspec/changes/archive/2026-09-08-instance-picker-filter-search/proposal.md
---

The interactive EC2 instance picker has no type-to-filter/search. Users can only scroll with the arrow keys to find an instance in the list.

## Context

The instance picker is `ui.SelectFromTable` (`internal/ui/table.go`), built on `bubbles/table`, which has no filtering. By contrast the profile picker uses `ui.Choose` (`internal/ui/list.go`) built on `bubbles/list` with `SetFilteringEnabled(true)`, so profile lists already filter as you type.

The shared selector `selectInstanceID` (`cmd/instancehelper.go` -> `ui.SelectFromTable`) is used by: `ec2connect`, `ec2`, `ec2stopstart`, `ec2scheduler`. Adding filtering to the shared widget fixes all four at once.

Related prior bean: bmc-q3f7 (filter-select-lists) was scrapped 2026-05-13; profile lists got filtering since, but the table-based instance selector never did.

## Desired behavior

- In the interactive instance table, the user can type to filter rows in real-time (match against the visible columns / instance name / id / IPs, consistent with the positional `[search]` fragment logic already in ec2connect).
- Filtering, clearing the filter, and selecting a filtered row all work; Esc/Ctrl+C cancel semantics preserved.
- Works across all commands using `selectInstanceID`.

## Notes / open design questions (for the proposal)
- Option A: add a `textinput` filter to the existing `bubbles/table` model in `table.go`.
- Option B: switch the instance selector from `bubbles/table` to a filterable `bubbles/list` (like the profile picker), rendering columns as list items.
- Keep the multi-column table readability either way.
