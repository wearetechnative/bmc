## 1. History package

- [ ] 1.1 Create `internal/history/history.go` with `Load(name string) []string` — reads `~/.local/share/bmc/<name>-history.json`, returns empty slice on any error
- [ ] 1.2 Implement `Save(name string, entry string) error` — prepends entry, deduplicates, caps at 10, writes atomically (write to temp file, rename)
- [ ] 1.3 Resolve XDG data dir: use `$XDG_DATA_HOME` if set, otherwise `~/.local/share`

## 2. Console command integration

- [ ] 2.1 In `cmd/console.go`, before calling `selectProfileInteractive`, load console history via `history.Load("console")`
- [ ] 2.2 Build the profile item list with recent profiles prepended (Desc: "recent") and deduped from the full list
- [ ] 2.3 Call `ui.Choose` directly with the reordered list instead of `selectProfileInteractive` (console-specific selector)
- [ ] 2.4 After `awsops.OpenConsole` returns without error, call `history.Save("console", selectedProfile.Name)`

## 3. Verification

- [ ] 3.1 First run with no history file: selector shows normal list, history file created after successful open
- [ ] 3.2 Second run: previously used profile appears at top with "recent" label
- [ ] 3.3 Using `-p` flag or `AWS_PROFILE` set: history file not modified
- [ ] 3.4 History capped at 10: opening an 11th distinct profile drops the oldest entry
