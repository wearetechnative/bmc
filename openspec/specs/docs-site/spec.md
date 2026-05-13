## ADDED Requirements

### Requirement: Hugo docs site exists in repository
The repository SHALL contain a Hugo site in the `docs/` directory using the Congo theme.

#### Scenario: Site directory present
- **WHEN** the repository is cloned
- **THEN** a `docs/` directory SHALL exist containing a valid Hugo site with `hugo.toml` and Congo configured as a Hugo module

#### Scenario: Congo installed as Hugo module
- **WHEN** `hugo mod tidy` is run inside `docs/`
- **THEN** Congo SHALL be resolved as a Go module dependency without errors

### Requirement: TechNative branding applied
The docs site SHALL use TechNative visual identity consistently.

#### Scenario: Custom color scheme active
- **WHEN** the site is built
- **THEN** the primary color SHALL be `#1C3D52` (dark navy) and the accent/link color SHALL be `#D4921A` (gold)

#### Scenario: Logo displayed in navigation
- **WHEN** any page is loaded
- **THEN** the `tn-bmc.png` logo SHALL appear in the site header/navigation

#### Scenario: Site title and description
- **WHEN** any page is loaded
- **THEN** the site title SHALL be "BMC" and the description SHALL reference TechNative

### Requirement: Multi-page navigation structure
The docs site SHALL provide clear navigation across all content sections.

#### Scenario: Landing page with quick-start
- **WHEN** a user visits `bmc.technative.cloud/`
- **THEN** they SHALL see a landing page with a one-liner description, a quick install command, and links to main sections

#### Scenario: Installation section
- **WHEN** a user navigates to `/installation/`
- **THEN** they SHALL find install instructions for Homebrew, Nix (nix-env, nix profile, NixOS), and binary download

#### Scenario: Setup section with sub-pages
- **WHEN** a user navigates to `/setup/`
- **THEN** they SHALL find sub-pages for: shell integration, configuration (config.json reference table), and MFA setup

#### Scenario: Commands section with per-command pages
- **WHEN** a user navigates to `/commands/`
- **THEN** they SHALL find dedicated pages for: profsel, console, ec2 (covering ec2, ec2ls, ec2connect, ec2find, ec2stopstart, ec2scheduler including --json flags), and ecsconnect

#### Scenario: Advanced section
- **WHEN** a user navigates to `/advanced/`
- **THEN** they SHALL find pages for: Chrome profile isolation, NixOS/home-manager setup, and migration from the bash version

### Requirement: Automated deployment via GitHub Actions
The docs site SHALL be built and deployed automatically on every push to `main`.

#### Scenario: Workflow triggers on docs change
- **WHEN** a commit is pushed to `main` that modifies files under `docs/`
- **THEN** the GitHub Actions workflow SHALL trigger, build the Hugo site, and push the output to the `gh-pages` branch

#### Scenario: Site available after successful workflow
- **WHEN** the GitHub Actions workflow completes successfully
- **THEN** the updated site SHALL be live at `bmc.technative.cloud`

### Requirement: Custom domain configured
The site SHALL be served at `bmc.technative.cloud`.

#### Scenario: CNAME file present in build output
- **WHEN** the Hugo site is built
- **THEN** a `CNAME` file containing `bmc.technative.cloud` SHALL be present in the output

#### Scenario: HTTPS enforced
- **WHEN** a user visits `http://bmc.technative.cloud`
- **THEN** they SHALL be redirected to `https://bmc.technative.cloud` automatically

### Requirement: README points to docs site
The repository README SHALL direct users to the docs site as the primary reference.

#### Scenario: README contains docs link
- **WHEN** a user reads the GitHub repository README
- **THEN** they SHALL find a prominent link to `https://bmc.technative.cloud` in the first screen of content
