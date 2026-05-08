#!/usr/bin/env bash
set -euo pipefail

# ── helpers ────────────────────────────────────────────────────────────────────
red()   { printf '\033[31m%s\033[0m\n' "$*"; }
green() { printf '\033[32m%s\033[0m\n' "$*"; }
bold()  { printf '\033[1m%s\033[0m\n' "$*"; }
die()   { red "ERROR: $*"; exit 1; }

confirm() {
  local answer
  printf '%s [y/N] ' "$1"
  read -r answer
  [[ "$answer" =~ ^[Yy]$ ]]
}

# ── pre-flight checks ──────────────────────────────────────────────────────────
bold "==> Pre-flight checks"

[[ -f VERSION-bmc ]] || die "VERSION-bmc not found"

if [[ -n "$(git status --porcelain --untracked-files=no)" ]]; then
  die "Uncommitted changes detected. Commit or stash them first."
fi

git ls-remote origin HEAD &>/dev/null || die "Remote 'origin' is unreachable"

green "  Pre-flight checks passed"

# ── version bump ───────────────────────────────────────────────────────────────
bold "==> Version bump"

RAW_VERSION="$(cat VERSION-bmc | tr -d '[:space:]')"

# Normalize: strip legacy 4th component (e.g. 0.2.12.0 → 0.2.12)
DOT_COUNT="$(echo "$RAW_VERSION" | tr -cd '.' | wc -c)"
if [[ "$DOT_COUNT" -ge 3 ]]; then
  CURRENT_VERSION="${RAW_VERSION%.*}"
else
  CURRENT_VERSION="$RAW_VERSION"
fi

IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

printf '  Current version: %s' "$CURRENT_VERSION"
if [[ "$RAW_VERSION" != "$CURRENT_VERSION" ]]; then
  printf ' (normalized from %s)' "$RAW_VERSION"
fi
printf '\n\n'

echo "  Select bump type:"
echo "  1) patch  →  $MAJOR.$MINOR.$((PATCH + 1))"
echo "  2) minor  →  $MAJOR.$((MINOR + 1)).0"
echo "  3) major  →  $((MAJOR + 1)).0.0"
echo ""
printf '  Choice [1-3]: '
read -r BUMP_CHOICE

case "$BUMP_CHOICE" in
  1) NEW_VERSION="$MAJOR.$MINOR.$((PATCH + 1))" ;;
  2) NEW_VERSION="$MAJOR.$((MINOR + 1)).0" ;;
  3) NEW_VERSION="$((MAJOR + 1)).0.0" ;;
  *) die "Invalid choice: $BUMP_CHOICE" ;;
esac

printf '  New version: '
bold "$NEW_VERSION"
echo ""

confirm "  Proceed with version $NEW_VERSION?" || { echo "Aborted."; exit 0; }

# ── tag-exists guard ───────────────────────────────────────────────────────────
if git tag --list "v$NEW_VERSION" | grep -q .; then
  die "Tag v$NEW_VERSION already exists"
fi

# ── CHANGELOG update ───────────────────────────────────────────────────────────
bold "==> Updating CHANGELOG"

grep -q "^## NEXT VERSION" CHANGELOG.md \
  || die "CHANGELOG.md has no '## NEXT VERSION' section — add your entries first"

RELEASE_DATE="$(date '+%d %b %Y')"
sed -i "0,/^## NEXT VERSION/s/^## NEXT VERSION/## [$NEW_VERSION] - $RELEASE_DATE/" CHANGELOG.md
green "  CHANGELOG.md: ## [$NEW_VERSION] - $RELEASE_DATE"

# ── update VERSION-bmc ─────────────────────────────────────────────────────────
echo "$NEW_VERSION" > VERSION-bmc
green "  VERSION-bmc: $NEW_VERSION"

# ── vendorHash auto-update ─────────────────────────────────────────────────────
bold "==> Checking Nix vendorHash"

LAST_TAG="$(git describe --tags --abbrev=0 2>/dev/null || echo "")"
NEEDS_HASH_UPDATE=false

if [[ -z "$LAST_TAG" ]]; then
  NEEDS_HASH_UPDATE=true
  echo "  No previous tag found — will recalculate vendorHash"
elif git diff --quiet "$LAST_TAG" -- go.mod go.sum 2>/dev/null; then
  green "  go.mod/go.sum unchanged since $LAST_TAG — skipping vendorHash update"
else
  NEEDS_HASH_UPDATE=true
  echo "  go.mod or go.sum changed since $LAST_TAG — recalculating vendorHash"
fi

