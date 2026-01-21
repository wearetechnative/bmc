# Proposal: Add EC2 Instance Scheduler Tag Toggle

## Why

Many AWS environments use the `InstanceScheduler` tag to automatically stop EC2 instances on a defined schedule (e.g., stop development instances overnight to save costs). However, there are times when engineers need to temporarily disable this automatic scheduling - for example, when running long-running tasks, debugging issues, or performing maintenance that extends beyond normal working hours.

Currently, users must navigate to the AWS Console or use raw AWS CLI commands to manually rename the tag from `InstanceScheduler` to `InstanceScheduler_DISABLED` and back. This is cumbersome and error-prone.

## What Changes

Add a new `bmc ec2scheduler` command that provides convenient management of the `InstanceScheduler` tag on EC2 instances:

- **List all instances** showing which have scheduling enabled, disabled, or not configured
- **Toggle tag** between `InstanceScheduler` and `InstanceScheduler_DISABLED` for instances with existing scheduler tags
- **Guide users** to add scheduler tags manually via AWS Console for instances without them
- **Console integration** - offer to open AWS Console directly to the selected instance details page using `assumego` for immediate tag addition
- **Interactive selection** using `gum` for user-friendly instance selection
- **Tag value preservation** - when toggling, the original tag value is preserved (only the tag name changes)

The command will:
1. List all EC2 instances (not just those with scheduler tags)
2. Show scheduler status (enabled/disabled/none) and schedule value for each instance in a table
3. For instances with scheduler tags: allow toggling between enabled/disabled states
4. For instances without scheduler tags: provide instructions to add tags manually via AWS Console
5. Offer to open AWS Console directly to the instance details page using `assumego` with proper environment variables and console destination URL
6. Display current tag details and request confirmation before making changes
7. Only work with the specific tag names `InstanceScheduler` and `InstanceScheduler_DISABLED`
8. Provide clear feedback about what changed

## Impact

- **Affected specs**: New `ec2scheduler` capability spec will be created
- **Affected code**:
  - `bmc` - Add new command registration for `ec2scheduler`
  - New script: `ec2scheduler.sh` - Main implementation
  - `_bmclib.sh` - May need helper functions for tag operations (if not already present)

This is a new capability and does not modify existing functionality.
