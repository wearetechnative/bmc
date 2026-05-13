---
title: "ecsconnect"
weight: 40
---

`bmc ecsconnect` opens an interactive shell into a running ECS container.

## Usage

```bash
bmc ecsconnect     # Interactive: pick cluster → service → task → container
```

The command walks you through selecting:
1. ECS cluster
2. Service
3. Task
4. Container

Then launches an interactive shell via AWS SSM.

## Prerequisites

- `aws` CLI v2
- `session-manager-plugin`

Check with:

```bash
bmc doctor
```

## Install prerequisites

**aws CLI v2**

```bash
# Homebrew
brew install awscli

# Nix
nix profile add nixpkgs#awscli2
```

**session-manager-plugin**

```bash
# Homebrew
brew install --cask session-manager-plugin

# Manual
# https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html
```
