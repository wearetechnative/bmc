## Context

`bmc` currently has one source of documentation: `README.md`. It is comprehensive but monolithic — a single long file with no navigation, search, or deep-linking. As the project grows and becomes more public-facing, discoverability and usability suffer.

The goal is a proper docs site at `bmc.technative.cloud` that becomes the primary reference, with `README.md` reduced to a short introduction with a link.

## Goals / Non-Goals

**Goals:**
- Ship a Hugo site in `docs/` with TechNative branding
- Deploy automatically to GitHub Pages on every push to `main`
- Serve at `bmc.technative.cloud` via custom domain
- Cover all existing functionality across multiple navigable pages
- Match TechNative visual identity: navy `#1C3D52`, gold `#D4921A`, `tn-bmc.png` logo

**Non-Goals:**
- Versioned docs per release (always tracks `main`)
- Search backend (Congo's built-in client-side search is sufficient)
- Internationalisation

## Decisions

### Theme: Congo

**Decision**: Use the [Congo](https://jpanther.github.io/congo/) Hugo theme.

**Rationale**: Congo is modern, actively maintained, and supports custom color schemes via CSS variables. It has a clean sidebar + top-nav layout suited for CLI documentation. It is lighter than Docsy and more doc-oriented than PaperMod.

**Alternative considered**: Hextra — similar quality but Congo has better branding flexibility and more layout options (hero, background layouts for the landing page).

### Congo installed as Hugo module

**Decision**: Add Congo via `hugo mod` (Go modules), not as a git submodule.

**Rationale**: Cleaner dependency management, no nested git repo, version-pinned in `go.mod`. Requires Go to be available in the GitHub Actions runner, which it is by default.

### Custom color scheme

**Decision**: Create `docs/assets/css/schemes/technative.css` defining CSS variables for navy + gold.

**Rationale**: Congo's color system maps to Tailwind palette names by default, but none match TN colors precisely. A custom scheme file is the documented way to override.

Colors:
- `--color-primary`: `#1C3D52` (navy)
- Accent/links: `#D4921A` (gold)
- Congo config: `colorScheme = "technative"`

### Page structure

```
/                      Landing (hero layout, install snippet, quick links)
/installation/         Homebrew, Nix, binary, NixOS
/setup/
  shell-integration/   install-shell-integration, eval wrapper
  configuration/       config.json fields table
  mfa/                 MFA, TOTP, clipboard (wl-copy / xclip)
/commands/
  profsel/
  console/
  ec2/                 ec2, ec2ls (+ --json), ec2connect, ec2find (+ --json),
                       ec2stopstart, ec2scheduler
  ecsconnect/
/advanced/
  chrome-profiles/
  nixos-home-manager/
  migration/           bash → Go version
```

### GitHub Actions workflow

**Decision**: Use `peaceiris/actions-hugo` + `peaceiris/actions-gh-pages`.

**Rationale**: Established, well-documented actions for Hugo + GitHub Pages. The workflow triggers on push to `main` when files under `docs/**` change (path filter to avoid rebuilding on unrelated commits).

### README strategy

**Decision**: Keep README short — badge, one-liner description, install snippet, link to `bmc.technative.cloud`.

**Rationale**: README is still useful for GitHub landing page context, but docs site is the canonical reference. Duplication causes drift.

## Risks / Trade-offs

- **DNS dependency**: `bmc.technative.cloud` CNAME must be set before the custom domain works. GitHub Pages will serve on `wearetechnative.github.io/bmc` immediately, custom domain once DNS propagates. → Document DNS step clearly; site is usable without it.
- **Hugo module requires Go in CI**: The GitHub Actions runner has Go pre-installed; no extra setup needed.
- **Congo version pinning**: Theme updates may change visual output. → Pin to a specific tag in `go.mod`.

## Migration Plan

1. Add DNS CNAME: `bmc` → `wearetechnative.github.io` on `technative.cloud`
2. Enable GitHub Pages in repo settings, set source to `gh-pages` branch, set custom domain to `bmc.technative.cloud`, enable HTTPS
3. Push `docs/` + workflow → GitHub Actions builds and deploys automatically
4. Verify at `bmc.technative.cloud`
5. Update README to point to docs site

Rollback: disable GitHub Pages in repo settings; `docs/` directory can be removed without affecting the Go binary.
