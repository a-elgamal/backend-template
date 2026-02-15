#!/usr/bin/env bash
set -euo pipefail

# ─── Colors & helpers ────────────────────────────────────────────────────────

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

info()  { printf "${BLUE}▶${NC} %s\n" "$*"; }
ok()    { printf "${GREEN}✔${NC} %s\n" "$*"; }
warn()  { printf "${YELLOW}⚠${NC} %s\n" "$*"; }
err()   { printf "${RED}✖${NC} %s\n" "$*" >&2; }

# ─── Pre-flight checks ──────────────────────────────────────────────────────

# Check if setup already ran
if ! grep -q 'alielgamal\.com/myservice' go.mod 2>/dev/null; then
    err "It looks like setup has already been run (go.mod no longer contains alielgamal.com/myservice)."
    exit 1
fi

# Warn about uncommitted changes
if [[ -n "$(git status --porcelain 2>/dev/null)" ]]; then
    warn "You have uncommitted git changes. Consider committing or stashing before running setup."
    printf "  Continue anyway? [y/N] "
    read -r answer
    if [[ "$answer" != "y" && "$answer" != "Y" ]]; then
        echo "Aborted."
        exit 0
    fi
fi

# ─── Collect inputs ─────────────────────────────────────────────────────────

echo ""
printf "${BOLD}Backend Template Setup Wizard${NC}\n"
echo "─────────────────────────────────────"
echo ""

# 1. Go module path
while true; do
    printf "${BOLD}Go module path${NC} (e.g. github.com/myorg/myapi): "
    read -r MODULE_PATH
    if [[ -z "$MODULE_PATH" ]]; then
        err "Module path is required."
        continue
    fi
    # Basic validation: at least one slash, no spaces
    if [[ ! "$MODULE_PATH" =~ ^[a-zA-Z0-9._-]+(/[a-zA-Z0-9._-]+)+$ ]]; then
        err "Invalid module path format. Expected something like github.com/org/repo"
        continue
    fi
    break
done

# 2. Service name
while true; do
    printf "${BOLD}Service name${NC} (lowercase, alphanumeric + hyphens, e.g. order-api): "
    read -r SERVICE_NAME
    if [[ -z "$SERVICE_NAME" ]]; then
        err "Service name is required."
        continue
    fi
    if [[ ! "$SERVICE_NAME" =~ ^[a-z][a-z0-9-]*$ ]]; then
        err "Service name must start with a letter and contain only lowercase letters, numbers, and hyphens."
        continue
    fi
    break
done

# 3. Project prefix
printf "${BOLD}Project prefix${NC} for cloud resources [${SERVICE_NAME}]: "
read -r PROJECT_PREFIX
PROJECT_PREFIX="${PROJECT_PREFIX:-$SERVICE_NAME}"

# 4. Gist ID
printf "${BOLD}Gist ID${NC} for badge updates (leave empty to skip): "
read -r GIST_ID

# 5. Cloud provider
echo ""
echo "Cloud provider:"
echo "  1) GCP only"
echo "  2) AWS only"
echo "  3) Both"
while true; do
    printf "${BOLD}Choose${NC} [1/2/3]: "
    read -r CLOUD_CHOICE
    case "$CLOUD_CHOICE" in
        1) CLOUD_PROVIDER="gcp" ; break ;;
        2) CLOUD_PROVIDER="aws" ; break ;;
        3) CLOUD_PROVIDER="both"; break ;;
        *) err "Please enter 1, 2, or 3." ;;
    esac
done

# ─── Confirmation ────────────────────────────────────────────────────────────

echo ""
echo "─────────────────────────────────────"
printf "${BOLD}Summary:${NC}\n"
echo "  Module path:    $MODULE_PATH"
echo "  Service name:   $SERVICE_NAME"
echo "  Project prefix: $PROJECT_PREFIX"
echo "  Gist ID:        ${GIST_ID:-<skip>}"
echo "  Cloud provider: $CLOUD_PROVIDER"
echo "─────────────────────────────────────"
printf "Proceed? [y/N] "
read -r confirm
if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    echo "Aborted."
    exit 0
fi

echo ""

# ─── Execute replacements ───────────────────────────────────────────────────

# Helper: cross-platform sed in-place
sedi() {
    if [[ "$(uname)" == "Darwin" ]]; then
        sed -i '' "$@"
    else
        sed -i "$@"
    fi
}

# 1. Replace Go module path
info "Replacing Go module path..."
sedi "s|alielgamal.com/myservice|${MODULE_PATH}|g" go.mod
find . -name '*.go' -not -path './.git/*' | while IFS= read -r f; do
    sedi "s|alielgamal.com/myservice|${MODULE_PATH}|g" "$f"
done
ok "Module path updated."

# 2. Replace project prefix in terraform and workflows
info "Replacing project prefix..."
find ./terraform -name '*.tf' -not -path './.git/*' | while IFS= read -r f; do
    sedi "s|knz-myservice|${PROJECT_PREFIX}|g" "$f"
done
find ./.github/workflows -name '*.yml' | while IFS= read -r f; do
    sedi "s|knz-myservice|${PROJECT_PREFIX}|g" "$f"
