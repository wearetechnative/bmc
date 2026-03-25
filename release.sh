#!/usr/bin/env bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_error "Not in a git repository"
    exit 1
fi

# Check for uncommitted changes
if [[ -n $(git status --porcelain) ]]; then
    print_error "You have uncommitted changes. Please commit or stash them first."
    git status --short
    exit 1
fi

# Check remote connectivity
print_info "Checking remote connectivity..."
if ! git ls-remote --exit-code origin &>/dev/null; then
    print_error "Cannot reach remote repository. Check your network connection."
    exit 1
fi
print_success "Remote is reachable"

# Read current version from VERSION-bmc
if [[ ! -f VERSION-bmc ]]; then
    print_error "VERSION-bmc file not found"
    exit 1
fi

CURRENT_VERSION=$(cat VERSION-bmc | tr -d '[:space:]')
print_info "Current version: ${GREEN}${CURRENT_VERSION}${NC}"

# Parse version components
IFS='.' read -r -a VERSION_PARTS <<< "$CURRENT_VERSION"
MAJOR="${VERSION_PARTS[0]}"
MINOR="${VERSION_PARTS[1]}"
PATCH="${VERSION_PARTS[2]}"
BUILD="${VERSION_PARTS[3]:-0}"

echo ""
echo "Select release type:"
echo "  1) Patch   (${MAJOR}.${MINOR}.$((PATCH + 1)).${BUILD}) - Bug fixes, no new features"
echo "  2) Minor   (${MAJOR}.$((MINOR + 1)).0.${BUILD}) - New features, backwards compatible"
echo "  3) Major   ($((MAJOR + 1)).0.0.${BUILD}) - Breaking changes"
echo ""
read -p "Enter choice (1-3): " RELEASE_TYPE

case $RELEASE_TYPE in
    1)
        RELEASE_NAME="patch"
        NEW_VERSION="${MAJOR}.${MINOR}.$((PATCH + 1)).${BUILD}"
        ;;
    2)
        RELEASE_NAME="minor"
        NEW_VERSION="${MAJOR}.$((MINOR + 1)).0.${BUILD}"
        ;;
    3)
        RELEASE_NAME="major"
        NEW_VERSION="$((MAJOR + 1)).0.0.${BUILD}"
        ;;
    *)
        print_error "Invalid choice"
        exit 1
        ;;
esac

print_info "New version will be: ${GREEN}${NEW_VERSION}${NC} (${RELEASE_NAME} release)"

# Check if tag already exists
if git rev-parse "v${NEW_VERSION}" >/dev/null 2>&1; then
    print_error "Tag v${NEW_VERSION} already exists!"
    print_info "Existing tags:"
    git tag -l | grep -E "v${MAJOR}\.${MINOR}\." | tail -5
    exit 1
fi
print_success "Version tag is available"

# Check for OpenSpec completed changes that need archiving
if command -v openspec &> /dev/null; then
    print_info "Checking OpenSpec changes..."

    COMPLETED_CHANGES=$(openspec list 2>/dev/null | grep "✓ Complete" | awk '{print $1}' || true)

    if [[ -n "$COMPLETED_CHANGES" ]]; then
        print_warning "Found completed OpenSpec changes that may need archiving:"
        echo ""
        openspec list | grep "✓ Complete" || true
        echo ""
        read -p "Archive these changes before release? (y/n): " ARCHIVE_CHANGES

        if [[ $ARCHIVE_CHANGES == "y" || $ARCHIVE_CHANGES == "Y" ]]; then
            print_info "Archiving completed changes..."
            for change in $COMPLETED_CHANGES; do
                print_info "Archiving ${change}..."
                if openspec archive "$change" --yes 2>/dev/null; then
                    print_success "Archived ${change}"
                else
                    print_warning "Could not archive ${change} (may already be archived)"
                fi
            done

            # Commit archived changes if any
            if [[ -n $(git status --porcelain openspec/) ]]; then
                print_info "Committing archived changes..."
                git add openspec/
                git commit -m "chore: archive completed openspec changes before release"
                print_success "Archived changes committed"
            fi
        else
            print_warning "Skipping OpenSpec archiving"
        fi
    else
        print_success "No completed OpenSpec changes to archive"
    fi
