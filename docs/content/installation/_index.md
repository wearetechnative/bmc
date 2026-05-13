---
title: "Installation"
weight: 10
---

BMC is available via Homebrew, Nix, binary download, and a one-liner install script for Linux.

## Homebrew (macOS / Linux)

```bash
brew install wearetechnative/tap/bmc
```

## Linux — one-liner

Download and install the latest release directly:

```bash
curl -fsSL https://github.com/wearetechnative/bmc/releases/latest/download/bmc_$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/').tar.gz \
  | tar -xz bmc && sudo mv bmc /usr/local/bin/
```

Or manually: pick the right archive for your architecture from [GitHub Releases](https://github.com/wearetechnative/bmc/releases/latest):

| Architecture | File |
|---|---|
| x86\_64 (most systems) | `bmc_<version>_linux_amd64.tar.gz` |
| ARM64 (Raspberry Pi, AWS Graviton) | `bmc_<version>_linux_arm64.tar.gz` |

```bash
tar -xzf bmc_*.tar.gz bmc
sudo mv bmc /usr/local/bin/
```

## Nix

### nix profile (recommended)

```bash
nix profile add github:wearetechnative/bmc
```

### nix-env

```bash
nix-env -iA bmc -f https://github.com/wearetechnative/bmc/archive/main.tar.gz
```

### NixOS — configuration.nix

```nix
{
  inputs.bmc.url = "github:wearetechnative/bmc";
  # ...
  environment.systemPackages = [ inputs.bmc.packages.${system}.bmc ];
}
```

## macOS binary

Download from [GitHub Releases](https://github.com/wearetechnative/bmc/releases/latest):

| Architecture | File |
|---|---|
| Apple Silicon (M1/M2/M3) | `bmc_<version>_darwin_arm64.tar.gz` |
| Intel | `bmc_<version>_darwin_amd64.tar.gz` |

## Verify installation

```bash
bmc version
bmc doctor
```

`bmc doctor` checks all required and optional dependencies and prints install hints for anything missing.

## Next step

→ [Setup — Shell integration](/setup/shell-integration/)
