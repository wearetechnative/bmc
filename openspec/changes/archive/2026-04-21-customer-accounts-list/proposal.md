## Why

There is no quick way to see which AWS accounts (profiles) belong to a specific customer. Users currently have to mentally filter through all profiles or use `profsel` and browse groups interactively. A dedicated subcommand that lets you pick a customer (profile group) and immediately see all their accounts would speed up daily operations.

## What Changes

- Add a new `accountls` subcommand to `bmc`
- The subcommand presents an interactive customer (profile group) selector using `gum filter`
- After selection, it prints a table of all AWS profiles belonging to that customer with their account details (profile name, ARN/account ID)
- The command exits after displaying the listing (no further action)

## Capabilities

### New Capabilities
- `account-listing`: Interactive customer selection followed by a formatted table of AWS accounts/profiles for that customer

### Modified Capabilities
<!-- No existing spec-level requirements are changing -->

## Impact

- `bmc` main dispatcher: new subcommand registration
- `_bmclib.sh`: may add a new function or reuse `printAWSProfiles()` with group filtering
- Dependencies: existing tools only (`gum`, `jq`, `jsonify-aws-dotfiles`)
