---
title: "Installation"
weight: 10
---

BMC is available via Homebrew, Nix, and as a pre-built binary.

## Homebrew (macOS / Linux)

```bash
brew install wearetechnative/tap/bmc
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

## Binary download

Download from [GitHub Releases](https://github.com/wearetechnative/bmc/releases) for your platform:

| Platform | Architecture |
|---|---|
| Linux | `amd64`, `arm64` |
| macOS | `amd64` (Intel), `arm64` (Apple Silicon) |

Download the archive, extract, and place `bmc` somewhere on your `$PATH`.

## Verify installation

```bash
bmc version
bmc doctor
```

`bmc doctor` checks all required and optional dependencies and prints install hints for anything missing.

## Next step

→ [Setup — Shell integration](/setup/shell-integration/)
