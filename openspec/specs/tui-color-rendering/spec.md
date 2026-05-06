## Requirements

### Requirement: TUI kleuren werken in piped stdout context
De TUI SHALL kleuren en styling renderen wanneer stdout niet een TTY is, zolang `/dev/tty` beschikbaar is.

#### Scenario: eval context met /dev/tty beschikbaar
- **WHEN** `bmc profsel` wordt uitgevoerd via `eval "$(./bmc profsel)"`
- **THEN** de TUI toont kleuren, borders en highlighting via `/dev/tty`

#### Scenario: /dev/tty niet beschikbaar
- **WHEN** `/dev/tty` niet geopend kan worden
- **THEN** de TUI valt terug op plain output zonder kleuren (bestaand gedrag)