if [[ "$NEEDS_HASH_UPDATE" == true ]]; then
  if ! command -v nix &>/dev/null; then
    printf '\033[33mWARNING: nix not in PATH — cannot auto-update vendorHash\033[0m\n'
    echo "  Update vendorHash manually before releasing:"
    echo "    1. Set: vendorHash = lib.fakeHash; in package.nix"
    echo "    2. Run: nix build .#bmc"
    echo "    3. Copy the 'got: sha256-...' hash into package.nix"
  else
    echo "  Temporarily setting fakeHash in package.nix..."
    # Capture current hash for potential rollback
    ORIGINAL_HASH="$(grep -oP '(?<=vendorHash = ")sha256-[^"]+' package.nix || echo "")"

    sed -i 's|vendorHash = "sha256-[^"]*";|vendorHash = lib.fakeHash;|' package.nix
    sed -i 's|vendorHash = null;|vendorHash = lib.fakeHash;|' package.nix

    echo "  Running nix build to calculate hash (may take a while)..."
    NIX_OUT="$(nix build .#bmc 2>&1 || true)"
    NEW_HASH="$(echo "$NIX_OUT" | grep -oP '(?<=got:    )sha256-\S+' || echo "")"

    if [[ -z "$NEW_HASH" ]]; then
      # Restore original
      if [[ -n "$ORIGINAL_HASH" ]]; then
        sed -i "s|vendorHash = lib.fakeHash;|vendorHash = \"$ORIGINAL_HASH\";|" package.nix
      fi
      die "Could not extract vendorHash from nix build output. Run manually:\n  nix build .#bmc"
    fi

    sed -i "s|vendorHash = lib.fakeHash;|vendorHash = \"$NEW_HASH\";|" package.nix
    green "  vendorHash updated: $NEW_HASH"
  fi
fi

# ── optional flake.lock update ─────────────────────────────────────────────────
bold "==> flake.lock"
if confirm "  Run 'nix flake update'?"; then
  nix flake update
  green "  flake.lock updated"
else
  echo "  Skipped"
fi

# ── local Go build verification ────────────────────────────────────────────────
bold "==> Build verification"

echo "  Building bmc.test with version ldflags..."
go build -ldflags "-X github.com/wearetechnative/bmc/cmd.Version=$NEW_VERSION" -o bmc.test . \
  || { rm -f bmc.test; die "go build failed"; }

REPORTED="$(./bmc.test version 2>&1 | grep -oP '\d+\.\d+\.\d+' | head -1 || echo "")"
rm -f bmc.test

if [[ "$REPORTED" != "$NEW_VERSION" ]]; then
  die "Version mismatch: binary reports '$REPORTED', expected '$NEW_VERSION'"
fi
green "  Binary reports correct version: $NEW_VERSION"

# Rebuild main bmc binary
go build -ldflags "-X github.com/wearetechnative/bmc/cmd.Version=$NEW_VERSION" -o bmc .
green "  bmc binary rebuilt"

# ── release summary and confirmation ──────────────────────────────────────────
bold "==> Release summary"
echo ""
printf '  Version:  %s\n' "$NEW_VERSION"
printf '  Tag:      v%s\n' "$NEW_VERSION"
printf '  Date:     %s\n' "$RELEASE_DATE"
echo ""
git diff -- VERSION-bmc CHANGELOG.md package.nix flake.lock 2>/dev/null || true
echo ""

if ! confirm "  Commit and tag v$NEW_VERSION?"; then
  echo "  Rolling back..."
  git checkout -- VERSION-bmc CHANGELOG.md package.nix flake.lock 2>/dev/null || true
  echo "Aborted."
  exit 0
fi

# ── git commit ─────────────────────────────────────────────────────────────────
bold "==> Creating release commit"

git add VERSION-bmc CHANGELOG.md
git diff --quiet HEAD -- package.nix 2>/dev/null || git add package.nix
git diff --quiet HEAD -- flake.lock 2>/dev/null  || git add flake.lock

git commit -m "release: bump version to $NEW_VERSION"
green "  Commit created"

# ── annotated tag ──────────────────────────────────────────────────────────────
bold "==> Creating tag v$NEW_VERSION"

TAG_MSG="$(awk "/^\#\# \[$NEW_VERSION\]/{found=1; next} found && /^\#\# /{exit} found{print}" CHANGELOG.md)"

git tag -a "v$NEW_VERSION" -m "$(printf 'Release v%s\n\n%s' "$NEW_VERSION" "$TAG_MSG")"
green "  Tag v$NEW_VERSION created"

# ── push ───────────────────────────────────────────────────────────────────────
bold "==> Push to origin"
echo ""
if confirm "  Push commit and tag to origin?"; then
  git push origin main
  git push origin "v$NEW_VERSION"
  green "  Pushed. GitHub Actions will build and publish the release."
else
  echo ""
  echo "  Push manually when ready:"
  bold "    git push origin main"
  bold "    git push origin v$NEW_VERSION"
fi

echo ""
green "Release v$NEW_VERSION complete."
