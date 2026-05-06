# Homebrew Tap Setup

One-time setup to enable `brew install wearetechnative/tap/bmc`.

## 1. Create the tap repository

1. Go to [github.com/wearetechnative](https://github.com/wearetechnative)
2. Create a new **public** repository named `homebrew-tap`
3. Initialize with a README (so the repo is not empty)
4. Create a `Formula/` directory in the repository (add a `.gitkeep` file if needed)

The repository must be at `github.com/wearetechnative/homebrew-tap` for GoReleaser's `brews:` config to work.

## 2. Generate a Personal Access Token

1. Go to **GitHub Settings → Developer settings → Personal access tokens → Fine-grained tokens**
2. Click **Generate new token**
3. Set:
   - **Resource owner**: `wearetechnative`
   - **Repository access**: Only selected repositories → `wearetechnative/homebrew-tap`
   - **Permissions**: Contents → Read and write
4. Generate and copy the token

## 3. Add the token as a repository secret

1. Go to `github.com/wearetechnative/bmc` → **Settings → Secrets and variables → Actions**
2. Click **New repository secret**
3. Name: `HOMEBREW_TAP_TOKEN`
4. Value: paste the token from step 2
5. Click **Add secret**

## 4. Verify

After the next release (`./release.sh` + push), GoReleaser will automatically create or update `Formula/bmc.rb` in `wearetechnative/homebrew-tap`.

Users can then install bmc with:

```sh
brew tap wearetechnative/tap
brew install bmc
```

## Notes

- The `GITHUB_TOKEN` secret is provided automatically by GitHub Actions — no setup needed
- The `HOMEBREW_TAP_TOKEN` must be a PAT because the default `GITHUB_TOKEN` cannot write to other repositories
- If the secret is missing when a release is pushed, GoReleaser will fail at the Homebrew step; the binaries will still be published to GitHub Releases
