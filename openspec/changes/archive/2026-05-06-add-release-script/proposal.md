## Why

bmc has no release process. Tagging, changelog finalization, Nix vendorHash updates, and binary publishing are all manual and error-prone. The project needs a repeatable release script that mirrors the pattern used in parsh, adapted for bmc's specifics (VERSION-bmc, package.nix, Homebrew tap).

This task is tracked in [bmc-5yzf](../../../.beans/bmc-5yzf--release-script.md).

## What Changes

- **`release.sh`**: Interactive release script — bumps VERSION-bmc (semver, dropping legacy 4th component), finalizes CHANGELOG, auto-updates `package.nix` vendorHash when Go deps change, verifies local build, commits, tags, and optionally pushes
- **`main.go`**: Option B version resolution — prefer ldflags-injected version over embedded file so GoReleaser and Nix builds stamp the binary correctly
- **`package.nix`**: Add version ldflags injection (`-X cmd.Version=${version}`); fix `vendorHash` from `null` placeholder to a real sha256 hash
- **`.goreleaser.yml`**: Add ldflags version injection (`-X github.com/wearetechnative/bmc/cmd.Version={{.Version}}`) as belt-and-suspenders alongside embed
- **`.github/workflows/release.yml`**: GitHub Actions workflow that triggers GoReleaser on `v*` tag push, publishing multi-platform binaries and updating the Homebrew tap
- **`docs/homebrew-tap-setup.md`**: One-time setup guide for creating the `wearetechnative/homebrew-tap` GitHub repo and configuring the `HOMEBREW_TAP_TOKEN` secret

## Capabilities

### New Capabilities

- `release-automation`: Repeatable release process covering version bump, changelog finalization, Nix vendorHash auto-update, local build verification, git tag creation, GitHub Actions + GoReleaser pipeline
- `homebrew-tap`: Setup documentation and goreleaser configuration for publishing bmc to a Homebrew tap at `wearetechnative/homebrew-tap`

### Modified Capabilities

## Impact

- `main.go`: version resolution logic change (embed → prefer ldflags if set)
- `package.nix`: ldflags addition, vendorHash fix
- `.goreleaser.yml`: ldflags addition
- New files: `release.sh`, `.github/workflows/release.yml`, `docs/homebrew-tap-setup.md`
- Requires: `HOMEBREW_TAP_TOKEN` secret configured in the GitHub repo settings (one-time)
