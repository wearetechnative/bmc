## 1. Update config struct

- [x] 1.1 In `internal/config/config.go`: rename `ClipboardCommand string` → `CopyCommand string` with JSON tag `copy_command`
- [x] 1.2 Add `PasteCommand string` field with JSON tag `paste_command` to `MFAConfig`

## 2. Update MFA session logic

- [x] 2.1 In `internal/mfa/session.go`: update `acquireTOTP` to use `cfg.MFA.CopyCommand` instead of `cfg.MFA.ClipboardCommand`
- [x] 2.2 Update `copyToClipboard` signature/name to `copyAndPaste(copyCmd, pasteCmd, code string, outfd *os.File)`
- [x] 2.3 After successful copy: `time.Sleep(300 * time.Millisecond)`, then run `pasteCmd` if non-empty (no stdin)
- [x] 2.4 On paste success: print `-- Pasted to active window`; on paste failure: print `-- Note: Paste failed (...)`

## 3. Update doctor check

- [x] 3.1 In `cmd/doctor.go`: update `checkClipboard` to read `cfg.MFA.CopyCommand` instead of `cfg.MFA.ClipboardCommand`
- [x] 3.2 Update the doctor check label/message from `clipboard_command` to `copy_command`

## 4. Fix totp_script execution

- [x] 4.1 Replace `strings.Fields` + `exec.Command(parts[0], parts[1:]...)` with `exec.Command("sh", "-c", script)` in `runTOTPScript`
- [x] 4.2 Set `cmd.Stdin = os.Stdin` and `cmd.Stderr = os.Stderr` so interactive TUI scripts (e.g. gum filter) work correctly

## 5. Update documentation and verify

- [x] 5.1 Update `README.md` config reference table: replace `mfa.clipboard_command` with `mfa.copy_command` and add `mfa.paste_command`
- [x] 5.2 Add breaking change entry to `CHANGELOG.md` under `## NEXT VERSION`
- [x] 5.3 Confirm `go build ./...` succeeds
