## 1. Pass filter data into the picker

- [x] 1.1 Update `selectInstanceID` (`cmd/instancehelper.go`) to build a per-row filter string from `InstanceID + Name + PrivateIP + PublicIP` (case-insensitive) and pass it alongside the rendered rows to the widget; verify `go build ./...`
- [x] 1.2 Extend the `ui.SelectFromTable` signature (or add a filterable variant) to accept the per-row filter values without breaking existing non-instance callers; verify `go build ./...` and `grep -rn SelectFromTable cmd/ internal/` compiles for all callers

## 2. Implement type-to-filter in the widget

- [x] 2.1 Add filter state to the interactive table selector in `internal/ui/table.go` (per design decision A or B): capture printable keys into a filter query, recompute visible rows on each keystroke via case-insensitive substring match, and show the active query in the footer; verify manually that typing narrows the list
- [x] 2.2 Preserve key semantics: Enter selects the highlighted (filtered) row, Esc/Ctrl+C cancel returns no selection, and clearing the query restores the full list; verify each interaction manually
- [x] 2.3 Leave the non-interactive `selectPlainTable` fallback unchanged and confirm it still works when stdout is not a TTY (e.g. piped)

## 3. Verify across commands

- [x] 3.1 Verify type-to-filter works in the interactive picker for `ec2connect`, `ec2`, `ec2stopstart`, and `ec2scheduler` (each calls `selectInstanceID`) with no per-command code changes
- [x] 3.2 Verify the interactive filter matches the same fields as the positional `[search]` fragment: filtering by a private IP works even when the private-IP column is hidden via `cfg.EC2.Columns`

## 4. Docs & changelog

- [x] 4.1 Add a `## NEXT VERSION` entry to `CHANGELOG.md` under `### Added` describing type-to-filter in the interactive instance picker
- [x] 4.2 Update the docs site (`docs/content/`) to mention that the instance picker is filterable by typing
