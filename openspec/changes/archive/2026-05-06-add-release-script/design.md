## Context

bmc is a Go binary distributed via Homebrew tap and Nix flake. It has no release automation today. The version is embedded at compile time via `//go:embed VERSION-bmc`. The file `package.nix` reads the same file via `lib.fileContents` for the Nix build. `.goreleaser.yml` already has the Homebrew tap configured but lacks version ldflags injection. There is no `release.sh`, no GitHub Actions workflow, and the Homebrew tap repo does not yet exist.

The reference implementation is `parsh/release.sh`, adapted here for bmc's specifics:
- Version file is `VERSION-bmc` (not `VERSION`)
- Nix config lives in `package.nix` (not `flake.nix` directly)
- Version was `0.2.12.0` (4-part legacy) — migrating to 3-part semver

## Goals / Non-Goals

**Goals:**
- Repeatable, interactive release script that handles all pre-publish steps
- Auto-detect and update `package.nix` vendorHash when Go deps change
- Inject version into binaries consistently across all build paths (local, GoReleaser, Nix)
- GitHub Actions workflow that triggers GoReleaser on tag push
- Homebrew tap setup documentation

**Non-Goals:**
- Fully automated (unattended) releases — the script remains interactive
- Windows binary distribution (linux/darwin amd64/arm64 only, matching existing goreleaser config)
- Automated tap repo creation — the setup doc covers this as a one-time manual step

## Decisions

### Decision: Semver migration — drop 4th version component

`VERSION-bmc` currently contains `0.2.12.0`. GoReleaser expects `vX.Y.Z` tags. The 4th component is a legacy artifact from the bash era.

**Decision**: `release.sh` parses the existing version and drops the 4th component before presenting bump options. The first release tag will be `v0.2.12` (or next patch/minor/major from there).

**Alternative considered**: Keep 4-part format. Rejected — GoReleaser and semver tooling do not support it, and it adds complexity with no benefit.

### Decision: Option B version resolution in main.go

bmc currently uses `//go:embed VERSION-bmc` in `main.go` which unconditionally sets `cmd.Version` at startup, overwriting any ldflags-injected value. GoReleaser and Nix both inject version via ldflags, but the embed overrides them.

**Decision**: Change `main.go` to check if `cmd.Version` was already set by ldflags (i.e., not `"dev"`) before overwriting with the embedded value. This makes the ldflags injection from GoReleaser and Nix effective without removing the embed fallback for local builds.

```
if cmd.Version == "dev" {
    cmd.Version = strings.TrimSpace(versionRaw)
}
```

**Alternative considered**: Remove embed entirely, use ldflags only. Rejected — local `go build` without ldflags would show `"dev"`, which is confusing. The embed ensures the correct version is always present for local builds.

**Alternative considered**: Keep embed only, update VERSION-bmc before tagging. Partially valid — the embed would be correct for GoReleaser builds since the file is updated before the tag. But Nix builds can be done at any time from any commit, so they need ldflags to be reliable.

### Decision: vendorHash lives in package.nix, not flake.nix

The parsh release script targets `flake.nix` for vendorHash. In bmc, `flake.nix` is a thin wrapper that calls `pkgs.callPackage ./package.nix {}`. The vendorHash is in `package.nix`.

**Decision**: `release.sh` applies all vendorHash sed operations to `package.nix`.

### Decision: vendorHash auto-calculation strategy

The parsh approach: temporarily zero the vendorHash, run `nix build`, parse `got: sha256-...` from the error output, restore. This works because Nix prints the expected hash when it gets a mismatch.

**Decision**: Adopt the same strategy verbatim, targeting `package.nix`. The current `vendorHash = null` is a broken placeholder (implies vendor dir, which doesn't exist). This must be set to a real hash. The release script will calculate it on first run if deps changed since the last tag.

### Decision: release.sh targets package.nix version ldflags

`package.nix` currently has `ldflags = [ "-s" "-w" ]`. To make Nix-built binaries report the correct version, add `-X github.com/wearetechnative/bmc/cmd.Version=${version}`. This is consistent with Option B — the ldflags value wins over the embed when the binary is built by Nix.

## Risks / Trade-offs

- **vendorHash = null is currently broken**: Nix builds fail today. The release script will fix this on first run, but the state between now and first release is broken. Mitigation: fix vendorHash as part of implementing this change.
- **First release drops 4th version component**: `0.2.12.0` → `0.2.12`. Homebrew/Nix users upgrading will see a version number change. Mitigation: document in CHANGELOG as a versioning normalization, not a version downgrade.
- **HOMEBREW_TAP_TOKEN must be configured before first release**: If the secret is missing, GoReleaser will fail to push to the tap. Mitigation: the setup doc covers this, and the script instructs the user to push after verifying secrets.

## Migration Plan

1. Fix `package.nix` vendorHash (calculate correct hash)
2. Update `main.go` with Option B logic
3. Update `package.nix` and `.goreleaser.yml` with ldflags
4. Write `release.sh` and make executable
5. Create `.github/workflows/release.yml`
6. Create `docs/homebrew-tap-setup.md`
7. Follow tap setup doc to create `wearetechnative/homebrew-tap` and add secret
8. Run `./release.sh` for first release

## Open Questions

- Should `flake.lock` be pinned to a newer nixpkgs than `nixos-24.05`? (Out of scope for this change, but worth a future task.)
