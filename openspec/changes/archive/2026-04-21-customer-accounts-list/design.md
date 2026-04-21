## Context

BMC provides several subcommands for AWS operations (`profsel`, `ec2ls`, `console`, etc.). Profile selection is handled by `selectAWSProfile()` in `_bmclib.sh`, which uses a two-stage flow: select profile group (customer) → select profile within group. The profile data comes from `jsonify-aws-dotfiles` which parses `~/.aws/config` and `~/.aws/credentials` into JSON.

Currently there is no way to just list all accounts for a given customer without entering the full profile selection flow.

## Goals / Non-Goals

**Goals:**
- Provide a new `accountls` subcommand that shows all profiles for a selected customer
- Reuse existing profile group selection and JSON data infrastructure
- Display output as a formatted table using `gum table`

**Non-Goals:**
- Setting `AWS_PROFILE` or performing any AWS API calls
- Adding filtering/search within the account listing
- Exporting account data to files

## Decisions

### 1. Subcommand name: `accountls`
Follows existing naming convention (`ec2ls`, `ec2connect`). Alternatives considered: `customers`, `accounts` — rejected because `accountls` is consistent with the `*ls` pattern for listing commands.

### 2. Implement as inline function in `_bmclib.sh`
The logic is straightforward (group select + filtered table display), similar to `printAWSProfiles()`. No need for a separate `.sh` file. The dispatcher in `bmc` will call the function directly.

### 3. Reuse profile group selection from `selectAWSProfile()`
Extract only the group selection step (stage 1) from the existing flow. Use the same `gum filter` pattern and `jsonify-aws-dotfiles` JSON source. This ensures consistent UX across commands.

### 4. Table output with `gum table`
Display columns: Profile Name, Account ID (extracted from ARN), Role. Consistent with how `ec2ls` displays results. The command exits immediately after displaying the table.

## Risks / Trade-offs

- [Duplicating group selection logic] → Keep it minimal; extract the group selection into its own small helper or inline the few lines needed. The existing `selectAWSProfile()` function combines group + profile selection tightly, so a small amount of duplication is acceptable for clarity.
