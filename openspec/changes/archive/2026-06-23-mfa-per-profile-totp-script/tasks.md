## 1. Config

- [x] 1.1 Add `ProfileScripts map[string]string` field with `json:"profile_scripts,omitempty"` to `MFAConfig` in `internal/config/config.go`

## 2. MFA Session Logic

- [x] 2.1 Add `sourceProfile string` parameter to `acquireTOTP()` signature in `internal/mfa/session.go`
- [x] 2.2 Implement profile-script lookup in `acquireTOTP()`: check `cfg.MFA.ProfileScripts[sourceProfile]` first, then fall back to `cfg.MFA.TOTPScript`
- [x] 2.3 Update the call site in `EnsureValid()` to pass `sourceProfile` to `acquireTOTP()`

## 3. Tests

- [x] 3.1 Add unit test for profile-specific script lookup (profile entry exists)
- [x] 3.2 Add unit test for global fallback when profile not in map
- [x] 3.3 Add unit test for nil/empty `ProfileScripts` map (backwards compatibility)

## 4. Documentation

- [x] 4.1 Update `bmc doctor` output to mention `profile_scripts` if configured (or note it is unconfigured)
- [x] 4.2 Update `docs/content/` with config example showing `profile_scripts` usage
