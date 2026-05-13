---
title: "BMC"
description: "TechNative AWS/Terraform DevOps tools — a single Go binary that simplifies working with AWS"
layout: "background"
---

**BMC (Bill McCloud)** is a single Go binary by TechNative that simplifies everyday AWS operations — profile selection, EC2/ECS access, MFA, and console access — all from the terminal.

[**Get Started →**](/installation/) &nbsp;&nbsp; [GitHub](https://github.com/wearetechnative/bmc)

---

## Install in seconds

```bash
# Homebrew
brew install wearetechnative/tap/bmc

# Nix
nix profile add github:wearetechnative/bmc
```

---

## What BMC does

| | |
|---|---|
| **Profile selection** | Interactively switch AWS profiles and set `AWS_PROFILE` in your shell |
| **EC2 access** | List, connect (SSH/SSM), start/stop, find instances across profiles |
| **AWS Console** | Open the AWS console for any profile in Firefox or Chrome |
| **MFA** | Automatic session refresh with TOTP support |
| **ECS** | Interactive shell into ECS containers |

---

## Explore the docs

- [Installation](/installation/) — Homebrew, Nix, binary download
- [Setup](/setup/) — Shell integration, configuration, MFA
- [Commands](/commands/) — All bmc commands with examples
- [Advanced](/advanced/) — Chrome profiles, NixOS, migration
