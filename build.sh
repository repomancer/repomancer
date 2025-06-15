#!/usr/bin/env bash
set -euxo pipefail

FULL_NAME=Repomancer
# If current commit is tagged with a version tag, set that as the version
VERSION="0.0.0"
VERSION_TAG=$(git tag --points-at HEAD)
if [[ "${VERSION_TAG}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]
then
  VERSION=${VERSION_TAG/#v}
fi

echo ${VERSION}
FILENAME="${FULL_NAME}-${VERSION}"

# Cleanup
rm -rf "${FULL_NAME}.app" "${FILENAME}-Darwin-arm64.zip" "${FILENAME}-Darwin-arm64.dmg"

# Test
go test ./...

# Package
fyne package -os darwin --app-build 1 --app-version "${VERSION}"
zip --symlinks -r "${FILENAME}-Darwin-arm64.zip" "${FULL_NAME}.app/"
hdiutil create -volname "${FULL_NAME}" -srcfolder "${FULL_NAME}.app" -ov -format UDZO "${FILENAME}-Darwin-arm64.dmg"
