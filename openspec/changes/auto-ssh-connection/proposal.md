# Change: Auto-select SSH connection when SSH-specific flags are provided

## Why
When users provide SSH-specific flags (`-u` for username or `-i` for identity file) to `bmc ec2connect`, they have already indicated their intention to use SSH. Prompting them to choose between SSH and SSM adds unnecessary friction and extra steps to their workflow.

## What Changes
- When `-u` (username) flag is provided, automatically select SSH connection method
- When `-i` (identity file) flag is provided, automatically select SSH connection method
- When both flags are provided, automatically select SSH connection method
- Connection type prompt (SSM/SSH) will only appear when neither flag is provided
- Existing SSM functionality remains unchanged when no SSH flags are present

## Impact
- Affected specs: ec2connect (new capability spec)
- Affected code: ec2connect.sh:59 (connection method selection logic)
- Backward compatible: Users can still manually select SSM by not providing -u or -i flags
- Improves user experience by reducing unnecessary prompts
