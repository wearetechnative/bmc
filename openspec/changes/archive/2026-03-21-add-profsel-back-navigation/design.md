# Design: add-profsel-back-navigation

## Overview
This change adds back navigation capability to the `selectAWSProfile` function by wrapping the group and profile selection steps in a loop that allows users to return to group selection when they cancel at the profile selection stage.

## Architecture

### Current Flow
```
selectAWSProfile()
├─ if preferedProfile set → use it directly
└─ else:
   ├─ Select group (gum filter) → awsProfileGroups
   │  └─ if empty → unset selectedProfileName, return
   ├─ Select profile from group (gum filter) → selectedProfile
   │  └─ if empty → unset selectedProfileName, return
   └─ Parse and set profile variables
```

### Proposed Flow
```
selectAWSProfile()
├─ if preferedProfile set → use it directly
└─ else:
   └─ while true:
      ├─ Select group (gum filter) → awsProfileGroups
      │  └─ if empty → unset selectedProfileName, return (EXIT LOOP)
      ├─ Select profile from group (gum filter) → selectedProfile
      │  ├─ if empty → continue (BACK TO GROUP SELECTION)
      │  └─ else → break (EXIT LOOP)
      └─ Parse and set profile variables
```

## Implementation Details

### Modified Function: `selectAWSProfile` in `_bmclib.sh`

The key changes:
1. Add `while true` loop around the selection logic (after checking `preferedProfile`)
2. Keep group selection cancellation behavior unchanged (returns immediately)
3. Change profile selection cancellation behavior:
   - Instead of `return`, use `continue` to restart the loop
   - Optionally add a message like "Returning to group selection..."
4. Add `break` statement after successful profile selection to exit the loop

### Code Structure
```bash
function selectAWSProfile {
  if [[ -z $preferedProfile ]]; then
    while true; do
      # Select group
      awsProfileGroups=$(jsonify-aws-dotfiles | jq -r '[.config[].group] | unique | sort | .[]' | grep -v null | gum filter --height 25)

      # Check if group selection was cancelled
      if [[ -z $awsProfileGroups ]]; then
        unset selectedProfileName
        return  # Exit function entirely
      fi

      # Build and display profile table
      selectedProfileTable=$(...)
      header=$(echo "  $selectedProfileTable" | head -n1)
      selectedProfile=$(echo "$selectedProfileTable" | tail -n +2 | gum filter --header="$header")

      # Check if profile selection was cancelled
      if [[ -z $selectedProfile ]]; then
        continue  # Return to group selection
      fi

      # Profile selected successfully - parse and exit loop
      # ... (existing parsing logic) ...
      break  # Exit loop with valid selection
    done
  else
    # Existing preferred profile logic
    # ... (unchanged) ...
  fi

  # Existing post-selection logic
  # ... (unchanged) ...
}
```

## Edge Cases

1. **Rapid cancellation:** User cancels group → exits immediately (existing behavior)
2. **Multiple back navigations:** User can cancel profile selection multiple times to try different groups
3. **Preferred profile flag (`-p`):** Bypasses loop entirely (existing behavior)
4. **Empty group:** If selected group has no profiles, table would be empty - already handled by existing `if [ ${#aws_profiles[@]} -eq 0 ]` check in other functions (not directly applicable here, but profile list will be empty)

## Testing Considerations

### Unit Test Scenarios (Manual)
1. Normal flow: group → profile → success
2. Back navigation: group → cancel profile → group → profile → success
3. Multiple backs: group → cancel → group → cancel → group → profile → success
4. Exit at group: cancel group → exit
5. Preferred profile: `-p profile-name` → bypass loop
6. Sourced mode: ensure `return` doesn't close shell
7. Executed mode: ensure proper exit codes

### Shell Compatibility
- bash 4+: `while true`, `continue`, `break` are all POSIX-compatible
- zsh 5.8+: Same constructs work identically

## Performance Impact
Negligible - adds one loop structure with minimal overhead. The `gum filter` calls are the primary performance bottleneck, and they remain unchanged.

## Security Impact
None - this is a navigation enhancement with no security implications.

## Backward Compatibility
Fully backward compatible:
- Existing successful selection paths unchanged
- Existing cancellation at group level unchanged
- Existing preferred profile flag behavior unchanged
- New behavior only activates on profile-level cancellation (previously would exit, now loops back)

## Future Enhancements (Out of Scope)
1. Add explicit "← Back" option in profile menu
2. Add visual indicator showing current group name in profile selection header
3. Add keyboard shortcut (e.g., ESC for back, Ctrl-C for exit)
4. Add history/breadcrumb showing navigation path
