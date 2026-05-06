## ADDED Requirements

### Requirement: Release script exists and is executable
A `release.sh` script SHALL exist at the repository root and be executable. It SHALL guide a developer through a complete release interactively.

#### Scenario: Script is found and runs
- **WHEN** a developer runs `./release.sh` from the repository root
- **THEN** the script starts and presents the version bump menu

### Requirement: Pre-flight checks before release
The script SHALL verify preconditions before making any changes.

#### Scenario: Uncommitted changes detected
- **WHEN** `git status --porcelain` returns output
- **THEN** the script prints an error and exits without modifying any files

#### Scenario: Remote is unreachable
- **WHEN** `git ls-remote origin` fails
- **THEN** the script prints an error and exits

#### Scenario: VERSION-bmc missing
- **WHEN** `VERSION-bmc` does not exist
- **THEN** the script prints an error and exits

### Requirement: Semver version bump
The script SHALL read `VERSION-bmc`, normalize it to 3-part semver (dropping any legacy 4th component), and present patch/minor/major bump options.

#### Scenario: Version bump selection
- **WHEN** the developer selects a bump type (patch/minor/major)
- **THEN** the script calculates the new version and confirms before proceeding

#### Scenario: Tag already exists
- **WHEN** a git tag for the new version already exists
- **THEN** the script prints an error and exits without making changes

### Requirement: CHANGELOG finalization
The script SHALL update `CHANGELOG.md` by replacing `## NEXT VERSION` with the new version and release date.

#### Scenario: NEXT VERSION section present
- **WHEN** `CHANGELOG.md` contains `## NEXT VERSION`
- **THEN** the script replaces it with `## [X.Y.Z] - DD Mon YYYY`

#### Scenario: NEXT VERSION section missing
- **WHEN** `CHANGELOG.md` does not contain `## NEXT VERSION`
- **THEN** the script prints an error and exits without making changes

### Requirement: Nix vendorHash auto-update
When Go module dependencies (`go.mod` or `go.sum`) have changed since the last release tag, the script SHALL automatically recalculate and update the `vendorHash` in `package.nix`.

#### Scenario: Dependencies unchanged
- **WHEN** `go.mod` and `go.sum` are identical to the last release tag
- **THEN** the script skips vendorHash recalculation and continues

#### Scenario: Dependencies changed and Nix available
- **WHEN** `go.mod` or `go.sum` differ from the last release tag AND `nix` is in PATH
- **THEN** the script zeroes vendorHash in `package.nix`, runs `nix build`, parses the correct `sha256-` hash from the error output, updates `package.nix`, and stages the file for commit

#### Scenario: Dependencies changed but Nix unavailable
- **WHEN** `go.mod` or `go.sum` differ from the last release tag AND `nix` is NOT in PATH
- **THEN** the script prints a warning with manual instructions and continues without failing

### Requirement: Local Go build verification
The script SHALL build the binary locally with the new version to verify there are no compilation errors.

#### Scenario: Build succeeds
- **WHEN** `go build` with version ldflags succeeds
- **THEN** the script verifies the binary reports the correct version and removes the test binary

#### Scenario: Build fails
- **WHEN** `go build` exits non-zero
- **THEN** the script prints an error and exits without committing

### Requirement: Release commit and tag creation
The script SHALL create a git commit and annotated tag containing the release contents.

#### Scenario: Commit created
- **WHEN** the developer confirms the diff
- **THEN** the script commits `VERSION-bmc`, `CHANGELOG.md`, and any staged Nix files with message `release: bump version to X.Y.Z`

#### Scenario: Annotated tag created
- **WHEN** the commit is created
- **THEN** the script creates an annotated tag `vX.Y.Z` with the CHANGELOG section as the tag message

#### Scenario: Developer declines commit
- **WHEN** the developer answers "n" to the commit confirmation
- **THEN** the script rolls back all file changes and exits cleanly

### Requirement: Optional push to remote
The script SHALL offer to push both the commit and the tag to origin, but SHALL NOT push without explicit confirmation.

#### Scenario: Developer confirms push
- **WHEN** the developer answers "y" to the push confirmation
- **THEN** the script runs `git push origin main` and `git push origin vX.Y.Z`

#### Scenario: Developer declines push
- **WHEN** the developer answers "n" to the push confirmation
- **THEN** the script prints the push commands for the developer to run manually and exits

### Requirement: Version injection via ldflags in all build paths
The version SHALL be injectable via ldflags so that GoReleaser and Nix builds stamp the binary correctly. `main.go` SHALL prefer a ldflags-set version over the embedded file value.

#### Scenario: Binary built with ldflags
- **WHEN** a binary is built with `-X github.com/wearetechnative/bmc/cmd.Version=X.Y.Z`
- **THEN** `bmc version` reports `X.Y.Z`

#### Scenario: Binary built without ldflags (local go build)
- **WHEN** a binary is built with `go build` and no ldflags
- **THEN** `bmc version` reports the version embedded from `VERSION-bmc`

### Requirement: GitHub Actions release workflow
A `.github/workflows/release.yml` workflow SHALL trigger GoReleaser when a `v*` tag is pushed, building multi-platform binaries and updating the Homebrew tap.

#### Scenario: Tag push triggers release
- **WHEN** a tag matching `v*.*.*` is pushed to the repository
- **THEN** GitHub Actions runs GoReleaser, publishes binaries, and pushes a formula update to `wearetechnative/homebrew-tap`
