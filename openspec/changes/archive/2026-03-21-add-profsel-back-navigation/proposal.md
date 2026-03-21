# Change: Add back navigation in profile selection

## Why
When using `bmc profsel`, users select an AWS profile through a two-step process: first selecting a profile group, then selecting a specific profile from that group. Currently, after selecting a group, the only way to return to the group selection menu is to press Ctrl-C, which cancels the entire operation and forces the user to restart the command.

This creates a poor user experience when:
- Users accidentally select the wrong group
- Users want to explore profiles in different groups
- Users realize mid-selection they need a profile from a different group

## What Changes
- Add navigation loop to `selectAWSProfile` function in `_bmclib.sh`
- Allow users to cancel at profile selection stage to return to group selection
- Maintain existing cancellation behavior at group selection level (exits completely)
- Preserve compatibility with `-p` (preferred profile) flag that bypasses selection

## Impact
- Affected specs: profile-selection (new back navigation requirements)
- Affected code: `_bmclib.sh` (selectAWSProfile function at lines 306-338)
- User experience: Users can now navigate back from profile selection to group selection by canceling, improving UX when exploring different profile groups
- Backward compatibility: Fully compatible - all existing behavior preserved
