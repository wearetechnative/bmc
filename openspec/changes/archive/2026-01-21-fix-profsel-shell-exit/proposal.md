# Change: Fix profsel shell exit on cancellation

## Why
When `bmc profsel` is sourced and the user interrupts the profile selection (Ctrl-C) or doesn't select a profile, the `setMFA` function calls `exit 1`, which causes the current shell to exit. This is disruptive and prevents users from gracefully cancelling the operation.

## What Changes
- Add graceful cancellation handling to `profsel` command
- Prevent `setMFA` from being called when no profile is selected
- Ensure the sourced script uses `return` instead of `exit` to avoid closing the user's shell
- Provide clear feedback when profile selection is cancelled

## Impact
- Affected specs: profile-selection (new capability)
- Affected code: `bmc` (profsel function at line 91-129), `_bmclib.sh` (selectAWSProfile function at line 306-344, setMFA function at line 346-383)
- User experience: Users can now safely cancel profile selection without losing their shell session
