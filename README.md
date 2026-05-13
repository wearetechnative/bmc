# BMC (Bill McCloud) — TechNative AWS/Terraform DevOps tools

A single Go binary that simplifies working with AWS — profile selection, EC2/ECS operations, MFA, and console access.

**[→ Full documentation at bmc.technative.cloud](https://bmc.technative.cloud)**

## Quick install

```bash
# Homebrew
brew install wearetechnative/tap/bmc

# Nix
nix profile add github:wearetechnative/bmc
```

See [installation docs](https://bmc.technative.cloud/installation/) for Nix flake, NixOS, and binary download options.

## What it does

- **`bmc profsel`** — Interactively select an AWS profile and export it to your shell
- **`bmc ec2`** — List, connect (SSH/SSM), start/stop EC2 instances
- **`bmc ec2find`** — Find instances across all profiles in an account group
- **`bmc console`** — Open the AWS console for any profile in Firefox or Chrome
- **`bmc ecsconnect`** — Interactive shell into an ECS container

## After installing

```bash
bmc install-shell-integration   # Required for profsel
bmc doctor                      # Check setup and dependencies
```

Full setup guide: [bmc.technative.cloud/setup](https://bmc.technative.cloud/setup/)
