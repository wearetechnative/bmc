## ADDED Requirements

### Requirement: Homebrew tap setup documented
A `docs/homebrew-tap-setup.md` document SHALL exist describing the one-time steps to create and configure the `wearetechnative/homebrew-tap` GitHub repository and the required `HOMEBREW_TAP_TOKEN` secret.

#### Scenario: Developer follows setup guide
- **WHEN** a developer follows `docs/homebrew-tap-setup.md` from scratch
- **THEN** they can complete the setup (create tap repo, generate PAT, add secret) without needing additional information

### Requirement: GoReleaser configured to publish to Homebrew tap
`.goreleaser.yml` SHALL have a `brews:` section that pushes a formula to `wearetechnative/homebrew-tap` using the `HOMEBREW_TAP_TOKEN` environment variable.

#### Scenario: Release publishes formula
- **WHEN** GoReleaser runs on a tagged release
- **THEN** it creates or updates `Formula/bmc.rb` in `wearetechnative/homebrew-tap` with the correct version, URL, and sha256 checksum

### Requirement: bmc installable via Homebrew tap
After setup, users SHALL be able to install bmc using:
```
brew tap wearetechnative/tap
brew install bmc
```

#### Scenario: Install via Homebrew
- **WHEN** a user runs `brew tap wearetechnative/tap && brew install bmc`
- **THEN** the `bmc` binary is installed and `bmc version` reports the correct release version
