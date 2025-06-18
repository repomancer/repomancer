#!/usr/bin/env bash
set -euo pipefail

# Check if current HEAD is tagged. If not, abort
VERSION_TAG=$(git tag --points-at HEAD)
if [[ "${VERSION_TAG}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  VERSION=${VERSION_TAG/#v/}
  ARTIFACT="Repomancer-${VERSION}-Darwin-arm64.dmg"
else
  echo "Current HEAD does not have tag matching ^v\d+.\d+.\d+$"
  echo "Aborting..."
  exit 1
fi

if [ ! -f "${ARTIFACT}" ]; then
  echo "${ARTIFACT} not found! Aborting"
  exit 1
fi

echo "Release version ${VERSION}"
# shellcheck disable=SC2034
read -rsp $'Press any key to continue...\n' -n1 key

# Create release and upload artifact
gh release create "${VERSION_TAG}"
gh release upload "${VERSION_TAG}" "${ARTIFACT}"

HASH=$(sha256 --quiet "${ARTIFACT}")
echo "${ARTIFACT} hash:"
echo "${HASH}"
echo Update release in https://github.com/repomancer/homebrew-repomancer