done
ok "Project prefix updated."

# 3. Rename terraform module directories
info "Renaming terraform module directories..."
if [[ -d "terraform/modules/myservice-aws" ]]; then
    mv "terraform/modules/myservice-aws" "terraform/modules/${SERVICE_NAME}-aws"
fi
if [[ -d "terraform/modules/myservice" ]]; then
    mv "terraform/modules/myservice" "terraform/modules/${SERVICE_NAME}"
fi
ok "Terraform modules renamed."

# 4. Replace service name across all relevant files
info "Replacing service name..."
FILES_TO_UPDATE=(
    Makefile
    Dockerfile
    docker-compose.local.yml
    application.yaml
    CLAUDE.md
    README.md
    main.go
)

for f in "${FILES_TO_UPDATE[@]}"; do
    if [[ -f "$f" ]]; then
        sedi "s|myservice|${SERVICE_NAME}|g" "$f"
    fi
done

# cmd/*.go
find ./cmd -name '*.go' | while IFS= read -r f; do
    sedi "s|myservice|${SERVICE_NAME}|g" "$f"
done

# internal/**/*.go
find ./internal -name '*.go' | while IFS= read -r f; do
    sedi "s|myservice|${SERVICE_NAME}|g" "$f"
done

# terraform/**/*.tf
find ./terraform -name '*.tf' | while IFS= read -r f; do
    sedi "s|myservice|${SERVICE_NAME}|g" "$f"
done

# .github/workflows/*.yml
find ./.github/workflows -name '*.yml' | while IFS= read -r f; do
    sedi "s|myservice|${SERVICE_NAME}|g" "$f"
done

# portal/**/*.dart
find ./portal -name '*.dart' 2>/dev/null | while IFS= read -r f; do
    sedi "s|myservice|${SERVICE_NAME}|g" "$f"
done

ok "Service name updated."

# 5. Replace gist ID (if provided)
if [[ -n "$GIST_ID" ]]; then
    info "Replacing gist ID..."
    find ./.github/workflows -name '*.yml' | while IFS= read -r f; do
        sedi "s|9db84bf86640df151185b96b762b6c1e|${GIST_ID}|g" "$f"
    done
    ok "Gist ID updated."
fi

# ─── Cloud provider cleanup ─────────────────────────────────────────────────

if [[ "$CLOUD_PROVIDER" == "gcp" ]]; then
    info "Removing AWS-specific files..."
    rm -rf terraform/aws-dev/ terraform/aws-repo/ "terraform/modules/${SERVICE_NAME}-aws/"
    rm -f .github/workflows/publish-aws.yml .github/workflows/deploy-aws.yml .github/workflows/deploy-aws-dev.yml
    rm -rf internal/aws/
    ok "AWS files removed."
elif [[ "$CLOUD_PROVIDER" == "aws" ]]; then
    info "Removing GCP-specific files..."
    rm -rf terraform/dev/ terraform/repo/ "terraform/modules/${SERVICE_NAME}/"
    rm -f .github/workflows/publish.yml .github/workflows/deploy.yml .github/workflows/deploy-dev.yml .github/workflows/deploy-repo.yml
    rm -rf internal/google/
    ok "GCP files removed."
else
    ok "Keeping both cloud providers."
fi

# ─── README cleanup ─────────────────────────────────────────────────────────

info "Updating README..."

# Generate a title-cased version of the service name (hyphens to spaces, capitalize words)
TITLE_NAME=$(echo "$SERVICE_NAME" | tr '-' ' ' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) substr($i,2)}1')

# Replace title
sedi "s|^# My Service|# ${TITLE_NAME}|" README.md

# Replace the "Using This Template" section (lines starting from "## Using This Template"
# through "### 5. Configure deployment" paragraph end) with a short note
# We use awk for multi-line replacement
awk '
/^## Using This Template/ { skip=1; print "## Origin\n\nThis project was initialized from the [backend-template](https://github.com/alielgamal/backend-template).\n"; next }
/^## / && skip { skip=0 }
skip { next }
{ print }
' README.md > README.md.tmp && mv README.md.tmp README.md

ok "README updated."

# ─── Self-destruct ───────────────────────────────────────────────────────────

info "Cleaning up setup files..."

# Remove setup target from Makefile
awk '
/^\.PHONY: setup/ { skip=1; next }
/^setup:/ { skip=1; next }
skip && /^[^\t ]/ { skip=0 }
skip { next }
{ print }
' Makefile > Makefile.tmp && mv Makefile.tmp Makefile

# Remove blank lines that may have been left before help target
sedi '/^$/N;/^\n$/d' Makefile

# Delete this script
rm -f setup.sh

ok "Setup files removed."

# ─── Done ────────────────────────────────────────────────────────────────────

echo ""
printf "${GREEN}${BOLD}Setup complete!${NC}\n"
echo ""
echo "Next steps:"
echo "  1. Run 'go mod tidy' to update dependencies"
echo "  2. Review the changes and commit: git add -A && git commit -m 'Initialize project from template'"
echo "  3. Fill in TODO placeholders in terraform/ files"
echo "  4. Configure GitHub repository secrets and variables (see README)"
echo ""
