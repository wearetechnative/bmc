# Proposal: Add Ignore_scheduler Tag Support

## Why

The current `bmc ec2scheduler` command allows toggling the scheduler on/off by renaming the `InstanceScheduler` tag to `InstanceScheduler_DISABLED` (and vice versa). However, this completely disables scheduling for the instance, which is not ideal when you only need to temporarily extend an instance's runtime beyond a scheduled stop time.

The upstream terraform-aws-module-scheduler (https://github.com/wearetechnative/terraform-aws-module-scheduler) supports a more granular feature: the `Ignore_scheduler` tag. This tag allows users to keep an instance running until a specific time (e.g., "22:00 Europe/Amsterdam"), after which the tag is automatically removed and the instance resumes its normal schedule.

This is more user-friendly for common scenarios like:
- Running a long task that extends beyond scheduled stop time
- Debugging issues that require the instance to stay running late
- Temporary overrides without completely disabling the scheduler

Currently, BMC users must manually add/edit the `Ignore_scheduler` tag through the AWS Console or raw AWS CLI commands, which is cumbersome.

## What Changes

Modify the `bmc ec2scheduler` command to support managing the `Ignore_scheduler` tag in addition to (or instead of) the current toggle functionality:

1. **Display `Ignore_scheduler` status** in the instance table
   - Show if an `Ignore_scheduler` tag exists and its value (until when the instance will stay running)
   - Update the table columns to include this information

2. **Add/modify `Ignore_scheduler` tag** for selected instances
   - Allow users to set a time until which the instance should ignore scheduled stops
   - Prompt for time in format: `HH:MM Timezone` (e.g., "22:00 Europe/Amsterdam")
   - Create or update the `Ignore_scheduler` tag with the specified value

3. **Remove `Ignore_scheduler` tag** when no longer needed
   - Allow users to manually remove the override and return to normal schedule immediately

4. **Replace existing toggle functionality** with the `Ignore_scheduler` approach
   - Remove the InstanceScheduler/InstanceScheduler_DISABLED toggle mechanism
   - Implement menu-based interaction: "Set ignore until time", "Remove ignore override", "Cancel"
   - Use free-form text entry for timezone (with helpful examples/suggestions)

## Impact

- **Affected specs**:
  - `ec2scheduler` spec - MODIFIED requirements for managing scheduler overrides
    - Remove requirements related to toggling between InstanceScheduler/InstanceScheduler_DISABLED
    - Add requirements for managing Ignore_scheduler tags
    - Update display requirements to show Ignore_scheduler status

- **Affected code**:
  - `ec2scheduler.sh` - Major rewrite:
    - Replace tag renaming logic with Ignore_scheduler tag add/remove
    - Update table display to show Ignore_scheduler status and expiry time
    - Add menu for selecting action (set/remove ignore override)
    - Add prompts for time (HH:MM) and timezone input with validation
    - Update confirmation and feedback messages

- **Breaking changes**:
  - **BREAKING**: Users can no longer toggle scheduler enabled/disabled via tag renaming
  - The new approach requires users to set a specific time for override instead
  - Existing behavior for instances without scheduler tags remains unchanged (still shows instructions to add InstanceScheduler tag)

## Design Decisions

Based on user input, the following decisions have been made:

1. **Replace toggle entirely**: The Ignore_scheduler approach is more aligned with upstream best practices and provides better UX for temporary overrides

2. **Menu-based interaction**: When user selects an instance, show options:
   - "Set ignore until time" - prompts for time and timezone
   - "Remove ignore override" - removes existing Ignore_scheduler tag
   - "Cancel" - exit without changes

3. **Free-form timezone entry**: Allow users to type timezone (e.g., "Europe/Amsterdam") with helpful examples shown in the prompt

4. **Basic validation**:
   - Validate time format (HH:MM)
   - Show examples and common timezones as guidance
   - Note: Detailed timezone validation against TZ database is optional (can be added if needed)
