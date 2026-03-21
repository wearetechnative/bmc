# Proposal: add-profsel-back-navigation

## Problem
When using `bmc profsel`, users select an AWS profile through a two-step process:
1. Select a profile group via `gum filter`
2. Select a specific profile from that group via `gum filter`

Currently, after selecting a group in step 1, the only way to return to the group selection menu is to press Ctrl-C, which cancels the entire operation. There is no "back" or "return to previous menu" option available at the profile selection stage.

This creates a poor user experience when:
- Users accidentally select the wrong group
- Users want to explore profiles in different groups
- Users realize mid-selection they need a profile from a different group

## Proposed Solution
Add a navigation loop to the `selectAWSProfile` function that allows users to return to the group selection menu from the profile selection menu without canceling the entire operation.

The solution will:
1. Wrap the profile group and profile selection in a loop
2. Detect when the user cancels at the profile selection stage (empty selection)
3. Provide feedback that they're returning to the group selection
4. Allow the loop to exit when either a valid profile is selected or the user cancels at the group selection stage

## Scope
- **In scope:**
  - Modify `selectAWSProfile` function in `_bmclib.sh` to support back navigation
  - Maintain existing cancellation behavior (Ctrl-C at group selection exits gracefully)
  - Preserve compatibility with the `-p` (preferred profile) flag that bypasses selection

- **Out of scope:**
  - Adding explicit "Back" menu items in the gum filter UI (relying on cancel/empty behavior is sufficient)
  - Changing other commands that use profile selection
  - Modifying the console command's profile selection flow

## User Impact
- **Positive:** Users can navigate back to group selection without restarting the command
- **Neutral:** No breaking changes - existing behavior for successful selections and group-level cancellations remains unchanged
- **Risk:** Minimal - this adds a loop but maintains all existing exit conditions

## Technical Approach
The implementation will add a `while` loop around the two selection steps in `selectAWSProfile`:
- Continue looping if profile selection returns empty (user pressed Ctrl-C at profile stage)
- Exit loop if group selection returns empty (user pressed Ctrl-C at group stage)
- Exit loop if a valid profile is selected

## Alternatives Considered
1. **Add explicit "Back" option in profile menu** - Would require modifying the gum filter output and parsing it differently, adding complexity
2. **Add "Back" option in group menu** - Would cause a loop without allowing group re-selection
3. **Keep current behavior** - Forces users to restart the command, poor UX

## Dependencies
None - this is a self-contained change to the `selectAWSProfile` function.

## Testing Strategy
- Manual testing in both bash and zsh
- Test scenarios:
  - Select group, cancel at profile selection, verify return to group selection
  - Select group, select profile, verify normal completion
  - Cancel at group selection, verify graceful exit
  - Test with `-p` flag to ensure bypass still works
  - Test in both sourced and executed modes