else
    print_warning "OpenSpec CLI not found, skipping OpenSpec checks"
fi
echo ""

# Ask for changelog entry
print_info "CHANGELOG.md update"
echo ""
echo "Choose changelog option:"
echo "  1) I have already added entries under '## NEXT VERSION' in CHANGELOG.md"
echo "  2) Enter changelog entries now (interactive)"
echo ""
read -p "Enter choice (1-2): " CHANGELOG_CHOICE

if [[ $CHANGELOG_CHOICE == "2" ]]; then
    echo ""
    print_info "Enter changelog entries (one per line, press Ctrl+D when done):"
    print_info "Format: '- Feature: description' or '- Fix: description' or '- Enhancement: description'"
    echo ""

    CHANGELOG_ENTRIES=""
    while IFS= read -r line; do
        if [[ -n "$line" ]]; then
            CHANGELOG_ENTRIES="${CHANGELOG_ENTRIES}${line}\n"
        fi
    done

    if [[ -z "$CHANGELOG_ENTRIES" ]]; then
        print_error "No changelog entries provided"
        exit 1
    fi
fi

# Get current date
RELEASE_DATE=$(date +"%d %b %Y" | sed 's/ 0/ /g')

# Update CHANGELOG.md
print_info "Updating CHANGELOG.md..."

if [[ $CHANGELOG_CHOICE == "1" ]]; then
    # Replace "## NEXT VERSION" with actual version and date
    if grep -q "## NEXT VERSION" CHANGELOG.md; then
        sed -i "s/## NEXT VERSION/## ${NEW_VERSION} - ${RELEASE_DATE}/" CHANGELOG.md
        print_success "Updated CHANGELOG.md (replaced NEXT VERSION)"
    else
        # No NEXT VERSION found, add new section at top
        print_warning "No '## NEXT VERSION' found in CHANGELOG.md"
        echo "Please add changelog entries manually before continuing."
        exit 1
    fi
else
    # Insert new version section with provided entries
    TEMP_FILE=$(mktemp)
    {
        head -n 2 CHANGELOG.md
        echo ""
        echo "## ${NEW_VERSION} - ${RELEASE_DATE}"
        echo -e "$CHANGELOG_ENTRIES"
        echo ""
        tail -n +3 CHANGELOG.md
    } > "$TEMP_FILE"
    mv "$TEMP_FILE" CHANGELOG.md
    print_success "Updated CHANGELOG.md with new entries"
fi

# Update VERSION-bmc
print_info "Updating VERSION-bmc..."
echo "$NEW_VERSION" > VERSION-bmc
print_success "Updated VERSION-bmc to ${NEW_VERSION}"

# Update package.nix
print_info "Updating package.nix..."
if [[ -f package.nix ]]; then
    sed -i "s/version = \"${CURRENT_VERSION}\";/version = \"${NEW_VERSION}\";/" package.nix
    print_success "Updated package.nix to ${NEW_VERSION}"
else
    print_warning "package.nix not found, skipping"
fi

# Verify Nix flake if available
if [[ -f flake.nix ]] && command -v nix &> /dev/null; then
    print_info "Verifying Nix flake..."

    # Check flake
    if nix flake check 2>&1 | grep -q "error:"; then
        print_warning "Nix flake check found issues (this may be normal)"
    else
        print_success "Nix flake check passed"
    fi

    # Optionally update flake.lock
    echo ""
    read -p "Update flake.lock (nix flake update)? (y/n): " UPDATE_FLAKE

    if [[ $UPDATE_FLAKE == "y" || $UPDATE_FLAKE == "Y" ]]; then
        print_info "Updating flake.lock..."
        nix flake update
        print_success "flake.lock updated"

        if [[ -n $(git status --porcelain flake.lock) ]]; then
            git add flake.lock
            print_info "flake.lock will be included in release commit"
        fi
    fi
