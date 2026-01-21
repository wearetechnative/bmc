# Proposal: Fix TOTP Script Handling and Clipboard Integration

## Problem

The current TOTP script functionality in BMC has several implementation issues that prevent it from working correctly:

1. **Array handling broken**: `totpScript` is defined as an array in `config.env` (e.g., `totpScript=("/path/to/script" "-t" "code")`) but executed as `${totpScript}` in `_bmclib.sh:377`, which doesn't properly expand arrays with arguments
2. **Undefined variable**: Line 378 references `${clipboardCommand}` which is not defined in config.env; the config file defines `clipboardCopyCommand` and `clipboardPasteCommand` instead
3. **Logic error**: Line 382 displays `"Code: ${totpCode}"` in the else block (when totpScript is not set), but `totpCode` is only defined inside the if block (line 377), resulting in an empty or undefined message
4. **Inconsistent behavior**: When totpScript is not configured, the code attempts to show a TOTP code that was never generated

These issues prevent users from:
- Using TOTP scripts with command-line arguments
- Having TOTP codes automatically copied to clipboard
- Seeing meaningful output when TOTP script is not configured

## Proposed Solution

Fix the TOTP script execution and clipboard integration to work correctly with array-based configuration and proper variable handling:

### 1. Array Execution
Change `${totpScript}` to `"${totpScript[@]}"` to properly expand array with all arguments

### 2. Clipboard Command
Replace `${clipboardCommand}` with `"${clipboardCopyCommand[@]}"` to match the config.env variable name and support array-based commands

### 3. Logic Flow
Fix the if/else logic to only display messages when they make sense:
- If totpScript is configured: generate code, copy to clipboard, display code
- If totpScript is not configured: prompt user that manual MFA entry is needed

### 4. Documentation
Update config.env comments and README to document the proper array syntax for totpScript

## Benefits

- **Working TOTP integration**: Users can properly configure external TOTP generators (like password managers)
- **Proper clipboard support**: TOTP codes are correctly copied to clipboard using the configured command
- **Clear feedback**: Users see appropriate messages based on their configuration
- **Better UX**: Eliminates confusing undefined variable messages
- **Example compatibility**: Works with existing config patterns like rbw-menu.sh with arguments

## Scope

This change modifies:
1. `_bmclib.sh` lines 376-383 - Fix array expansion and variable references
2. Config documentation - Add examples and clarify array syntax for totpScript and clipboard commands

## Dependencies

- Existing config.env structure (no breaking changes)
- Bash array expansion support (already required by BMC)

## Out of Scope

- Adding new TOTP generation methods
- Changing the MFA authentication flow
- Adding automatic clipboard paste functionality
- Supporting non-array configuration formats (arrays are more flexible)
