## MODIFIED Requirements

### Requirement: TUI colors work in piped stdout context
The TUI SHALL render colors and styling when stdout is not a TTY, as long as `/dev/tty` is available.

#### Scenario: eval context with /dev/tty available
- **WHEN** `bmc profsel` is executed via `eval "$(./bmc profsel)"`
- **THEN** the TUI renders colors, borders, and highlighting via `/dev/tty`

#### Scenario: /dev/tty not available
- **WHEN** `/dev/tty` cannot be opened
- **THEN** the TUI falls back to plain output without colors (existing behavior)