else
    if [[ ! -f flake.nix ]]; then
        print_warning "flake.nix not found, skipping Nix checks"
    else
        print_warning "Nix not installed, skipping flake checks"
    fi
fi

# Show summary
echo ""
echo "════════════════════════════════════════════════════════════"
print_info "Release Summary"
echo "════════════════════════════════════════════════════════════"
echo ""
echo "  Release type:    ${RELEASE_NAME}"
echo "  Old version:     ${CURRENT_VERSION}"
echo "  New version:     ${GREEN}${NEW_VERSION}${NC}"
echo "  Release date:    ${RELEASE_DATE}"
echo ""
echo "  Files changed:"
echo "    - CHANGELOG.md"
echo "    - VERSION-bmc"
echo "    - package.nix"
if [[ -n $(git status --porcelain flake.lock 2>/dev/null) ]]; then
    echo "    - flake.lock"
fi
if [[ -n $(git status --porcelain openspec/ 2>/dev/null) ]]; then
    echo "    - openspec/ (archived changes)"
fi
echo ""
echo "════════════════════════════════════════════════════════════"
echo ""

# Show git diff
print_info "Changes to be committed:"
echo ""
git diff HEAD CHANGELOG.md VERSION-bmc package.nix flake.lock 2>/dev/null || git diff CHANGELOG.md VERSION-bmc package.nix
if [[ -n $(git status --porcelain openspec/ 2>/dev/null) ]]; then
    echo ""
    print_info "OpenSpec changes:"
    git status --short openspec/
fi
echo ""

# Confirm before committing
read -p "Commit these changes? (y/n): " CONFIRM_COMMIT

if [[ $CONFIRM_COMMIT != "y" && $CONFIRM_COMMIT != "Y" ]]; then
    print_warning "Aborting release. Rolling back changes..."
    git checkout CHANGELOG.md VERSION-bmc package.nix 2>/dev/null || true
    git checkout flake.lock 2>/dev/null || true
    print_info "Changes rolled back"
    exit 0
fi

# Create git commit
print_info "Creating git commit..."
git add CHANGELOG.md VERSION-bmc package.nix

# Add flake.lock if it was updated
if [[ -n $(git status --porcelain flake.lock 2>/dev/null) ]]; then
    git add flake.lock
fi

COMMIT_MESSAGE="release: bump version to ${NEW_VERSION}

Update VERSION-bmc and CHANGELOG.md for ${RELEASE_NAME} release."

git commit -m "$COMMIT_MESSAGE"
print_success "Commit created"

# Create git tag
print_info "Creating git tag v${NEW_VERSION}..."

TAG_MESSAGE="Release v${NEW_VERSION}

$(sed -n "/## ${NEW_VERSION}/,/## [0-9]/p" CHANGELOG.md | sed '$d' | tail -n +2)"

git tag -a "v${NEW_VERSION}" -m "$TAG_MESSAGE"
print_success "Tag v${NEW_VERSION} created"

# Show final summary
echo ""
echo "════════════════════════════════════════════════════════════"
print_success "Release v${NEW_VERSION} prepared successfully!"
echo "════════════════════════════════════════════════════════════"
echo ""
echo "  Commit: $(git log -1 --pretty=format:'%h - %s')"
echo "  Tag:    v${NEW_VERSION}"
echo ""
echo "════════════════════════════════════════════════════════════"
echo ""

# Confirm before pushing
read -p "Push to remote (commit + tag)? (y/n): " CONFIRM_PUSH

if [[ $CONFIRM_PUSH != "y" && $CONFIRM_PUSH != "Y" ]]; then
    print_warning "Skipping push. You can push later with:"
    echo ""
    echo "  git push origin main"
    echo "  git push origin v${NEW_VERSION}"
    echo ""
    exit 0
fi

# Push to remote
print_info "Pushing to remote..."
git push origin main
git push origin "v${NEW_VERSION}"

echo ""
print_success "Release v${NEW_VERSION} pushed to remote!"
print_success "Done! 🚀"
echo ""
