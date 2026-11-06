# Tasks

## 1. Already implemented in the working tree (verify)

- [x] 1.1 `ec2connect` connection method picker lists `ssm` before `ssh` in `cmd/ec2connect.go`; verify: `bmc ec2connect` without `-u`/`-i`/`-k` shows `ssm` highlighted first
- [x] 1.2 `ui.PrintTable` pads every cell with one space on both sides; verify: `AWS_PROFILE=TN-NonProduction bmc ec2ls` shows `│ InstanceId ... │`
- [x] 1.3 Add `ui.PrintText(columns, rows, headers)` and `ec2ls --txt` (mutually exclusive with `--json`); verify: `bmc ec2ls --txt` shows aligned columns without borders and `bmc ec2ls --txt --json` errors
- [x] 1.4 Add `ec2ls --no-headers` requiring `--txt`; verify: `bmc ec2ls --txt --no-headers` prints only instance lines and `bmc ec2ls --no-headers` errors with "--no-headers requires --txt"

## 2. Global debug flag

- [x] 2.1 Create `internal/debug` package with `Enabled bool` and `Printf(w io.Writer, format string, args ...any)` that writes only when enabled; verify with a unit test covering enabled and disabled
- [x] 2.2 Register persistent `--debug` flag on `rootCmd` in `cmd/root.go` bound to `debug.Enabled`; verify: `bmc ec2ls --debug --help` and `bmc --debug version` run without error

## 3. MFA messages behind debug

- [x] 3.1 In `internal/mfa/session.go`, route "-- Using AWS source-profile", "Current MFA Session Valid, until" and the TOTP code echo through `debug.Printf`; keep refresh, TOTP-script, manual-entry, device-not-found and refreshed-until messages unconditional; verify: `go test ./internal/mfa/...`
- [x] 3.2 Remove the `mfa.Quiet` variable and its assignment in `cmd/ec2ls.go`; verify: `grep -rn "mfa.Quiet\|Quiet bool" cmd internal` returns nothing and `go build ./...` succeeds
- [x] 3.3 Verify behaviour end-to-end: `AWS_PROFILE=TN-NonProduction bmc ec2ls --txt 2>&1 >/dev/null` prints nothing on stderr with a valid session, and with `--debug` prints the source-profile and session-valid lines

## 4. Documentation

- [x] 4.1 Document `--txt` and `--no-headers` in `docs/content/commands/ec2ls.md`
- [x] 4.2 Update `docs/content/commands/ec2ls.md`: remove the statement that `--no-headers` hides MFA messages; MFA status is hidden by default for all commands
- [x] 4.3 Document the global `--debug` flag in `docs/content/commands/_index.md` (what it shows and that output goes to stderr); verify the page renders in `hugo server` or by reading the markdown
- [x] 4.4 Mention in `docs/content/commands/ec2connect.md` that SSM is the default in the connection method picker
- [x] 4.5 Update `CHANGELOG.md` under `## NEXT VERSION`: add `--debug`, describe MFA status lines as hidden by default (output change), and correct the `--no-headers` entry

## 5. Final verification

- [x] 5.1 Run `gofmt -l cmd internal` (no new files listed), `go vet ./...` and `go test ./...`; all pass
- [x] 5.2 `openspec validate ec2ls-text-output-and-mfa-debug --strict` passes
