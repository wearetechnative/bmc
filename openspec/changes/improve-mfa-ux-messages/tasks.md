# Implementation Tasks

## 1. Remove debug output from bmc script
- [x] Remove `echo $0` statement on line 12 of bmc script
- [x] This eliminates the confusing "bmc" output at script start
- [x] **Validation**: Run bmc profsel and verify "bmc" no longer appears
- **Dependency**: None (standalone change)

## 2. Improve source profile message
- [x] Change `echo "sourceProfile $sourceProfile"` to `echo "-- Using AWS source-profile: $sourceProfile"`
- [x] Makes message clear and user-friendly with consistent `--` prefix
- [x] **Validation**: Verify message displays correctly during profile selection
- **Dependency**: None (can be done in parallel with task 1)

## 3. Remove MFA boolean debug message
- [x] Remove `echo "MFA: ${mfa}"` on line 369 of _bmclib.sh
- [x] Users don't need to see internal boolean state
- [x] **Validation**: Verify MFA flow works without displaying boolean
- **Dependency**: None (can be done in parallel)

## 4. Replace command echo with user-friendly message
- [x] Change `echo aws-mfa --profile...` to `echo "-- Refreshing MFA session for ${sourceProfile}..."`
- [x] Hides technical command details, shows clear action status
- [x] **Validation**: Verify message appears when MFA session needs refresh
- **Dependency**: None (can be done in parallel)

## 5. Add TOTP script execution message
- [x] Add `echo "-- Executing TOTP script..."` before totpScript execution
- [x] Provides feedback while script runs (especially important for slow password managers)
- [x] **Validation**: Verify message appears before TOTP script runs
- **Dependency**: None (can be done in parallel)

## 6. Improve clipboard validation and messaging
- [x] Check clipboard command exit status before showing success
- [x] Display "-- Copied to clipboard" only on successful copy
- [x] Display "-- Note: Clipboard copy failed..." on error
- [x] Suppress stderr with `2>/dev/null` to avoid noise
- [x] **Validation**: Test with working and broken clipboard commands
- **Dependency**: Requires understanding of clipboard flow (task 5)

## 7. Update documentation
- [x] Update CHANGELOG.md with user experience improvements
- [x] Document the improved message flow
- [x] **Validation**: Documentation is clear and accurate
- **Dependency**: Can be done after tasks 1-6 or in parallel

## 8. Test complete MFA flow
- [x] Test with TOTP script configured → should show all new messages
- [x] Test with working clipboard → should show success message
- [x] Test with broken clipboard → should show failure message
- [x] Test without TOTP script → should show manual entry message
- [x] Verify no debug output appears
- [x] **Validation**: All scenarios display appropriate user-friendly messages
- **Dependency**: Requires all tasks 1-6 complete

## Notes

- Tasks 1-5 are independent and can be done in parallel
- Task 6 builds on clipboard handling from previous TOTP script fixes
- Task 7 can be done anytime
- Task 8 is the final verification step
- All changes maintain backward compatibility
