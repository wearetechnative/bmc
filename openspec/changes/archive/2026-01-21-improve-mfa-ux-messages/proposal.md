# Proposal: Improve MFA User Experience with Better Messages

## Why

The MFA authentication flow in BMC displays debug-oriented messages that are confusing for end users:

1. **Debug output**: Prints `bmc` (from `echo $0`) which provides no useful information
2. **Technical jargon**: Shows `sourceProfile technative` using developer terminology instead of user-friendly language
3. **Raw boolean**: Displays `MFA: true` which is implementation detail rather than user information
4. **Command echo**: Prints the full `aws-mfa` command with all flags and ARNs, creating visual noise
5. **No script feedback**: When TOTP script executes, users don't know what's happening until it completes
6. **False success message**: Shows "Copied to clipboard" even when clipboard command fails

These messages make the MFA flow feel like debugging output rather than a polished user experience, and can mislead users about whether operations succeeded.

## What Changes

Replace debug-oriented messages with clear, actionable user feedback:

### 1. Remove Debug Output
- Remove `echo $0` that prints "bmc" with no context
- Remove `MFA: ${mfa}` boolean flag display

### 2. Improve Message Clarity
- Change `sourceProfile technative` to `-- Using AWS source-profile: technative`
- Change raw command echo to `-- Refreshing MFA session for ${sourceProfile}...`
- Add `-- Executing TOTP script...` before running totpScript
- Validate clipboard success before showing confirmation

### 3. Conditional Messaging
- Only show "Copied to clipboard" when clipboard command succeeds
- Show helpful error when clipboard fails: "Note: Clipboard copy failed (command not found or error)"
- Suppress clipboard command stderr to avoid noise

## Impact

**Benefits:**
- **Professional UX**: Messages feel intentional and user-focused
- **Clear status**: Users understand what's happening at each step
- **Accurate feedback**: Success messages only shown when operations actually succeed
- **Reduced confusion**: Technical implementation details hidden from users
- **Better troubleshooting**: Error messages are helpful rather than just command failures

**Affected code:**
1. `bmc` line 12 - Remove `echo $0` debug output
2. `_bmclib.sh` lines 358-391 - Improve all MFA-related user messages
3. Clipboard validation - Check success before showing confirmation

**Affected specs:**
- Creates new `mfa-authentication` spec documenting MFA user experience requirements
