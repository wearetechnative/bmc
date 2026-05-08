## Why

The bmc configuration file is currently in TOML format (`~/.config/bmc/config.toml`). While TOML works well in Go, JSON is a more universally supported format that integrates better with external tooling, scripts, and editors without requiring TOML-specific support.

## What Changes

- **BREAKING**: Config file renamed from `~/.config/bmc/config.toml` to `~/.config/bmc/config.json`
- **BREAKING**: Config file format changes from TOML to JSON
- `config.Load()` reads JSON instead of TOML
- `bmc doctor` checks for `config.json` instead of `config.toml`
- README and documentation updated to show JSON examples
- Migration path: if `config.toml` exists and `config.json` does not, bmc prints a migration hint

## Capabilities

### New Capabilities

- `config-format-json`: Config file is read and written as JSON; existing TOML config is detected and a migration hint is shown

### Modified Capabilities

- `mfa-authentication`: Config loading path changes (same fields, different format)

## Impact

- `internal/config/config.go`: replace TOML library with `encoding/json` (stdlib)
- `go.mod` / `go.sum`: remove TOML dependency if no longer used elsewhere
- `README.md`: update all config examples from TOML to JSON syntax
- Users must migrate their `config.toml` to `config.json`

## Bean

[bmc-ueqs — config-toml2json](.beans/bmc-ueqs--config-toml2json.md)
