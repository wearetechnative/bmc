## 1. Implementation

- [x] 1.1 Modify `profsel` function in `bmc` to check if `selectedProfileName` is empty before calling `setMFA`
- [x] 1.2 Add early return in `profsel` when profile selection is cancelled or empty
- [x] 1.3 Update `selectAWSProfile` in `_bmclib.sh` to detect and handle gum filter cancellation (exit code check)
- [x] 1.4 Review and update `setMFA` function to use `return` instead of `exit` when called from sourced context, or ensure it's never called with empty sourceProfile
- [x] 1.5 Add user feedback message when selection is cancelled
- [x] 1.6 Test cancellation behavior in both sourced (`. bmc profsel`) and executed (`bmc profsel`) modes
- [x] 1.7 Test with Ctrl-C interruption during profile group selection
- [x] 1.8 Test with Ctrl-C interruption during profile selection
- [x] 1.9 Test in both bash and zsh shells

## 2. Documentation

- [x] 2.1 Update CHANGELOG.md with fix description under "NEXT VERSION"
- [x] 2.2 Verify README accurately describes profsel behavior
