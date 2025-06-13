#!/usr/bin/env bash
set -euxo pipefail

FULL_NAME=Repomancer
VERSION=$(date +%Y.%m%d.%H%M)

rm -rf ${FULL_NAME}.app

fyne package -os darwin --app-build 1 --app-version "${VERSION}"
zip --symlinks -r ${FULL_NAME}-${VERSION}-darwin-arm64.zip "${FULL_NAME}.app/"
