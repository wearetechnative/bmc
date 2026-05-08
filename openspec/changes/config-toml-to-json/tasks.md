## 1. Config package

- [ ] 1.1 In `internal/config/config.go`: replace `toml:"..."` struct tags with `json:"..."`
- [ ] 1.2 Replace `github.com/BurntSushi/toml` import with `encoding/json`
- [ ] 1.3 Update `ConfigPath()` to return `~/.config/bmc/config.json`
- [ ] 1.4 Update `Load()` to parse JSON with `json.Unmarshal`
- [ ] 1.5 Add migration hint in `Load()`: if `config.json` absent and `config.toml` present, print hint to stderr

## 2. Dependency cleanup

- [ ] 2.1 Run `go mod tidy` to remove `BurntSushi/toml` from `go.mod` and `go.sum`
- [ ] 2.2 Verify no other package imports `BurntSushi/toml`

## 3. Documentation

- [ ] 3.1 Update `README.md`: replace all TOML config examples with JSON equivalents
- [ ] 3.2 Update `CHANGELOG.md` under `## NEXT VERSION`

## 4. Verification

- [ ] 4.1 `bmc doctor` reports correct config file path (`config.json`)
- [ ] 4.2 With valid `config.json`: settings are applied correctly
- [ ] 4.3 With absent `config.json` and present `config.toml`: migration hint shown, defaults used
- [ ] 4.4 With absent both files: no error, defaults used
