## Context

See proposal.md — Why. The interactive instance selector is `ui.SelectFromTable` in `internal/ui/table.go`, built on `bubbles/table`, which has no filtering — the `tableModel.Update` only handles enter/esc/ctrl+c/q and delegates the rest to the table (scrolling). `selectInstanceID` (`cmd/instancehelper.go`) feeds it pre-rendered rows via `awsops.InstanceRows(instances, cols)`, so the widget currently only knows the visible column strings, not the underlying instance fields.

The profile picker already solves the same UX problem with `ui.Choose` (`internal/ui/list.go`), which uses `bubbles/list` with `SetFilteringEnabled(true)` and a per-item `FilterValue()`. The positional `[search]` pre-filter in `ec2connect.go` already matches `InstanceID`/`Name`/`PrivateIP`/`PublicIP` (case-insensitive substring) — the interactive filter should match the same fields for consistency.

## Goals / Non-Goals

**Goals:**
- Real-time type-to-filter in the interactive instance picker, matching the same fields as the positional `[search]` fragment.
- One implementation in the shared widget so all four commands benefit without call-site logic changes.
- Preserve current selection, scrolling, and cancel (Esc/Ctrl+C) semantics.

**Non-Goals:**
- No change to the positional `[search]` argument behavior (that pre-filter stays as-is).
- No change to which columns are displayed or to `cfg.EC2.Columns`.
- No fuzzy/ranked matching — plain case-insensitive substring, consistent with existing behavior.

## Decisions

### Decision: Match against instance fields, not just rendered columns
`selectInstanceID` must give the widget a per-row filter string built from `InstanceID + Name + PrivateIP + PublicIP` (the same fields `ec2connect` already concatenates), independent of which columns `cfg.EC2.Columns` happens to display. Otherwise a user filtering on private IP would fail whenever that column is hidden.

- **Alternative — filter on the visible column text only**: rejected; behavior would vary with column config and diverge from the positional `[search]`.

### Decision: How to add filtering to the table widget
Two viable approaches; the implementer picks during apply based on which keeps the multi-column layout cleanest:

- **Option A — add a `textinput` filter to the existing `bubbles/table` model.** Keep `bubbles/table` for rendering; add a `bubbles/textinput`, and on each keystroke recompute the visible `[]btable.Row` from the full row set using the per-row filter strings. Preserves the exact table look; more manual wiring (filter state, focus, footer hint).
- **Option B — render instances through the filterable `ui.Choose`/`bubbles/list` path.** Reuse the already-working `SetFilteringEnabled(true)` machinery by giving each list item a `FilterValue()` of the concatenated fields and a `Title` that renders the columns. Less new code, unifies both pickers on one widget; risk is column alignment inside list items vs. a real table.

Recommendation: prefer **Option A** if preserving the aligned table columns matters, **Option B** if unifying on one filterable widget is preferred. Decide in apply; the spec is agnostic to the choice.

### Decision: Keep the shared entry point
Filtering lives entirely inside `ui.SelectFromTable` (+ the row/field data passed by `selectInstanceID`). The four commands call `selectInstanceID` unchanged, so the feature lands everywhere at once with no per-command edits.

## Risks / Trade-offs

- **Key binding collision** (typing a filter vs. `q`-to-quit / navigation keys) → Mitigation: while the filter input is focused, printable keys go to the filter; reserve Esc/Ctrl+C for cancel and Enter for select, mirroring how `ui.Choose` handles `list.Filtering` state in `list.go`.
- **Non-interactive fallback** (`selectPlainTable`) has no filtering → Mitigation: leave the plain fallback as-is; type-to-filter is a TTY-only enhancement, and the positional `[search]` still works everywhere.
- **Option B column alignment** → Mitigation: if chosen, pad column cells within each list item; otherwise choose Option A.

## Open Questions

- None that affect the spec or task breakdown. The A-vs-B widget choice is an implementation detail resolved during apply and does not change observable behavior.
