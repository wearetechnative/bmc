---
# bmc-m7qr
title: tui-colors-eval-context
status: in-progress
type: task
priority: normal
created_at: 2026-05-06T20:26:20Z
updated_at: 2026-05-06T20:26:20Z
---

Bij `eval "$(./bmc profsel)"` worden geen kleuren getoond in de TUI. De TUI schrijft correct naar `/dev/tty` via `tea.WithOutput(tty)`, maar de lipgloss default renderer kijkt naar `os.Stdout` voor kleurdetectie. Omdat stdout gepiped is in een `$()` subshell, detecteert termenv geen TTY en valt het terug op plain output.

Fix: na het openen van `/dev/tty`, de lipgloss default renderer koppelen aan die tty met `lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(tty))`. Dit speelt in `internal/ui/list.go`, `table.go` en `spinner.go`.
