## Why

When `AWS_PROFILE` is already set in the shell environment, `ensureAWSProfile()` returns early without calling `mfa.EnsureValid()`. This causes `InvalidClientTokenId` errors from AWS when the MFA session has expired, with no helpful feedback to the user.

## What Changes

- `ensureAWSProfile()` in `cmd/profilehelper.go` will call `mfa.EnsureValid()` even when `AWS_PROFILE` is already set in the environment, so expired MFA sessions are detected and refreshed before any AWS operation

## Capabilities

### New Capabilities

### Modified Capabilities

- `mfa-authentication`: MFA session validity is now checked regardless of how the profile was set (interactively or via environment variable)

## Impact

- `cmd/profilehelper.go`: load profile from awsconfig when `AWS_PROFILE` is pre-set, resolve source profile, and call `mfa.EnsureValid()`
- Bean: none — this is a standalone bug fix
