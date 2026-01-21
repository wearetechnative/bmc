# Proposal: Prompt to Start Stopped Instances in ec2connect

## Why

When users select a stopped EC2 instance in `bmc ec2connect`, the script immediately exits with an error message: "!!! Instance chosen is not running. Current state is : stopped. Not executing the SSH-command". This forces users to:

1. Exit ec2connect
2. Run `bmc ec2stopstart` separately
3. Wait for the instance to start
4. Re-run `bmc ec2connect`
5. Re-select the same instance

This multi-step workflow is cumbersome and interrupts the user's flow when they want to connect to an instance that happens to be stopped.

## What Changes

Enhance `bmc ec2connect` to detect when a selected instance is in a stopped state and offer to start it automatically. The workflow becomes:

1. User selects a stopped instance
2. Script prompts: "Instance is stopped. Start it? (y/n)"
3. If yes: Start the instance using the same logic as `ec2stopstart`, wait for it to reach running state, then proceed with connection
4. If no: Exit gracefully without starting

Additionally, provide a configuration option `BMC_AUTO_START_STOPPED_INSTANCES` in `~/.config/bmc/config.env` that allows users to:
- Skip the prompt and automatically start stopped instances (value: "always")
- Skip the prompt and never start stopped instances (value: "never")
- Always prompt (value: "prompt" or unset - default behavior)

## Impact

**Affected code:**
This change modifies the `ec2connect.sh` script to:
1. Detect stopped instance state (already done at line 53-54)
2. Prompt user to start the instance (new)
3. Read and respect config option from `~/.config/bmc/config.env` (new)
4. Call instance start logic and wait for running state (reuse existing patterns from `ec2StopStartInstance`)
5. Continue with normal connection flow after instance is running
6. Improve error message by removing "Not executing the SSH-command" text from the non-running instance error message (line 55)

**Affected specs:**
- Creates or updates `ec2connect` spec documenting the auto-start feature
