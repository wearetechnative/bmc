---
title: "profsel"
weight: 10
---

`bmc profsel` interactively selects an AWS profile and exports it as `AWS_PROFILE` in your current shell.

## Usage

```bash
bmc profsel              # Interactive profile selection
bmc profsel -p myprofile # Pre-select a profile by name
bmc profsel -l           # List all profiles in tabular format
bmc profsel --json       # Output selected profile as JSON
```

## Shell integration required

`bmc profsel` must be invoked through the shell wrapper for `AWS_PROFILE` to be set in your current shell. Install it once with:

```bash
bmc install-shell-integration
```

See [Shell Integration](/setup/shell-integration/) for details.

## JSON output

`--json` outputs the selected profile as a JSON object:

```json
{
  "source_profile": "technative-long-term",
  "profile_name": "technative-prod",
  "profile_arn": "arn:aws:iam::123456789012:role/..."
}
```

Useful for scripting — pipe into `jq` or other tools.

## Recent profiles

The interactive picker shows recently used profiles at the top (last 10, labelled "recent"). History is shared across all commands that trigger interactive profile selection — a profile selected in `bmc console` or `bmc ec2connect` also appears recent in `bmc profsel`. History is stored in `~/.local/share/bmc/profile-history.json`.
