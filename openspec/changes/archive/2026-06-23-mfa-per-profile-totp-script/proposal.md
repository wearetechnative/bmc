## Why

The global `totp_script` in `~/.config/bmc/config.json` is a single command applied to every AWS source profile. Users who manage multiple AWS accounts through separate identities (e.g. a company account and a personal account) need different TOTP credentials for each, but currently have no way to configure this — they either use the wrong OTP entry or must manually enter a code each time.

## What Changes

- Add `profile_scripts` map to `MFAConfig` in `internal/config/config.go`: maps source profile name → totp script command string
- Update `acquireTOTP()` in `internal/mfa/session.go` to accept `sourceProfile` and resolve the correct script using lookup order: `profile_scripts[sourceProfile]` → `totp_script` (global fallback) → manual input
- Pass `sourceProfile` from `EnsureValid()` down to `acquireTOTP()`
- Existing configs without `profile_scripts` continue to work unchanged

## Capabilities

### New Capabilities

- `mfa-per-profile-totp-script`: Per-source-profile TOTP script configuration — users can map individual source profile names to specific totp script commands, with the global `totp_script` as fallback

### Modified Capabilities

- `mfa-authentication`: The TOTP acquisition requirement changes — the system must now resolve the TOTP script per source profile before execution

## Impact

- `internal/config/config.go` — `MFAConfig` struct gains `ProfileScripts map[string]string`
- `internal/mfa/session.go` — `acquireTOTP()` signature change (adds `sourceProfile string`), lookup logic added
- No CLI changes, no new flags, no breaking changes to existing config
- Backwards compatible: existing single-script configs work without modification
