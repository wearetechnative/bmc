## 1. Hugo site scaffold

- [ ] 1.1 Create `docs/` directory and initialise Hugo site (`hugo new site docs --force`)
- [ ] 1.2 Initialise Go module inside `docs/` (`hugo mod init github.com/wearetechnative/bmc/docs`)
- [ ] 1.3 Add Congo as Hugo module dependency in `docs/go.mod` (pin to latest stable tag)
- [ ] 1.4 Create `docs/hugo.toml` with baseURL, title, theme, and Congo module config

## 2. TechNative branding

- [ ] 2.1 Create `docs/assets/css/schemes/technative.css` with CSS variables for `#1C3D52` (navy) and `#D4921A` (gold)
- [ ] 2.2 Copy `tn-bmc.png` to `docs/static/logo.png`
- [ ] 2.3 Configure Congo in `hugo.toml`: `colorScheme = "technative"`, logo, site title, description
- [ ] 2.4 Create `docs/static/CNAME` containing `bmc.technative.cloud`

## 3. Content — Landing page

- [ ] 3.1 Create `docs/content/_index.md` — landing page with hero layout, one-liner, quick install snippet, links to main sections

## 4. Content — Installation

- [ ] 4.1 Create `docs/content/installation/_index.md` — Homebrew, Nix (nix-env, nix profile, NixOS config.nix), binary download

## 5. Content — Setup

- [ ] 5.1 Create `docs/content/setup/_index.md` — section overview
- [ ] 5.2 Create `docs/content/setup/shell-integration.md` — install-shell-integration, eval wrapper, NixOS/home-manager snippets
- [ ] 5.3 Create `docs/content/setup/configuration.md` — config.json fields reference table
- [ ] 5.4 Create `docs/content/setup/mfa.md` — MFA setup, TOTP script, clipboard integration (wl-copy, xclip)

## 6. Content — Commands

- [ ] 6.1 Create `docs/content/commands/_index.md` — commands overview
- [ ] 6.2 Create `docs/content/commands/profsel.md` — profsel flags and usage
- [ ] 6.3 Create `docs/content/commands/console.md` — console flags, Firefox containers, Chrome profiles
- [ ] 6.4 Create `docs/content/commands/ec2.md` — ec2, ec2ls (+ --json), ec2connect, ec2find (+ --json), ec2stopstart, ec2scheduler
- [ ] 6.5 Create `docs/content/commands/ecsconnect.md` — ecsconnect usage and prerequisites

## 7. Content — Advanced

- [ ] 7.1 Create `docs/content/advanced/_index.md` — advanced overview
- [ ] 7.2 Create `docs/content/advanced/chrome-profiles.md` — Chrome profile isolation feature
- [ ] 7.3 Create `docs/content/advanced/nixos-home-manager.md` — NixOS and home-manager integration details
- [ ] 7.4 Create `docs/content/advanced/migration.md` — migrating from bash version

## 8. GitHub Actions workflow

- [ ] 8.1 Create `.github/workflows/docs.yml` — build with Hugo, deploy to `gh-pages` branch using `peaceiris/actions-hugo` and `peaceiris/actions-gh-pages`, triggered on push to `main` with path filter `docs/**`

## 9. README update

- [ ] 9.1 Shorten `README.md` to: badge, one-liner, quick install snippet, prominent link to `https://bmc.technative.cloud`, brief "full docs" note
- [ ] 9.2 Update CHANGELOG under `## NEXT VERSION`
