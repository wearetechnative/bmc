## Why

`bmc` has no public documentation beyond the README, making it hard for new users to discover, install, and learn the tool. A dedicated docs site at `bmc.technative.cloud` provides a navigable, searchable, branded reference that grows with the project.

Bean: [bmc-d4m6](../../../.beans/bmc-d4m6--project-documentation.md)

## What Changes

- New Hugo site in `docs/` using the Congo theme with TechNative branding (`#1C3D52` navy + `#D4921A` gold, `tn-bmc.png` logo)
- Hosted on GitHub Pages at `bmc.technative.cloud` via custom CNAME
- GitHub Actions workflow (`.github/workflows/docs.yml`) builds and deploys on every push to `main`
- Multi-page structure covering: landing, installation, setup, all commands, advanced topics
- `README.md` shortened to essentials with a link to the docs site as the primary reference

## Capabilities

### New Capabilities

- `docs-site`: Public Hugo documentation site with TechNative branding, multi-page navigation, and GitHub Pages deployment

### Modified Capabilities

_(none — no existing spec-level behaviour changes)_

## Impact

- New: `docs/` directory (Hugo site source)
- New: `.github/workflows/docs.yml` (build + deploy workflow)
- New: `docs/static/CNAME` (`bmc.technative.cloud`)
- New: `docs/static/logo.png` (copied from `tn-bmc.png`)
- Modified: `README.md` (shortened, points to docs site)
- DNS: CNAME record `bmc` → `wearetechnative.github.io` on `technative.cloud` (manual step)
- No Go code changes
