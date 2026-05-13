---
title: "Installation"
weight: 10
---

BMC is available via Nix, as a `.deb`/`.rpm` package, via Homebrew, or as a plain binary.

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

## Linux

### Debian / Ubuntu (.deb)

Download the `.deb` for your architecture from [GitHub Releases](https://github.com/wearetechnative/bmc/releases/latest) and install:

```bash
sudo dpkg -i bmc_<version>_linux_amd64.deb
```

| Architecture | File |
|---|---|
| x86\_64 | `bmc_<version>_linux_amd64.deb` |
| ARM64 | `bmc_<version>_linux_arm64.deb` |

### Red Hat / Fedora / SUSE (.rpm)

```bash
sudo rpm -i bmc_<version>_linux_amd64.rpm
```

| Architecture | File |
|---|---|
| x86\_64 | `bmc_<version>_linux_amd64.rpm` |
| ARM64 | `bmc_<version>_linux_arm64.rpm` |

### One-liner (any distro)

```bash
curl -fsSL https://github.com/wearetechnative/bmc/releases/latest/download/bmc_$(curl -s https://api.github.com/repos/wearetechnative/bmc/releases/latest | grep tag_name | cut -d'"' -f4 | tr -d v)_linux_$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/').tar.gz \
  | tar -xz bmc && sudo mv bmc /usr/local/bin/
```

## macOS / Homebrew

```bash
brew install wearetechnative/tap/bmc
```

Or download the binary directly from [GitHub Releases](https://github.com/wearetechnative/bmc/releases/latest):

| Architecture | File |
|---|---|
| Apple Silicon (M1/M2/M3) | `bmc_<version>_darwin_arm64.tar.gz` |
| Intel | `bmc_<version>_darwin_amd64.tar.gz` |

## Manual binary install

Download the archive for your platform from [GitHub Releases](https://github.com/wearetechnative/bmc/releases/latest), extract, and place `bmc` on your `$PATH`:

```bash
tar -xzf bmc_*.tar.gz bmc
sudo mv bmc /usr/local/bin/
```

## Verify installation

```bash
bmc version
bmc doctor
```

`bmc doctor` checks all required and optional dependencies and prints install hints for anything missing.

## Next step

→ [Setup — Shell integration](/setup/shell-integration/)
