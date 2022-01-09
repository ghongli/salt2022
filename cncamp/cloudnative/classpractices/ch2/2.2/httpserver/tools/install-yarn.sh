#!/usr/bin/env bash

set -euo pipefail

REPO_ROOT="${REPO_ROOT:-"$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"}"
BUILD_ROOT="${REPO_ROOT}/build"
BUILD_BIN="${BUILD_ROOT}/bin"

NAME=yarn
RELEASE=1.22.11

DEST_DIR="${BUILD_ROOT}/.yarn/releases"
DEST_FILE="${DEST_DIR}/${NAME}-${RELEASE}.js"
WRAPPER_DEST_DIR="${BUILD_BIN}"
WRAPPER_DEST_FILE="${WRAPPER_DEST_DIR}/${NAME}.sh"

if [[ ! -f "${DEST_FILE}" ]]; then
  echo "Downloading yarn v${RELEASE} to build environment..."
  mkdir -p "${DEST_DIR}"
  curl -sSL -o "${DEST_FILE}" \
    "https://github.com/yarnpkg/${NAME}/releases/download/v${RELEASE}/${NAME}-${RELEASE}.js"
fi

# Install a wrapper script in build/bin that executes yarn if it doesn't exist already.
WRAPPER_SCRIPT="#!/bin/bash\nnode \"${DEST_FILE}\" \"\$@\"\n"
if [[ ! -f "${WRAPPER_DEST_FILE}" || $(< "${WRAPPER_DEST_FILE}") != $(printf "%b" "${WRAPPER_SCRIPT}") ]]; then
  mkdir -p "${WRAPPER_DEST_DIR}"
  printf "%b" "${WRAPPER_SCRIPT}" > "${WRAPPER_DEST_FILE}"
  chmod +x "${WRAPPER_DEST_FILE}"
fi

echo "${NAME}-${RELEASE}.js OK"