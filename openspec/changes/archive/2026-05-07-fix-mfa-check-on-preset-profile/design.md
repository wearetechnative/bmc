## Context

`ensureAWSProfile()` in `cmd/profilehelper.go` has two code paths:

1. **Environment path**: `AWS_PROFILE` already set → return immediately (no MFA check)
2. **Interactive path**: no profile in env → interactive selection → `mfa.EnsureValid()` called

`mfa.EnsureValid()` is safe to call at any time — it is a no-op when MFA is disabled, when the session is still valid, or when no MFA device is configured. The missing call on the environment path causes expired sessions to go undetected until AWS returns `InvalidClientTokenId`.

## Goals / Non-Goals

**Goals:**
- Call `mfa.EnsureValid()` on both code paths
- No change in behavior when MFA is not configured or session is still valid
- No additional prompts or latency when the session is valid

**Non-Goals:**
- Changing MFA session duration or credential storage
- Handling profiles without a source profile (no role assumption)

## Decisions

### Decision: resolve source profile even for pre-set profiles

To call `mfa.EnsureValid(sourceProfile, ...)`, we need the source profile name — not the role profile. The interactive path already calls `awsconfig.ResolveSourceProfile()` for this. We apply the same resolution on the environment path.

**Steps on the environment path:**
1. Load all profiles via `awsconfig.LoadProfiles()`
2. Find the profile matching `AWS_PROFILE`
3. Call `awsconfig.ResolveSourceProfile(profile)` to get the underlying credential profile
4. Load config and call `mfa.EnsureValid(sourceProfile, cfg, os.Stderr)`

**Alternative considered**: call `mfa.EnsureValid` with the role profile directly. Rejected — `EnsureValid` needs the source profile to load long-term credentials and the MFA device ARN.

### Decision: non-fatal if profile not found in config

If `AWS_PROFILE` is set to a profile that doesn't exist in `~/.aws/config` (unusual but possible), log a warning and continue rather than returning an error. The subsequent AWS call will fail with a clear message from the SDK.

## Risks / Trade-offs

- **Extra file read on startup**: `awsconfig.LoadProfiles()` reads `~/.aws/config`. This is already called on the interactive path; the added cost on the environment path is negligible.
- **MFA prompt on pre-set profile**: If the session is expired, the user will now see a TOTP prompt even when `AWS_PROFILE` was set externally. This is the desired behavior — previously it just failed with a cryptic AWS error.

## Migration Plan

Single-file change in `cmd/profilehelper.go`. No migration needed.
