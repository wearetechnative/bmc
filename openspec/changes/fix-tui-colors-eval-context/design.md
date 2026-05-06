## Context

bmc opent `/dev/tty` expliciet voor TUI I/O zodat de interface werkt wanneer stdout gepiped is (`eval "$(./bmc profsel)"`). Maar lipgloss v1.x gebruikt een global default renderer die kijkt naar `os.Stdout` voor kleurdetectie. Wanneer stdout niet een TTY is, valt termenv terug op plain output.

`lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(tty))` koppelt de renderer aan de tty fd, waarna alle lipgloss styles kleuren renderen via die tty.

## Goals / Non-Goals

**Goals:**
- Kleuren en styling werken in `eval "$(./bmc profsel)"` en andere piped contexten
- Fix geldt voor alle TUI componenten: list, table, spinner

**Non-Goals:**
- Wijzigingen aan terminal capability detection zelf
- Windows-ondersteuning (bmc target linux/darwin)

## Decisions

### Decision: SetDefaultRenderer als globale override

Na het openen van `/dev/tty`, direct `lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(tty))` aanroepen vóór het aanmaken van de tea.Program.

**Rationale**: Package-level style globals (`var tableStyle = lipgloss.NewStyle()...`) zijn al aangemaakt met de default renderer. `SetDefaultRenderer` overschrijft de global zodat alle rendering via de tty loopt — ook die van bestaande styles.

**Alternatief overwogen**: `renderer.NewStyle()` gebruiken voor alle styles. Afgewezen — vereist refactor van alle package-level globals en alle stijl-aanmaak verplaatsen naar binnen de functies.

### Decision: spinner ook naar /dev/tty

De spinner gebruikt nu `tea.WithOutput(os.Stderr)`. In piped context is stderr wel een TTY (niet gepiped), dus kleuren werken toevallig. Toch consistent maken: ook de spinner de tty laten gebruiken wanneer beschikbaar, zodat het gedrag uniform is.

## Risks / Trade-offs

- `SetDefaultRenderer` is een global state change. Als meerdere goroutines gelijktijdig TUI componenten starten, kan dit racen. bmc is single-threaded in TUI gebruik — geen risico.
- Na de program run blijft de renderer ingesteld op de tty. De tty wordt dan gesloten. Verdere lipgloss rendering na de TUI (bijv. `PrintTable`) zou dan een gesloten fd kunnen gebruiken. Oplossing: renderer pas instellen vlak voor `p.Run()`, en de tty open houden totdat de gehele functie klaar is (dit is al het geval).
