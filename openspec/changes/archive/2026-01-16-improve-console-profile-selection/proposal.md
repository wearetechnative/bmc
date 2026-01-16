# Change: Improve console profile selection to respect AWS_PROFILE

## Why
Currently, `bmc console` always prompts for profile selection even when `AWS_PROFILE` is already set in the environment. This creates unnecessary friction for users who have already selected a profile using `bmc profsel` or other means. Users must redundantly select the same profile again when opening the console.

## What Changes
- `bmc console` will check if `AWS_PROFILE` environment variable is set before prompting for profile selection
- If `AWS_PROFILE` is set, the command will use that profile directly without prompting
- A new `-p` flag (without value) allows users to force profile selection even when `AWS_PROFILE` is set
- The existing `-p <profile-name>` syntax continues to work for specifying a profile directly
- The `-s <service>` option continues to work unchanged

## Impact
- Affected specs: `aws-console-access` (new capability)
- Affected code: `bmc:131-160` (console function)
- User experience: Reduces friction when `AWS_PROFILE` is already configured
- Backwards compatibility: Existing `-p <profile>` and `-s <service>` usage remains unchanged
