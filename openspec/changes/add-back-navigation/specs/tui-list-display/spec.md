## MODIFIED Requirements

### Requirement: List height matches item count and terminal size
The system SHALL allocate list height based on item count and actual terminal height. For lists with 4 or fewer items: desired height is `len(items) + 3`, help bar and pagination hidden. For lists with more than 4 items: desired height is `min(len(items) + 6, 40)`, help bar and pagination visible. In all cases the height SHALL be clamped to `terminalHeight - 2`. No blank item-slots SHALL appear between items. The list SHALL use `/dev/tty` for I/O so it renders correctly in eval-context and sub-group navigation.

#### Scenario: Short list shows no blank rows
- **WHEN** `ui.Choose()` is called with 2 items (e.g. "ssh" / "ssm")
- **THEN** the rendered list contains no blank lines between items

#### Scenario: Short list hides help bar and pagination
- **WHEN** `ui.Choose()` is called with 4 or fewer items
- **THEN** the help bar and pagination dots are not shown

#### Scenario: Long list retains help bar and cap
- **WHEN** `ui.Choose()` is called with more than 4 items
- **THEN** the help bar and pagination are shown and height is `min(len(items) + 6, 40)` clamped to terminal height

#### Scenario: Terminal-aware resize
- **WHEN** the terminal is resized during list display
- **THEN** the list width and height update to fit the new terminal dimensions

#### Scenario: Eval-context rendering
- **WHEN** `eval "$(bmc profsel)"` is used and stdout is captured
- **THEN** the TUI renders correctly on the terminal via `/dev/tty`

### Requirement: Key binding semantics
ESC SHALL signal "go back" (`ui.ErrBack`) and Ctrl+C SHALL signal "cancel all" (`nil` error, no selection). Enter SHALL confirm the highlighted item.

#### Scenario: ESC signals back navigation
- **WHEN** the user presses ESC in any `ui.Choose()` list
- **THEN** `ui.Choose()` SHALL return `("", ui.ErrBack)`

#### Scenario: Ctrl+C signals full cancel
- **WHEN** the user presses Ctrl+C in any `ui.Choose()` list
- **THEN** `ui.Choose()` SHALL return `("", nil)` with no selection
