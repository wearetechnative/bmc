## ADDED Requirements

### Requirement: GoReleaser build matrix
The system SHALL be built for four targets via GoReleaser: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`. Binaries SHALL be named `bmc`.

#### Scenario: Release triggered
- **WHEN** a git tag is pushed matching `v*`
- **THEN** GoReleaser builds all four targets and publishes to GitHub Releases

### Requirement: Homebrew tap distribution
The system SHALL publish a Homebrew formula to `wearetechnative/homebrew-tap` automatically via GoReleaser on each release.

#### Scenario: Homebrew install
- **WHEN** user runs `brew install wearetechnative/tap/bmc`
- **THEN** Homebrew installs the correct binary for the user's platform

### Requirement: Nix flake distribution
The system SHALL provide a `flake.nix` using `buildGoModule` that can be used in three ways: `nix-env`, `nix profile add`, and NixOS `environment.systemPackages`.

#### Scenario: nix-env install
- **WHEN** user runs `nix-env -iA bmc` from the flake
- **THEN** bmc binary is installed in the user's profile

#### Scenario: nix profile install
- **WHEN** user runs `nix profile add github:wearetechnative/bmc`
- **THEN** bmc binary is installed via flakes

#### Scenario: NixOS configuration
- **WHEN** user adds the flake as an input and includes `bmc` in `environment.systemPackages`
- **THEN** bmc is available system-wide after `nixos-rebuild switch`

### Requirement: VERSION-bmc file
The system SHALL read the version from the `VERSION-bmc` file at build time (embedded via Go's `embed` package or ldflags) so `bmc version` displays the correct version.

#### Scenario: Version output
- **WHEN** user runs `bmc version`
- **THEN** system displays the version from `VERSION-bmc` in the same format as the current bash implementation
