## 1. Fix ensureAWSProfile in cmd/profilehelper.go

- [x] 1.1 When `AWS_PROFILE` is already set, load profiles via `awsconfig.LoadProfiles()` and find the matching profile
- [x] 1.2 Call `awsconfig.ResolveSourceProfile(profile)` to get the source credential profile
- [x] 1.3 Load config via `config.Load()` and call `mfa.EnsureValid(sourceProfile, cfg, os.Stderr)`
- [x] 1.4 If the profile is not found in `~/.aws/config`, log a warning to stderr and continue without MFA check

## 2. Verification

- [ ] 2.1 Test with valid MFA session: verify no TOTP prompt and session valid message is shown
- [ ] 2.2 Test with expired MFA session: verify TOTP prompt appears and session is refreshed before AWS call
- [ ] 2.3 Test with MFA disabled in config: verify no MFA prompt and no errors
