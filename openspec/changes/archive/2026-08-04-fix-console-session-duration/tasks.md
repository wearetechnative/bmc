## 1. Configuration

- [x] 1.1 Add `SessionDurationSeconds int` with a `session_duration_seconds` JSON tag to `ConsoleConfig` in `internal/config/config.go`
- [x] 1.2 Set the default to 3600 in `config.Defaults()`

## 2. Credential resolution

- [x] 2.1 Add `MinSessionDuration`, `MaxSessionDuration`, and `FallbackSessionDuration` constants to `internal/awsops/console.go`
- [x] 2.2 Add `sessionDuration` to read the config value and clamp it to 900–43200 seconds, returning one hour when unset, zero, or negative
- [x] 2.3 Extract credential resolution into `retrieveCredentials`, passing `awscfg.WithAssumeRoleCredentialOptions` to set `o.Duration` when the profile leaves it at zero or asks for more than the limit
- [x] 2.4 Add `isDurationRejected` to detect STS rejecting the requested `DurationSeconds`
- [x] 2.5 In `BuildFederationURL`, retry `retrieveCredentials` at one hour when the first attempt asked for more and was rejected on duration

## 3. Federation request

- [x] 3.1 Send `DurationSeconds` to `getSigninToken` only when the credentials carry no session token, with a comment explaining that AWS ignores it for temporary credentials

## 4. Reporting the expiry

- [x] 4.1 Change `OpenConsole` to return `(time.Time, error)` so callers receive the credential expiry
- [x] 4.2 In `cmd/console.go`, capture the expiry and print it to stderr when it is non-zero
- [x] 4.3 Give `registerConsoleSession` an `expiry` parameter, replacing the hardcoded `time.Now().Add(time.Hour)`, and fall back to `awsops.FallbackSessionDuration` when the expiry is zero

## 5. Watcher refresh scheduling

- [x] 5.1 Add `watcher.RefreshTime(expiry)` to `internal/watcher/state.go`, capping the refresh window at half the remaining lifetime
- [x] 5.2 Use `RefreshTime` in both `refreshSession` call sites in `internal/watcher/daemon.go` and in `registerConsoleSession`

## 6. Tests

- [x] 6.1 Add `internal/awsops/console_test.go` covering `sessionDuration` clamping, the guarantee that it never returns zero, and `isDurationRejected` against both STS message variants and unrelated errors
- [x] 6.2 Add `RefreshTime` tests to `internal/watcher/state_test.go` covering the full window, the halved window for short sessions, and an already-expired session
- [x] 6.3 Verify `go build ./...`, `go vet ./...`, and `go test ./...` all pass

## 7. Verification against live AWS

- [x] 7.1 Confirm the old code path yields ~15 minutes and the new one ~60 minutes for a real role-based profile
- [x] 7.2 Confirm a 12-hour request against a role capped at one hour is rejected, matched by `isDurationRejected`, and falls back to a working one-hour session

## 8. Documentation

- [x] 8.1 Add a "Session lifetime" section to `docs/content/commands/console.md` documenting the default, the `session_duration_seconds` option, its bounds, the `MaxSessionDuration` fallback, and the role-chaining cap
