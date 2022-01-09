#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="${REPO_ROOT:-"$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"}"
BUILD_ROOT="${REPO_ROOT}/build"
BUILD_BIN="${BUILD_ROOT}/bin"

# Will be set to false if any of the steps fail.
did_checks_pass=true

# Minimum versions.
MIN_GO_VERSION="1.16"
MIN_YARN_VERSION="1.22.0"

# param 1 - required version.
# param 2 - current version.
# returns true or false if the version is ok.
is_version_ok() {
  required_version=$1
  current_version=$2
  if [ "$(printf '%s\n' "$required_version" "$current_version" | sort -V | head -n1)" = "$required_version" ]; then
    return 0
  else
    return 1
  fi
}

os() {
  # If we're on OSX lets check for brew and coreutils as they are requirements.
  # docs/getting-started/local-build/#requirements.
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # check brew is installed
    if command -v brew --version &>/dev/null; then
      # check if coreutils is installed.
      if ! brew ls --versions coreutils >/dev/null; then
        echo "coreutils is not installed, this is a required dependency."
        echo "install by running [brew install coreutils]"
        did_checks_pass=false
      fi
    else
      echo "brew is not installed (visit https://brew.sh/ to install), unable to verify coreutils dependency."
      did_checks_pass=false
    fi
  fi
}

backend() {
  if ! command -v go -v &>/dev/null; then
    echo "golang is not installed or cannot be found in the current PATH, this is a required dependency."
    did_checks_pass=false
  else
    current_version=$(go version | {
      read -r _ _ v _
      echo "${v#go}"
    })
    if ! is_version_ok $MIN_GO_VERSION "$current_version"; then
      echo "golang version must be >= $MIN_GO_VERSION, current version $current_version"
      did_checks_pass=false
    fi
  fi

  YARN_BIN="${BUILD_BIN}/yarn.sh"
  if ! command -v "${YARN_BIN}" &>/dev/null; then
    echo "yarn is not installed or cannot be found in the current PATH, this is a required dependency."
    did_checks_pass=false
  else
    current_version=$("${YARN_BIN}" --version)
    if ! is_version_ok $MIN_YARN_VERSION "$current_version"; then
      echo "yarn version must be >= $MIN_YARN_VERSION, current version $current_version"
      did_checks_pass=false
    fi
  fi
}

main() {
  # always check OS level requirements
  os

  backend

  if [ "$did_checks_pass" = false ]; then
    printf "\nPlease refer to the development requirements /docs/getting-started/local-build/#requirements"
    return 1
  fi
}

main "$@"
