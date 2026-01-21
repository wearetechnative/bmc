# Proposal: Fix TOTP Script Handling and Clipboard Integration

## Why

The current TOTP script functionality in BMC has several implementation issues that prevent it from working correctly:

1. **Array handling broken**: `totpScript` is defined as an array in `config.env` (e.g., `totpScript=("/path/to/script" "-t" "code")`) but executed as `${totpScript}` in `_bmclib.sh:377`, which doesn't properly expand arrays with arguments
2. **Undefined variable**: Line 378 references `${clipboardCommand}` which is not defined in config.env; the config file defines `clipboardCopyCommand` and `clipboardPasteCommand` instead
3. **Logic error**: Line 382 displays `"Code: ${totpCode}"` in the else block (when totpScript is not set), but `totpCode` is only defined inside the if block (line 377), resulting in an empty or undefined message
4. **Inconsistent behavior**: When totpScript is not configured, the code attempts to show a TOTP code that was never generated

These issues prevent users from:
- Using TOTP scripts with command-line arguments
- Having TOTP codes automatically copied to clipboard
- Seeing meaningful output when TOTP script is not configured

## What Changes

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

## Impact

**Affected code:**
1. `_bmclib.sh` lines 376-383 - Fix array expansion and variable references
2. Config documentation - Add examples and clarify array syntax for totpScript and clipboard commands

**Affected specs:**
- Creates new `totp-integration` spec documenting TOTP script and clipboard functionality
