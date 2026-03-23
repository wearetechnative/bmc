# Implementation Tasks

## Task List

- [x] 1. **Update ec2scheduler.sh to show Scheduler status**
   - Modify jq query at line ~23-30 to extract `InstanceScheduler` tag value
   - Add logic to determine if tag exists (any value = "yes", missing = "no")
   - Update table header to include "Scheduler" column: `["InstanceId", "Name", "State", "Scheduler", "IgnoreUntil"]`
   - Update CSV data rows to include scheduler status in correct position
   - Adjust gum table column widths (currently `-w 20,35,12,30`) to accommodate new column: `-w 20,30,12,10,30`
   - **Validation**: Run `bmc ec2scheduler` and verify Scheduler column displays correctly for instances with and without InstanceScheduler tag

- [x] 2. **Update ec2ListInstances in _bmclib.sh for Hibernate normalization**
   - Modify AWS CLI query at line ~90-92 to include logic for normalizing Hibernate values
   - Add awk/jq transformation to convert:
     - "true" → "yes"
     - "false" → "no"
     - "None" or empty → "no"
   - Update output assembly at line ~106 to use normalized hibernate value
   - **Validation**: Run `bmc ec2ls` and verify Hibernate column shows "yes" or "no" instead of "true"/"false"/"None"

- [x] 3. **Update ec2ListInstances in _bmclib.sh to show Scheduler status**
   - Modify AWS CLI query to include `InstanceScheduler` tag: `Tags[?Key==\`InstanceScheduler\`].Value | [0]`
   - Add logic to convert tag presence to "yes"/"no" (any value = "yes", empty/null = "no")
   - Update table header at line ~95: `"InstanceId,PrivateIpAddress,PublicIpAddress,State,Hibernate,Name,Scheduler\n"`
   - Update CSV row assembly at line ~106 to append scheduler status
   - **Validation**: Run `bmc ec2ls` and verify Scheduler column displays correctly

- [x] 4. **Update ec2FindInstance in _bmclib.sh for new columns**
   - Update CSV header at line ~55 to include Scheduler column
   - Update CSV row assembly at line ~73 to include scheduler status value
   - Ensure selectProfileGroup call passes scheduler info correctly
   - **Validation**: Run `bmc ec2find <search-term>` and verify output includes Scheduler column

- [x] 5. **Create ec2ls spec**
   - Create new spec directory: `openspec/changes/enhance-instance-list-scheduler-status/specs/ec2ls/`
   - Write spec.md with ADDED requirements for:
     - Listing EC2 instances with formatted output
     - Normalizing Hibernate values to yes/no
     - Displaying Scheduler status
     - Interactive table display with gum
   - Include scenarios for instances with/without hibernation and scheduler tags
   - **Validation**: Ensure spec follows BMC conventions and includes all necessary scenarios

- [x] 6. **Update ec2scheduler spec**
   - Create spec delta: `openspec/changes/enhance-instance-list-scheduler-status/specs/ec2scheduler/spec.md`
   - Add MODIFIED requirement for "List All EC2 Instances"
   - Add scenario: "Display scheduler configuration status"
   - Update scenario: "Display ignore override information" to include Scheduler column
   - **Validation**: Ensure spec delta correctly describes modifications

- [x] 7. **Update documentation**
   - Update CHANGELOG.md under "## NEXT VERSION" with enhancement notice
   - Include note about new Scheduler column in ec2scheduler and ec2ls
   - Include note about Hibernate format change (true/false → yes/no)
   - Mention potential impact on automated parsing
   - **Validation**: Ensure changelog follows project conventions

- [x] 8. **Validate and test**
   - Run `openspec validate enhance-instance-list-scheduler-status --strict --no-interactive`
   - Test ec2scheduler with instances that have InstanceScheduler tag
   - Test ec2scheduler with instances without InstanceScheduler tag
   - Test ec2ls with instances with hibernation enabled/disabled
   - Test ec2ls with instances with/without InstanceScheduler tag
   - Verify table formatting looks good with various instance name lengths
   - **Validation**: All commands display correctly formatted tables with new columns

## Dependencies

- Tasks 1, 2-4 are independent and can be worked on in parallel
- Task 5 and 6 (spec creation) can be done in parallel with implementation
- Task 7 should be done after implementation is complete
- Task 8 must be done last after all other tasks complete

## Testing Notes

- Test with instances in different states: running, stopped, stopping, pending
- Test with instances that have InstanceScheduler tag with various values
- Test with instances without InstanceScheduler tag
- Test with instances with hibernation enabled (true)
- Test with instances without hibernation (false or None)
- Test with instances that have both InstanceScheduler and Ignore_scheduler tags
- Verify CSV output can still be parsed if users have automation
- Check table formatting with very long instance names
- Verify column widths provide good readability on standard terminal sizes
