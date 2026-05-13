## 1. Hugo site scaffold

- [x] 1.1 Create `docs/` directory and initialise Hugo site (`hugo new site docs --force`)
- [x] 1.2 Initialise Go module inside `docs/` (`hugo mod init github.com/wearetechnative/bmc/docs`)
- [x] 1.3 Add Congo as Hugo module dependency in `docs/go.mod` (pin to latest stable tag)
- [x] 1.4 Create `docs/hugo.toml` with baseURL, title, theme, and Congo module config

## 2. TechNative branding

- [x] 2.1 Create `docs/assets/css/schemes/technative.css` with CSS variables for `#1C3D52` (navy) and `#D4921A` (gold)
- [x] 2.2 Copy `tn-bmc.png` to `docs/static/logo.png`
- [x] 2.3 Configure Congo in `hugo.toml`: `colorScheme = "technative"`, logo, site title, description
- [x] 2.4 Create `docs/static/CNAME` containing `bmc.technative.cloud`

## 3. Content — Landing page

- [x] 3.1 Create `docs/content/_index.md` — landing page with hero layout, one-liner, quick install snippet, links to main sections

## 4. Content — Installation

- [x] 4.1 Create `docs/content/installation/_index.md` — Homebrew, Nix (nix-env, nix profile, NixOS config.nix), binary download

## 5. Content — Setup

- [x] 5.1 Create `docs/content/setup/_index.md` — section overview
- [x] 5.2 Create `docs/content/setup/shell-integration.md` — install-shell-integration, eval wrapper, NixOS/home-manager snippets
- [x] 5.3 Create `docs/content/setup/configuration.md` — config.json fields reference table
- [x] 5.4 Create `docs/content/setup/mfa.md` — MFA setup, TOTP script, clipboard integration (wl-copy, xclip)

## 6. Content — Commands

- [x] 6.1 Create `docs/content/commands/_index.md` — commands overview
- [x] 6.2 Create `docs/content/commands/profsel.md` — profsel flags and usage
- [x] 6.3 Create `docs/content/commands/console.md` — console flags, Firefox containers, Chrome profiles
- [x] 6.4 Create `docs/content/commands/ec2.md` — ec2, ec2ls (+ --json), ec2connect, ec2find (+ --json), ec2stopstart, ec2scheduler
- [x] 6.5 Create `docs/content/commands/ecsconnect.md` — ecsconnect usage and prerequisites

## 7. Content — Advanced

- [x] 7.1 Create `docs/content/advanced/_index.md` — advanced overview
- [x] 7.2 Create `docs/content/advanced/chrome-profiles.md` — Chrome profile isolation feature
- [x] 7.3 Create `docs/content/advanced/nixos-home-manager.md` — NixOS and home-manager integration details
- [x] 7.4 Create `docs/content/advanced/migration.md` — migrating from bash version

## 8. GitHub Actions workflow

- [x] 8.1 Create `.github/workflows/docs.yml` — build with Hugo, deploy to `gh-pages` branch using `peaceiris/actions-hugo` and `peaceiris/actions-gh-pages`, triggered on push to `main` with path filter `docs/**`

## 9. README update

- [x] 9.1 Shorten `README.md` to: badge, one-liner, quick install snippet, prominent link to `https://bmc.technative.cloud`, brief "full docs" note
- [x] 9.2 Update CHANGELOG under `## NEXT VERSION`
