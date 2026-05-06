## Why

After the Go rewrite (bmc-rewrite-in-go) several UX bugs and missing behaviours were found in daily use. These are small, focused fixes that improve the first-run experience and correct visual regressions in the TUI.

Relates to bean: `.beans/bmc-p4wn--post-rewrite-issues.md`

## What Changes

- **TUI list height**: fix blank rows between last item and help bar in short lists (`ec2connect` connection method, SSH user picker)
- **profsel missing wrapper hint**: when `bmc profsel` is run without the shell wrapper installed, warn the user on stderr instead of silently printing an unused export
- **install-shell-integration permission denied**: when the rc file is not writable (e.g. home-manager on NixOS), show a helpful fallback with manual snippets for home-manager, bash/zsh, and Fish instead of a generic error
- **console containerized tab**: `bmc console` opens the AWS console but lacks the isolated browser tab per profile that `assumego` (Granted) provided — this capability is missing entirely

## Capabilities

### New Capabilities
- none

### Modified Capabilities
- `tui-list-display`: height and help-bar rules for `ui.Choose()`
- `shell-integration`: behaviour of `profsel` output and `install-shell-integration` error handling

## Impact

- `internal/ui/list.go` — `Choose()` height formula and `SetShowHelp` logic
- `cmd/profsel.go` — TTY detection, stderr hint when wrapper is absent
- `cmd/install_shell.go` — `permission denied` fallback with manual snippets
- `cmd/console.go` — containerized tab support (scope TBD)
- `README.md` — NixOS shell integration documentation
