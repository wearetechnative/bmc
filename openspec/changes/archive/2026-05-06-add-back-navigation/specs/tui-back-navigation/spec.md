## ADDED Requirements

### Requirement: ErrBack sentinel for back navigation
The `internal/ui` package SHALL export `var ErrBack = errors.New("ui: user navigated back")`. `ui.Choose()` SHALL return `("", ErrBack)` when the user presses ESC, and `("", nil)` when the user presses Ctrl+C. Callers opt in to back navigation by checking `errors.Is(err, ui.ErrBack)`.

#### Scenario: ESC returns ErrBack
- **WHEN** the user presses ESC in a `ui.Choose()` list
- **THEN** `ui.Choose()` SHALL return `("", ui.ErrBack)`

#### Scenario: Ctrl+C returns nil error
- **WHEN** the user presses Ctrl+C in a `ui.Choose()` list
- **THEN** `ui.Choose()` SHALL return `("", nil)`

#### Scenario: Caller loops on ErrBack
- **WHEN** a multi-level caller receives `ui.ErrBack` from an inner list
- **THEN** the caller SHALL re-present the previous menu level

#### Scenario: ErrBack not silently swallowed
- **WHEN** a caller uses `if err != nil { return err }` without checking ErrBack
- **THEN** ErrBack SHALL propagate as a non-nil error (it is an error value)
