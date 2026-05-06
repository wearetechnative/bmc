## Context

bmc uses bubbletea for all interactive TUI prompts. `ui.Choose()` presents a filterable list and returns the selected item or `("", nil)` on cancel. Currently both ESC and Ctrl+C produce the same result, so multi-level menus (profsel, ecsconnect) have no way to go back — cancelling at any level exits the entire command.

## Goals / Non-Goals

**Goals:**
- ESC in a list goes back to the previous menu level
- Ctrl+C cancels the entire command from any level
- Callers opt in to back navigation via `errors.Is(err, ui.ErrBack)`
- Works for profsel (2 levels) and ecsconnect (4 levels)

**Non-Goals:**
- Animated transitions or breadcrumb display
- Back navigation for single-level selectors (ec2find, ec2stopstart, ec2connect method picker)
- Global undo/history beyond one level up

## Decisions

### Sentinel error over bool return or result type

`ui.Choose` returns `(string, error)`. Adding `ErrBack` fits this signature without change:
- `("", ErrBack)` — ESC pressed
- `("", nil)` — Ctrl+C pressed (cancel)
- `("item", nil)` — selection confirmed

**Alternatives considered:**
- `(string, bool, error)` — changes all call sites, more disruptive
- Result type enum — over-engineered for two states
- Loop on `""` without distinction — Ctrl+C at profile level would loop instead of exiting (undesirable)

`ErrBack` follows the `io.EOF` pattern: a sentinel that signals a non-error flow condition.

### ESC = back, Ctrl+C = cancel

In the current `listModel`, both keys set `quitting=true` and return `""`. We add `wentBack bool` to `listModel`. ESC sets `wentBack=true`; Ctrl+C does not. `Choose()` checks `wentBack` on the returned model and returns `ErrBack` accordingly.

### Callers implement the loop, not ui.Choose

`ui.Choose` is a primitive — it presents one list. Back navigation is a caller concern. Each multi-level caller wraps its levels in a `for` loop and uses `errors.Is(err, ui.ErrBack)` to decide whether to continue (go back) or return (propagate).

## Risks / Trade-offs

- **ErrBack propagated accidentally**: A caller that does `if err != nil { return err }` without checking `ErrBack` first will propagate it up as an error. Mitigation: document clearly; callers that don't need back navigation are unaffected (single-level selectors never receive ErrBack).
- **ecsconnect complexity**: 4 levels of nested loops. Mitigation: each level is an independent `for` block — readable and explicit.

## Migration Plan

No migration needed. All changes are internal to bmc. No config, data, or API changes.
