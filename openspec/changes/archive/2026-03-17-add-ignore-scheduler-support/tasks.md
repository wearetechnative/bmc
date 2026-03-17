# Implementation Tasks

## Task List

- [x] 1. **Update instance listing to show Ignore_scheduler status**
   - Modify jq query in ec2scheduler.sh to extract Ignore_scheduler tag
   - Add "IgnoreUntil" column to table display showing the Ignore_scheduler value (or "N/A")
   - Update table header and column widths to accommodate new field
   - Remove SchedulerStatus column (enabled/disabled/none) as it's no longer relevant
   - **Validation**: Run `bmc ec2scheduler` and verify table shows Ignore_scheduler values correctly

- [x] 2. **Remove old toggle logic**
   - Delete code that determines new_tag_key (InstanceScheduler vs InstanceScheduler_DISABLED)
   - Remove tag rename operations (delete old tag + create new tag with same value)
   - Remove related confirmation messages about enabling/disabling scheduler
   - **Validation**: Ensure no references to InstanceScheduler_DISABLED remain in the code

- [x] 3. **Implement action menu**
   - Add gum menu after instance selection with options: "Set ignore until time", "Remove ignore override", "Cancel"
   - Handle each menu choice with appropriate logic flow
   - Show current Ignore_scheduler value (if exists) before displaying menu
   - **Validation**: Test menu appears correctly and handles all three options

- [x] 4. **Implement "Set ignore until time" flow**
   - Prompt user for time in HH:MM format using gum input
   - Show example: "Example: 22:00"
   - Validate time format matches HH:MM pattern
   - Prompt user for timezone using gum input
   - Show examples: "Examples: Europe/Amsterdam, America/New_York, UTC"
   - Combine time and timezone into tag value: "HH:MM Timezone"
   - Create or update Ignore_scheduler tag with the combined value
   - Display success message showing instance ID and the ignore-until value
   - **Validation**: Test creating and updating Ignore_scheduler tags with various times/timezones

- [x] 5. **Implement "Remove ignore override" flow**
   - Check if Ignore_scheduler tag exists on the instance
   - If exists, delete the Ignore_scheduler tag
   - If doesn't exist, show message "No ignore override is currently set"
   - Display success message confirming removal
   - **Validation**: Test removing existing tags and attempting to remove when none exists

- [x] 6. **Update handling of instances without InstanceScheduler tag**
   - Keep existing logic that detects instances without scheduler tags
   - Maintain AWS Console opening functionality for adding InstanceScheduler tag
   - Ensure this flow still works correctly with new code structure
   - **Validation**: Test with instance that has no InstanceScheduler tag

- [x] 7. **Update user feedback messages**
   - Replace all toggle-related messages with ignore override messages
   - Update confirmation prompts to reflect new functionality
   - Ensure error messages are clear and actionable
   - **Validation**: Review all message outputs for clarity and correctness

- [x] 8. **Update spec and documentation**
   - Archive old spec requirements related to toggle functionality
   - Add new requirements for Ignore_scheduler management to ec2scheduler spec
   - Update CHANGELOG.md under "## NEXT VERSION" with breaking change notice
   - **Validation**: Run `openspec validate add-ignore-scheduler-support --strict`

## Dependencies

- Tasks 2-7 can be worked on in parallel after task 1 completes
- Task 8 should be done last after implementation is complete

## Testing Notes

- Test with instances in different states: running, stopped
- Test with instances that have InstanceScheduler, InstanceScheduler_DISABLED, no scheduler tags
- Test with instances that already have Ignore_scheduler tags
- Verify timezone strings are passed correctly to AWS tags
- Verify table formatting with long timezone names
