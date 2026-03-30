# Proposal: Fix ec2ls Name Value with Spaces

## Problem
The `bmc ec2ls` command incorrectly displays only the first word of EC2 instance Name tags when the tag contains spaces. For example, an instance with Name tag "My Server Instance" only shows "My" in the output table.

## Root Cause
In `_bmclib.sh:99-105`, the `ec2ListInstances` function uses:
```bash
instance_id=$(echo $line | awk '{print $1}')
private_ip=$(echo $line | awk '{print $2}')
# ... etc
name=$(echo $line | awk '{print $5}')
```

There are two issues:
1. **Missing field separator**: The AWS CLI `describe-instances` command with `--output text` returns tab-separated values. However, `awk` with default field separator (whitespace) treats each word as a separate field, causing only the first word to be captured when a field contains spaces.
2. **Unquoted variable expansion**: `echo $line` without quotes causes bash to perform word splitting, which converts tabs to spaces before awk even processes the data.

## Proposed Solution
Apply two fixes to all field extractions in the `ec2ListInstances` function (lines 99-105):

1. Add quotes around the `$line` variable to preserve tabs during expansion
2. Set awk field separator to tab (`-F'\t'`) to correctly parse tab-delimited fields

Update lines 99-105 in `_bmclib.sh` to:
```bash
instance_id=$(echo "$line" | awk -F'\t' '{print $1}')
private_ip=$(echo "$line" | awk -F'\t' '{print $2}')
public_ip=$(echo "$line" | awk -F'\t' '{print $3}')
state=$(echo "$line" | awk -F'\t' '{print $4}')
name=$(echo "$line" | awk -F'\t' '{print $5}')
hibernation_status=$(echo "$line" | awk -F'\t' '{print $6}')
scheduler_tag=$(echo "$line" | awk -F'\t' '{print $7}')
```

This ensures that tab-delimited fields from AWS CLI are correctly parsed while preserving spaces within field values.

## Impact
- **User-facing**: EC2 instance names with spaces will display correctly in `bmc ec2ls` output
- **Compatibility**: No breaking changes, only fixes incorrect behavior
- **Related commands**: The `ec2FindInstance` function already processes the output from `ec2ListInstances`, so fixing the source will automatically fix search results

## Alternatives Considered
1. **Quote the Name field in AWS query**: Would require restructuring the entire query and output processing logic - more complex than necessary
2. **Use different field positions**: Not applicable - the issue is the field separator, not the position
3. **Switch to JSON output**: Would require complete rewrite of parsing logic - unnecessary for this fix

## Testing
- Test with instances that have Name tags containing spaces
- Test with instances that have Name tags without spaces (ensure no regression)
- Test with instances that have no Name tag (ensure "None" handling still works)
- Test the `ec2FindInstance` command to ensure search still works correctly
