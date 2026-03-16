#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PASS=0
FAIL=0

setup_repo() {
  REPO_DIR=$(mktemp -d)
  cd "$REPO_DIR"
  git init -q
  git config user.email "test@test.com"
  git config user.name "Test"
}

teardown_repo() {
  rm -rf "$REPO_DIR"
}

# Calls calculate-version.sh and parses its stdout into NEW_TAG and BUMP_PART.
calculate_version() {
  local default_bump="${1:-minor}"
  local v_prefix="${2:-true}"

  local output
  output=$("$SCRIPT_DIR/calculate-version.sh" "$default_bump" "$v_prefix" 2>/dev/null)

  NEW_TAG=""
  BUMP_PART=""
  while IFS='=' read -r key value; do
    case "$key" in
      new_tag) NEW_TAG="$value" ;;
      part) BUMP_PART="$value" ;;
    esac
  done <<< "$output"
}

assert_eq() {
  local test_name="$1" expected="$2" actual="$3"
  if [[ "$expected" == "$actual" ]]; then
    echo "  PASS: $test_name"
    PASS=$((PASS+1))
  else
    echo "  FAIL: $test_name (expected '$expected', got '$actual')"
    FAIL=$((FAIL+1))
  fi
}

# --- Tests ---

echo "Test: first release with no existing tags (default minor)"
setup_repo
git commit -q --allow-empty -m "initial commit"
calculate_version "minor" "true"
assert_eq "tag" "v0.1.0" "$NEW_TAG"
assert_eq "part" "minor" "$BUMP_PART"
teardown_repo

echo "Test: first release with default patch"
setup_repo
git commit -q --allow-empty -m "initial commit"
calculate_version "patch" "true"
assert_eq "tag" "v0.0.1" "$NEW_TAG"
assert_eq "part" "patch" "$BUMP_PART"
teardown_repo

echo "Test: minor bump from existing tag"
setup_repo
git commit -q --allow-empty -m "initial"
git tag v1.2.3
git commit -q --allow-empty -m "new feature"
calculate_version "minor" "true"
assert_eq "tag" "v1.3.0" "$NEW_TAG"
assert_eq "part" "minor" "$BUMP_PART"
teardown_repo

echo "Test: patch bump from existing tag"
setup_repo
git commit -q --allow-empty -m "initial"
git tag v1.2.3
git commit -q --allow-empty -m "bugfix"
calculate_version "patch" "true"
assert_eq "tag" "v1.2.4" "$NEW_TAG"
assert_eq "part" "patch" "$BUMP_PART"
teardown_repo

echo "Test: #major in commit message overrides default"
setup_repo
git commit -q --allow-empty -m "initial"
git tag v1.2.3
git commit -q --allow-empty -m "breaking change #major"
calculate_version "patch" "true"
assert_eq "tag" "v2.0.0" "$NEW_TAG"
assert_eq "part" "major" "$BUMP_PART"
teardown_repo

echo "Test: #minor in commit message overrides default"
setup_repo
git commit -q --allow-empty -m "initial"
git tag v1.2.3
git commit -q --allow-empty -m "new feature #minor"
calculate_version "patch" "true"
assert_eq "tag" "v1.3.0" "$NEW_TAG"
assert_eq "part" "minor" "$BUMP_PART"
teardown_repo

echo "Test: #patch in commit message overrides default"
setup_repo
git commit -q --allow-empty -m "initial"
git tag v1.2.3
git commit -q --allow-empty -m "fix #patch"
calculate_version "major" "true"
assert_eq "tag" "v1.2.4" "$NEW_TAG"
assert_eq "part" "patch" "$BUMP_PART"
teardown_repo

echo "Test: #none skips release"
setup_repo
git commit -q --allow-empty -m "initial"
git tag v1.2.3
git commit -q --allow-empty -m "docs only #none"
calculate_version "minor" "true"
assert_eq "tag" "" "$NEW_TAG"
assert_eq "part" "none" "$BUMP_PART"
teardown_repo

echo "Test: no v-prefix"
setup_repo
git commit -q --allow-empty -m "initial"
git tag 1.0.0
git commit -q --allow-empty -m "next"
calculate_version "minor" "false"
assert_eq "tag" "1.1.0" "$NEW_TAG"
assert_eq "part" "minor" "$BUMP_PART"
teardown_repo

echo "Test: picks latest tag when multiple exist"
setup_repo
git commit -q --allow-empty -m "initial"
git tag v1.0.0
git commit -q --allow-empty -m "second"
git tag v1.1.0
git commit -q --allow-empty -m "third"
git tag v2.0.0
git commit -q --allow-empty -m "fourth"
calculate_version "minor" "true"
assert_eq "tag" "v2.1.0" "$NEW_TAG"
teardown_repo

echo "Test: ignores floating tags (no patch segment)"
setup_repo
git commit -q --allow-empty -m "initial"
git tag v1.0.0
git tag v1
git commit -q --allow-empty -m "next"
calculate_version "minor" "true"
assert_eq "tag" "v1.1.0" "$NEW_TAG"
teardown_repo

echo "Test: keyword in commit body (not first line)"
setup_repo
git commit -q --allow-empty -m "initial"
git tag v1.2.3
git commit -q --allow-empty -m "some feature

This is a detailed description.
#major"
calculate_version "patch" "true"
assert_eq "tag" "v2.0.0" "$NEW_TAG"
assert_eq "part" "major" "$BUMP_PART"
teardown_repo

echo ""
echo "Results: $PASS passed, $FAIL failed"
if [[ $FAIL -gt 0 ]]; then
  exit 1
fi
