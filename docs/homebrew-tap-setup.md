# Homebrew Tap Setup

The release workflow uses the `wearetechnative-releaser` GitHub App to push Homebrew formulas to `wearetechnative/homebrew-tap`. The app credentials are stored as **organization-level secrets**, so this setup works for all repos in the org.

## One-time org setup (already done)

The following are set up once and reused by every repo:

1. **GitHub App**: `wearetechnative-releaser` is installed on `wearetechnative/homebrew-tap` with Contents (write) permission
2. **Org secrets** at `github.com/organizations/wearetechnative/settings/secrets/actions`:
   - `RELEASER_APP_ID` — the App ID of `wearetechnative-releaser`
   - `RELEASER_PRIVATE_KEY` — the private key (`.pem` content) of the app

## Per-repo setup (for each new Go repo)

To add Homebrew distribution to another repo:

1. **Install the GitHub App on the new formula repo** (if it's a separate tap):
   - Go to the app settings → Install App → select the target repo

2. **Copy the workflow step** into the repo's release workflow:

```yaml
- name: Generate Homebrew tap token
  uses: actions/create-github-app-token@v1
  id: tap-token
  with:
    app-id: ${{ secrets.RELEASER_APP_ID }}
    private-key: ${{ secrets.RELEASER_PRIVATE_KEY }}
    owner: wearetechnative
    repositories: homebrew-tap

- name: Run GoReleaser
  uses: goreleaser/goreleaser-action@v6
  with:
    args: release --clean
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    HOMEBREW_TAP_TOKEN: ${{ steps.tap-token.outputs.token }}
```

3. **Configure GoReleaser** (`.goreleaser.yml`) with a `brews:` section pointing to `wearetechnative/homebrew-tap` using `{{ .Env.HOMEBREW_TAP_TOKEN }}`.

## How it works

- `actions/create-github-app-token@v1` generates a short-lived token (1 hour) scoped to `homebrew-tap`
- GoReleaser uses that token to push `Formula/bmc.rb` to the tap repo
- No personal access tokens involved — the app is owned by the org

## Install bmc via Homebrew

```sh
# One-liner (taps and installs in one step):
brew install wearetechnative/tap/bmc

# Or tap once, then install by name:
brew tap wearetechnative/tap
brew install bmc
```
