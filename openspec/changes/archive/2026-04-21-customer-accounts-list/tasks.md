## 1. Core Function

- [x] 1.1 Add `accountListAccounts()` function to `_bmclib.sh` that takes a profile group name, queries `jsonify-aws-dotfiles` JSON, filters profiles by group, and prints a table (Profile Name, Account ID, Role Name) using `gum table`
- [x] 1.2 Add `accountSelectGroup()` helper or inline logic that presents unique profile groups via `gum filter` and returns the selected group name

## 2. Command Registration

- [x] 2.1 Register `accountls` subcommand in `bmc` dispatcher with `make_command` and description
- [x] 2.2 Implement the `accountls` command function that calls group selection then account listing, handling cancellation (non-zero exit)

## 3. Edge Cases

- [x] 3.1 Handle empty result when selected group has no profiles (display informational message)
- [x] 3.2 Handle user cancellation during `gum filter` (exit cleanly with non-zero code)

## 4. Verification

- [x] 4.1 Test `accountls` appears in `bmc usage` output
- [x] 4.2 Test end-to-end: run `bmc accountls`, select a group, verify table output
