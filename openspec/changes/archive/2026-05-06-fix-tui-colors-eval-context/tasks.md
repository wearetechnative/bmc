## 1. Fix lipgloss renderer koppeling

- [x] 1.1 `internal/ui/list.go`: voeg `lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(tty))` toe direct na het openen van `/dev/tty` in `Choose()`
- [x] 1.2 `internal/ui/table.go`: zelfde fix in `runTable()`
- [x] 1.3 `internal/ui/spinner.go`: open `/dev/tty` voor spinner output (ipv `os.Stderr`) en stel renderer in

## 2. Verificatie

- [x] 2.1 Bouw binary: `go build -o bmc.test`
- [x] 2.2 Test kleuren in eval context: `eval "$(./bmc.test profsel)"` — TUI moet kleuren tonen
- [x] 2.3 Verwijder test binary
