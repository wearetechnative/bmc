## Why

Bij `eval "$(./bmc profsel)"` is stdout gepiped. De lipgloss default renderer kijkt naar `os.Stdout` voor kleurdetectie — detecteert geen TTY — en toont plain output zonder kleuren, ondanks dat de TUI correct naar `/dev/tty` schrijft.

This task is tracked in [bmc-m7qr](../../../.beans/bmc-m7qr--tui-colors-eval-context.md).

## What Changes

- **`internal/ui/list.go`**: Na het openen van `/dev/tty`, de lipgloss default renderer koppelen aan die tty via `lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(tty))`
- **`internal/ui/table.go`**: Zelfde fix
- **`internal/ui/spinner.go`**: Spinner gebruikt `os.Stderr` als output — ook koppelen aan tty wanneer beschikbaar

## Capabilities

### New Capabilities

- `tui-color-rendering`: TUI kleuren werken correct in eval/$() context door lipgloss renderer te koppelen aan /dev/tty

### Modified Capabilities

## Impact

- `internal/ui/list.go`, `table.go`, `spinner.go`
- Geen breaking changes — gedrag verbetert alleen in eval-context
