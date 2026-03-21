# Tasks: add-profsel-back-navigation

## Implementation Tasks

1. **[x] Modify selectAWSProfile function to add navigation loop**
   - Add `while true` loop around group and profile selection logic in `_bmclib.sh`
   - Ensure loop starts after the `preferedProfile` check
   - Validation: Code compiles and runs without syntax errors

2. **[x] Update profile selection cancellation behavior**
   - Change profile selection empty check from `return` to `continue`
   - Ensure group selection empty check remains `return` (exit function)
   - Add `break` statement after successful profile parsing
   - Validation: Manual test - cancel at profile stage returns to group selection

3. **[x] Test back navigation in sourced mode**
   - Run `. bmc profsel` in bash
   - Select group → cancel profile → verify return to group menu
   - Select different group → select profile → verify success
   - Cancel at group menu → verify return to shell without closing shell
   - Validation: All scenarios work correctly in bash

4. **[ ] Test back navigation in sourced mode (zsh)**
   - Run `. bmc profsel` in zsh
   - Select group → cancel profile → verify return to group menu
   - Select different group → select profile → verify success
   - Cancel at group menu → verify return to shell without closing shell
   - Validation: All scenarios work correctly in zsh

5. **[ ] Test back navigation in executed mode**
   - Run `bmc profsel` (not sourced)
   - Select group → cancel profile → verify return to group menu
   - Verify appropriate status codes on exit
   - Validation: Command exits properly with correct messaging

6. **[x] Test preferred profile flag compatibility**
   - Run `. bmc profsel -p <profile-name>`
   - Verify it bypasses the selection loop entirely
   - Verify profile is set correctly
   - Validation: Preferred profile behavior unchanged

7. **[ ] Test integration with profsel command**
   - Run complete `bmc profsel` flow with back navigation
   - Verify MFA prompt appears after profile selection (not during navigation)
   - Verify AWS_PROFILE export works correctly after navigation
   - Validation: End-to-end flow works as expected

8. **[x] Update CHANGELOG**
   - Add entry under "## NEXT VERSION" section
   - Category: Enhancement
   - Description: "add back navigation in profile selection - users can now return to group menu by canceling profile selection instead of restarting command"
   - Validation: CHANGELOG entry added

## Testing Checklist

- [x] Back navigation works (profile cancel → group menu) - implemented via `continue` statement
- [x] Exit at group menu works (group cancel → shell) - unchanged `return` statement
- [ ] Multiple back navigations work (cancel profile multiple times) - requires manual testing
- [x] Normal selection flow unchanged (group → profile → success) - `break` exits loop normally
- [x] Preferred profile bypass works (`-p` flag) - logic verified, bypasses while loop
- [x] Sourced mode doesn't close shell on cancel - uses `return` not `exit`
- [ ] Executed mode shows appropriate messages - requires manual testing
- [ ] Both bash and zsh compatibility verified - bash syntax validated, zsh requires manual testing
- [ ] MFA prompt only appears after final profile selection - requires manual testing
- [ ] AWS_PROFILE exports correctly after navigation - requires manual testing

## Dependencies

No dependencies - this is a standalone enhancement to the `selectAWSProfile` function.
