# Proposal: add-profsel-json-output

## Summary
Add a `--json` flag to `bmc profsel` command that outputs selected profile information as JSON instead of the normal interactive output. This enables programmatic integration and scripting use cases.

## Context
Currently, `bmc profsel` is designed for interactive terminal use with gum menus and human-readable output. Users who want to integrate profile selection into scripts or automation workflows need a machine-readable output format.

The command currently supports:
- Interactive profile selection (default)
- `-l` flag to list profiles
- `-p <profile>` flag to pre-select a profile

## Motivation
- Enable scripting and automation that needs to capture selected profile details
- Provide machine-readable output for integration with other tools
- Maintain consistency with modern CLI conventions (many tools offer `--json` flags)

## Scope
This change adds JSON output capability to the `bmc profsel` command:
- Add `--json` flag to profsel command
- Keep interactive profile selection (gum menus) when `--json` is used
- Output JSON to file descriptor 3 when available (allows progress to remain visible)
- Fallback to stdout when fd 3 is not available (backward compatible)
- Redirect progress/status messages intelligently based on JSON output destination
- Output final result as JSON object containing: source_profile, profile_name, profile_arn (role ARN)
- Maintain backward compatibility with existing flags and behavior
- Support both interactive mode (`--json` alone) and non-interactive mode (`--json -p <profile>`)

## Out of Scope
- JSON output for other bmc commands (can be added separately if needed)
- Changing the structure of existing profile selection logic
- Adding additional profile metadata beyond what's already available
- Changing how AWS_PROFILE environment variable is set (still not set when using --json)

## Affected Components
- `bmc` script (profsel command argument parsing)
- `_bmclib.sh` (selectAWSProfile function to support quiet mode)
- profile-selection spec (new requirement for JSON output)

## Alternatives Considered
1. **Use existing `-p` flag with parsing**: Could use `-p` to select and parse text output, but requires fragile text parsing
2. **New separate command**: Could create `bmc profsel-json`, but adds unnecessary command proliferation
3. **Always output JSON when piped**: Could detect pipe and auto-switch to JSON, but less explicit and harder to debug
4. **Suppress all interactive elements**: Could make --json non-interactive only, but reduces usability for scripting scenarios where user interaction is acceptable

## Dependencies
None. This is an additive feature that doesn't depend on other changes.

## Risks
- Low risk: Feature is purely additive and isolated to profsel command
- Need to redirect progress messages to stderr to keep stdout clean for JSON
- Need to handle error cases in JSON output (e.g., cancelled selection, no profile found)
- MFA messages and other status updates should go to stderr when --json is active
