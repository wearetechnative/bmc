## 1. Refactor shared selector into profilehelper.go

- [x] 1.1 Add `import "github.com/wearetechnative/bmc/internal/history"` to `cmd/profilehelper.go`
- [x] 1.2 Copy `recentGroups()` helper from `cmd/console.go` into `cmd/profilehelper.go`
- [x] 1.3 Add `selectProfileWithHistory(profiles []awsconfig.Profile) (awsconfig.Profile, bool, error)` to `cmd/profilehelper.go`, returning `(profile, wasInteractive, error)` — porting the full logic from `selectProfileForConsoleInteractive()` but using history key `"profile"`
- [x] 1.4 Delete `selectProfileInteractive()` from `cmd/profilehelper.go`

## 2. Update ensureAWSProfile()

- [x] 2.1 Replace the two `selectProfileInteractive()` call sites in `ensureAWSProfile()` with `selectProfileWithHistory()`, using the returned `wasInteractive` bool to gate `history.Save("profile", ...)`

## 3. Update console.go

- [x] 3.1 Replace `selectProfileForConsoleInteractive()` call in `runConsole()` with `selectProfileWithHistory()`
- [x] 3.2 Remove `recentGroups()` and `selectProfileForConsoleInteractive()` from `cmd/console.go`
- [x] 3.3 Change history save call in `runConsole()` from `history.Save("console", ...)` to `history.Save("profile", ...)`
- [x] 3.4 Remove unused `history` import from `cmd/console.go` if no longer needed (or leave if still used for Load — verify)

## 4. Update profsel.go

- [x] 4.1 Replace `selectProfileInteractive()` call in `runProfsel()` with `selectProfileWithHistory()`
- [x] 4.2 Add `history.Save("profile", selectedProfile.Name)` after a successful interactive selection in `runProfsel()`
- [x] 4.3 Add `import "github.com/wearetechnative/bmc/internal/history"` to `cmd/profsel.go`

## 5. Verify and clean up

- [x] 5.1 Run `go build ./...` and fix any compile errors
- [x] 5.2 Manually test `bmc profsel` interactive — verify recent profiles surface after first selection
- [x] 5.3 Manually test `bmc console` interactive — verify recent profiles still surface using the shared history
- [x] 5.4 Manually test a command using `ensureAWSProfile()` (e.g. `bmc ec2ls`) — verify recent profiles surface after first interactive selection
