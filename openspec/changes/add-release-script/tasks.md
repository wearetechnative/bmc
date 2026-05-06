## 1. Version Injection Alignment

- [x] 1.1 Update `main.go`: prefer ldflags-set version over embed — check `if cmd.Version == "dev"` before overwriting with `versionRaw`
- [x] 1.2 Update `package.nix` ldflags: add `-X github.com/wearetechnative/bmc/cmd.Version=${version}` alongside existing `-s -w`
- [x] 1.3 Update `.goreleaser.yml` ldflags: add `-X github.com/wearetechnative/bmc/cmd.Version={{.Version}}`

## 2. Fix Nix vendorHash

- [x] 2.1 Calculate correct `vendorHash` for current `go.mod`/`go.sum`: run `nix build` with `vendorHash = lib.fakeHash` or `""` in `package.nix`, capture the `got: sha256-...` from error output
- [x] 2.2 Update `package.nix` `vendorHash` from `null` to the calculated `sha256-...` hash
- [x] 2.3 Verify Nix build succeeds: `nix build .#bmc`

## 3. Release Script

- [x] 3.1 Create `release.sh` with pre-flight checks: uncommitted changes, remote reachable, `VERSION-bmc` exists
- [x] 3.2 Add version reading and semver normalization: parse `VERSION-bmc`, strip 4th component if present, present patch/minor/major bump menu
- [x] 3.3 Add tag-exists guard: fail if `vX.Y.Z` tag already exists
- [x] 3.4 Add CHANGELOG update: replace `## NEXT VERSION` with `## [X.Y.Z] - DD Mon YYYY`; fail if section not found
- [x] 3.5 Add `package.nix` vendorHash auto-update: detect go.mod/go.sum changes since last tag, temporarily zero hash, run `nix build`, parse correct hash, update `package.nix`
- [x] 3.6 Add optional `flake.lock` update: prompt to run `nix flake update` and stage result
- [x] 3.7 Add local Go build verification: `go build -ldflags "-X github.com/wearetechnative/bmc/cmd.Version=${NEW_VERSION}" -o bmc.test`, verify `./bmc.test version` reports correct version, remove test binary
- [x] 3.8 Add release summary, diff display, and commit confirmation with rollback on decline
- [x] 3.9 Create git commit staging `VERSION-bmc`, `CHANGELOG.md`, and any Nix files; message `release: bump version to X.Y.Z`
- [x] 3.10 Create annotated git tag `vX.Y.Z` with CHANGELOG section as tag message
- [x] 3.11 Prompt to push commit + tag to origin; print manual push commands if declined
- [x] 3.12 Make `release.sh` executable: `chmod +x release.sh`

## 4. GitHub Actions Workflow

- [x] 4.1 Create `.github/workflows/release.yml`: trigger on `push: tags: ['v*.*.*']`, run GoReleaser action, pass `GITHUB_TOKEN` and `HOMEBREW_TAP_TOKEN` secrets

## 5. Homebrew Tap Setup Documentation

- [x] 5.1 Create `docs/homebrew-tap-setup.md`: step-by-step guide to create `wearetechnative/homebrew-tap` repo, required `Formula/` directory structure, generate PAT with `repo` scope, add `HOMEBREW_TAP_TOKEN` secret to `wearetechnative/bmc` repository settings

## 6. Verification

- [x] 6.1 Run `./release.sh` dry-run check: verify pre-flight checks work (e.g., with dirty working tree)
- [x] 6.2 Verify `bmc version` reports correct version after `go build` without ldflags (embed path)
- [x] 6.3 Verify `bmc version` reports correct version after `go build -ldflags "-X ..."` (ldflags path)
