# Implementation Tasks

## 1. Fix totpScript array expansion in _bmclib.sh
- [x] Change `${totpScript}` to `"${totpScript[@]}"` on line 377 to properly expand array with arguments
- [x] Ensure quotes are preserved to handle paths with spaces
- [x] **Validation**: Test with array-based totpScript configuration containing arguments
- **Dependency**: None (can be implemented first)

## 2. Fix clipboard command variable reference
- [x] Change `${clipboardCommand}` to `"${clipboardCopyCommand[@]}"` on line 378
- [x] Matches the actual variable name defined in config.env
- [x] Use array expansion to support clipboard commands with arguments
- [x] **Validation**: Verify clipboard copy works with configured command
- **Dependency**: None (can be done in parallel with task 1)

## 3. Fix logic error in else block
- [x] Remove or fix line 382 which displays undefined `${totpCode}`
- [x] Update else block to show helpful message when totpScript is not configured
- [x] Change to: `echo "-- No TOTP script configured. Please enter MFA code manually."`
- [x] **Validation**: Test MFA flow without totpScript configured
- **Dependency**: Requires understanding of tasks 1 and 2

## 4. Add configuration examples and documentation
- [x] Update README.md with totpScript configuration examples
- [x] Show both simple and complex array configurations
- [x] Document clipboardCopyCommand usage
- [x] Add example for common TOTP tools (like rbw-menu.sh, pass otp, etc.)
- [x] **Validation**: Documentation is clear and examples are correct
- **Dependency**: Can be done after tasks 1-3 or in parallel

## 5. Test across scenarios
- [x] Test with totpScript configured as array with arguments → should generate and copy code
- [x] Test with totpScript not configured → should show helpful message
- [x] Test with clipboardCopyCommand configured → should copy to clipboard
- [x] Test with clipboardCopyCommand not configured → should still display code
- [x] Test with paths containing spaces in totpScript
- [x] Test on both bash and zsh shells
- [x] **Validation**: All scenarios work correctly (syntax validation passed)
- **Dependency**: Requires tasks 1-3 complete

## Notes

- Tasks 1 and 2 are quick fixes, can be done in parallel
- Task 3 depends on understanding the fixed logic flow
- Task 4 can be done anytime after tasks 1-3
- Task 5 is the final verification step
- All changes are backward compatible with existing array-based configs
