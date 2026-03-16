#!/usr/bin/env bash
set -euo pipefail

# Calculate the next semantic version based on git tags and commit message.
# Usage: calculate-version.sh <default-bump> <v-prefix>
#
# Outputs (one per line): new_tag, major, minor, part
# If GITHUB_OUTPUT is set, also writes to it for GitHub Actions.

DEFAULT_BUMP="${1:-minor}"
V_PREFIX="${2:-true}"

PREFIX=""
if [[ "$V_PREFIX" == "true" ]]; then
  PREFIX="v"
fi

LATEST_TAG=$(git tag --list "${PREFIX}[0-9]*.[0-9]*.[0-9]*" --sort=-v:refname | head -n 1)
if [[ -z "$LATEST_TAG" ]]; then
  LATEST_TAG="${PREFIX}0.0.0"
fi

VERSION="${LATEST_TAG#v}"
IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION"

COMMIT_MSG=$(git log -1 --format=%B)
if echo "$COMMIT_MSG" | grep -qi '#major'; then
  BUMP="major"
elif echo "$COMMIT_MSG" | grep -qi '#minor'; then
  BUMP="minor"
elif echo "$COMMIT_MSG" | grep -qi '#patch'; then
  BUMP="patch"
elif echo "$COMMIT_MSG" | grep -qi '#none'; then
  BUMP="none"
else
  BUMP="$DEFAULT_BUMP"
fi

if [[ "$BUMP" == "none" ]]; then
  echo "part=none"
  if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
    echo "part=none" >> "$GITHUB_OUTPUT"
  fi
  exit 0
fi

case "$BUMP" in
  major) MAJOR=$((MAJOR+1)); MINOR=0; PATCH=0 ;;
  minor) MINOR=$((MINOR+1)); PATCH=0 ;;
  patch) PATCH=$((PATCH+1)) ;;
esac

NEW_TAG="${PREFIX}${MAJOR}.${MINOR}.${PATCH}"

echo "new_tag=$NEW_TAG"
echo "major=${PREFIX}${MAJOR}"
echo "minor=${PREFIX}${MAJOR}.${MINOR}"
echo "part=$BUMP"

if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
  echo "new_tag=$NEW_TAG" >> "$GITHUB_OUTPUT"
  echo "major=${PREFIX}${MAJOR}" >> "$GITHUB_OUTPUT"
  echo "minor=${PREFIX}${MAJOR}.${MINOR}" >> "$GITHUB_OUTPUT"
  echo "part=$BUMP" >> "$GITHUB_OUTPUT"
fi

echo "Creating release $NEW_TAG (${BUMP} bump from ${LATEST_TAG})" >&2
