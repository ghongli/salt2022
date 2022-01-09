#!/bin/bash

set -euo pipefail

REPO_ROOT="${REPO_ROOT:-"$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"}"
BUILD_ROOT="${REPO_ROOT}/build"
BUILD_BIN="${BUILD_ROOT}/bin"

NAME=shellcheck
RELEASE=v0.7.2

ARCH=x86_64

RELEASE_BINARY="${BUILD_BIN}/${NAME}-${RELEASE}"

ensure_binary() {
  if [[ ! -f "${RELEASE_BINARY}" ]]; then
    echo "info: Downloading ${NAME} ${RELEASE} to build environment"
    mkdir -p "${BUILD_BIN}"
  fi

  case "${OSTYPE}" in
    "darwin"*) os_type="darwin" ;;
    "linux"*) os_type="linux" ;;
    *) echo "error: Unsupported OS '${OSTYPE}' for shellcheck install, please install manually" && exit 1 ;;
  esac

  release_archive_dir="/tmp/${NAME}"
  mkdir -p "${release_archive_dir}"
  release_archive="${release_archive_dir}/${RELEASE}.tar.xz"
  URL="https://github.com/koalaman/${NAME}/releases/download/${RELEASE}/${NAME}-${RELEASE}.${os_type}.${ARCH}.tar.xz"
  echo "URL: ${URL}"
  curl -sSL -o "${release_archive}" "${URL}"
  tar -xvf "${release_archive}" -C "${release_archive_dir}"

  find "${BUILD_BIN}" -maxdepth 0 -regex '.*/'${NAME}'-[A-Za-z0-9\.]+$' -exec rm {} \;  # cleanup older versions
  mv "${release_archive_dir}/${NAME}-${RELEASE}/${NAME}" "${RELEASE_BINARY}"
  chmod +x "${RELEASE_BINARY}"

  # Cleanup stale resources.
  rm "${release_archive}"
  rm -rf "${release_archive_dir}"
}

main() {
  cd "${REPO_ROOT}"

  ensure_binary

  "${RELEASE_BINARY}" ./**/*.sh
}

main "$@"