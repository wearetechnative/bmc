## 1. History package

- [x] 1.1 Create `internal/history/history.go` with `Load(name string) []string` — reads `~/.local/share/bmc/<name>-history.json`, returns empty slice on any error
- [x] 1.2 Implement `Save(name string, entry string) error` — prepends entry, deduplicates, caps at 10, writes atomically (write to temp file, rename)
- [x] 1.3 Resolve XDG data dir: use `$XDG_DATA_HOME` if set, otherwise `~/.local/share`

## 2. Console command integration

- [x] 2.1 In `cmd/console.go`, before calling `selectProfileInteractive`, load console history via `history.Load("console")`
- [x] 2.2 Build the profile item list with recent profiles prepended (Desc: "recent") and deduped from the full list
- [x] 2.3 Call `ui.Choose` directly with the reordered list instead of `selectProfileInteractive` (console-specific selector)
- [x] 2.4 After `awsops.OpenConsole` returns without error, call `history.Save("console", selectedProfile.Name)`

## 3. Verification

- [x] 3.1 First run with no history file: selector shows normal list, history file created after successful open
- [x] 3.2 Second run: previously used profile appears at top with "recent" label
- [x] 3.3 Using `-p` flag or `AWS_PROFILE` set: history file not modified
- [x] 3.4 History capped at 10: opening an 11th distinct profile drops the oldest entry
