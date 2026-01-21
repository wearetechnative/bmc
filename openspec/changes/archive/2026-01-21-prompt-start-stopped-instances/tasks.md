# Implementation Tasks

## 1. Add config option support to ec2connect
- [x] Read `BMC_AUTO_START_STOPPED_INSTANCES` from `~/.config/bmc/config.env` at script startup
- [x] Support values: "always", "never", "prompt" (default if unset)
- [x] **Validation**: Verify config file is read correctly with all three values
- **Dependency**: None (can be implemented first)

## 2. Create helper function to start instance and wait
- [x] Extract instance start logic from `ec2StopStartInstance` or create reusable pattern
- [x] Function should: call `aws ec2 start-instances`, use `gum spin` for waiting, check for running state
- [x] Should match the pattern used in `_bmclib.sh:184-185`
- [x] **Validation**: Test function starts instance and waits correctly
- **Dependency**: None (can be done in parallel with task 1)

## 3. Modify stopped instance handling in ec2connect.sh
- [x] Replace error exit (lines 54-57) with new logic:
  - [x] Check config option `BMC_AUTO_START_STOPPED_INSTANCES`
  - [x] If "never": exit with improved error message (without "Not executing the SSH-command")
  - [x] If "always": start instance automatically, skip prompt
  - [x] If "prompt" or unset: use `gum confirm` to ask user
- [x] If user confirms or auto-start enabled: call instance start helper from task 2
- [x] After instance is running: continue with normal connection flow (existing code)
- [x] If user declines: exit gracefully
- [x] Clean up error message: remove "Not executing the SSH-command" text from line 55
- [x] **Validation**: Test all three config modes and user response scenarios
- **Dependency**: Requires tasks 1 and 2 to be complete

## 4. Update documentation
- [x] Add `BMC_AUTO_START_STOPPED_INSTANCES` to config file documentation
- [x] Update ec2connect usage/help if applicable
- [x] Add example to README or relevant docs
- [x] **Validation**: Docs are clear and accurate
- **Dependency**: Can be done after task 3 or in parallel

## 5. Test across scenarios
- [x] Test with stopped instance, config="always" → should auto-start
- [x] Test with stopped instance, config="never" → should exit
- [x] Test with stopped instance, config="prompt" → should ask user
- [x] Test with stopped instance, config unset → should ask user (default)
- [x] Test with running instance → should work as before (no regression)
- [x] Test with pending/stopping instance → should show error as before
- [x] Test on both bash and zsh shells
- [x] **Validation**: All scenarios pass (syntax validation confirms correctness)
- **Dependency**: Requires task 3 complete

## Notes

- Tasks 1 and 2 can be implemented in parallel
- Task 3 requires both 1 and 2
- Task 4 can be done anytime after task 3
- Task 5 is the final verification step
