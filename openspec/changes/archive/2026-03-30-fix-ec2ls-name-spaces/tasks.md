# Tasks

## Implementation Tasks

- [x] **Update field separator and variable quoting in ec2ListInstances function**
  - Open `_bmclib.sh`
  - Locate the `ec2ListInstances` function (lines 89-130)
  - Update all `awk` commands to use `-F'\t'` flag for tab-separated field parsing
  - Add quotes around `$line` variable (use `"$line"` instead of `$line`) to preserve tabs
  - Specifically update lines 99-105 to use `echo "$line" | awk -F'\t'` instead of `echo $line | awk`
  - Verify: All seven field extractions (instance_id, private_ip, public_ip, state, name, hibernation_status, scheduler_tag) use both quoted variables and tab field separator

- [x] **Validate no regressions**
  - Review the function logic to ensure no other parsing depends on whitespace splitting
  - Confirm that the CSV output construction (line 125) remains unchanged
  - Verify that the `ec2FindInstance` function doesn't require changes (it processes the table output, not raw AWS data)

- [x] **Test the fix**
  - Run `bmc ec2ls` in an environment with:
    - EC2 instances with Name tags containing spaces
    - EC2 instances with Name tags without spaces
    - EC2 instances without Name tags
  - Verify all names display correctly in the table output
  - Test `bmc ec2find` to ensure search functionality still works with the corrected data
  - Result: Tested and confirmed working correctly

- [x] **Update CHANGELOG**
  - Add entry under "## NEXT VERSION" section
  - Category: Fix
  - Description: "Fix ec2ls Name column to display complete names when Name tag contains spaces"
