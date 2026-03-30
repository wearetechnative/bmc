# Proposal: Enhance Instance List with Scheduler Status

## Why

Currently, the `bmc ec2scheduler` and `bmc ec2ls` commands display EC2 instance information but lack consistent visibility into scheduler configuration:

1. **ec2scheduler**: Shows instances with their `Ignore_scheduler` tag status but doesn't indicate whether the base `InstanceScheduler` tag is configured. Users cannot quickly see which instances have scheduler enabled versus those that don't.

2. **ec2ls**:
   - Displays hibernation status as raw AWS values (true/false/None) which is inconsistent with BMC's user-friendly approach
   - Provides no visibility into scheduler configuration, making it difficult to understand which instances are managed by the scheduler

These limitations make it harder for users to:
- Quickly identify which instances are managed by the scheduler
- Understand the current scheduler configuration state at a glance
- Make informed decisions about which instances need scheduler tags added

## What Changes

### 1. ec2scheduler Command Enhancement

Add a "Scheduler" column to the instance table that shows whether the `InstanceScheduler` tag is configured:
- Display "yes" if the `InstanceScheduler` tag exists (regardless of value)
- Display "no" if the `InstanceScheduler` tag is missing
- This provides immediate visibility into base scheduler configuration

### 2. ec2ls Command Enhancement

Improve the instance listing display for better usability:

**Hibernate Column Normalization**:
- Convert "true" → "yes"
- Convert "false" → "no"
- Convert "None" (or missing) → "no"
- This aligns with BMC's user-friendly output style

**Add Scheduler Column**:
- Display "yes" if the `InstanceScheduler` tag exists
- Display "no" if the `InstanceScheduler` tag is missing
- Position as the last column in the table

## Impact

### Affected Specs

- **ec2scheduler** spec - MODIFIED requirements
  - Update "List All EC2 Instances" requirement to include Scheduler column
  - Add scenario for displaying scheduler configuration status

- **ec2ls** spec - NEW spec (currently no spec exists for ec2ls)
  - Create spec with requirements for instance listing
  - Include requirements for Hibernate display format
  - Include requirements for Scheduler status display

### Affected Code

- **ec2scheduler.sh**:
  - Modify jq query to extract both `InstanceScheduler` and `Ignore_scheduler` tags
  - Add "Scheduler" column to table header and data rows
  - Update table width configuration to accommodate new column
  - Affected lines: ~22-35 (table building and display)

- **_bmclib.sh** (`ec2ListInstances` function):
  - Modify AWS CLI query to include `InstanceScheduler` tag
  - Add transformation logic to convert Hibernate true/false/None to yes/no
  - Add Scheduler column showing yes/no based on `InstanceScheduler` tag presence
  - Update table header and column generation
  - Affected lines: ~88-111 (ec2ListInstances function)

### Breaking Changes

None. This is purely additive and formatting improvement:
- New columns added (Scheduler)
- Display format improved (Hibernate values normalized)
- No changes to command syntax or behavior
- Existing automation parsing CSV output may need adjustment for new column order

## Design Decisions

1. **Consistent "yes/no" format**: Use "yes/no" instead of "true/false" for better UX consistency across BMC commands

2. **Scheduler column placement**:
   - ec2scheduler: Place after State column, before IgnoreUntil
   - ec2ls: Place as last column to minimize disruption to existing layouts

3. **Tag presence check**: Only check for tag existence, not value validation
   - Any value in `InstanceScheduler` tag counts as "yes"
   - This matches current BMC behavior of treating any value as enabled

4. **No filtering**: Keep showing all instances
   - Don't filter based on scheduler status
   - Users should see complete picture of their infrastructure

5. **Column width allocation**:
   - Scheduler column: ~10 characters (sufficient for "yes"/"no")
   - Adjust other column widths as needed to maintain readability
