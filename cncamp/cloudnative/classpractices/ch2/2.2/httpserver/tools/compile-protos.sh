#!/usr/bin/env bash

set -euo pipefail

# https://github.com/protocolbuffers/protobuf/releases
PROTOC_RELEASE=3.17.3

# ideally matches https://github.com/bazelbuild/rules_go/blob/master/go/private/repositories.bzl (ultimately generated into grpc-go).
# However, these protos should be very stable, so drift is not a big concern.
GOOGLEAPIS_SHA=d4cd8d96ed6eb5dd7c997aab68a1d6bb0825090c

# https://github.com/angular/clang-format/releases
ANGULAR_CLANG_FORMAT_RELEASE=1.4.0

SCRIPT_ROOT="$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"

NAME=httpserver
REPO_ADDR=github.com/ghongli/salt2022/cncamp/cloudnative/classpractices/ch2/2.2/${NAME}

# .proto file and dir list.
PROTO_DOT_FILES=()
PROTO_DIRS=()

# parse options
ACTION="compile"
LINT_FIX=false
APP_API_ROOT=""
while getopts "lfc:" opt; do
  case $opt in
  l) ACTION="lint" ;;
  f) LINT_FIX=true ;;
  c) APP_API_ROOT="${OPTARG}" ;;
  *)
    echo "usage: $0 [-l]" >&2
    exit 1
    ;;
  esac
done
shift "$((OPTIND - 1))" # shift so that $@, $1, etc. refer to the non-option arguments

install_protoc() {
  export PROTOC_BIN="${GOBIN}/protoc-v${PROTOC_RELEASE}"
  export PROTOC_INCLUDE_DIR="${GOBIN}/protoc-v${PROTOC_RELEASE}-include"

  go install \
    github.com/bufbuild/buf/cmd/protoc-gen-buf-lint \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    github.com/envoyproxy/protoc-gen-validate \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

  if [[ ! -f "${PROTOC_BIN}" || ! -d "${PROTOC_INCLUDE_DIR}" ]]; then
    echo "info: Downloading protoc-v${PROTOC_RELEASE} to build environment"

    proto_arch=x86_64
    case "${OSTYPE}" in
    "darwin"*) proto_os="osx" ;;
    "linux"*) proto_os="linux" ;;
    *) echo "error: Unsupported OS '${OSTYPE}' for protoc install, please install manually" && exit 1 ;;
    esac

    proto_zip_out="/tmp/protoc-${PROTOC_RELEASE}.zip"
    curl -sSL -o "${proto_zip_out}" \
      "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_RELEASE}/protoc-${PROTOC_RELEASE}-${proto_os}-${proto_arch}.zip"

    proto_dir_out="/tmp/proto-${PROTOC_RELEASE}"
    mkdir -p "${proto_dir_out}"
    unzip -q -o "${proto_zip_out}" -d "${proto_dir_out}"

    if [[ ! -f ${PROTOC_BIN} ]]; then
      find "${GOBIN}" -maxdepth 0 -regex '.*/protoc-v[0-9\.]+$' -exec rm {} \; # cleanup older versions
      mv "${proto_dir_out}"/bin/protoc "${PROTOC_BIN}"
    fi

    if [[ ! -d "${PROTOC_INCLUDE_DIR}" ]]; then
      find "${GOBIN}" -maxdepth 0 -regex '.*/protoc-v.*?-include$' -type d -exec rm -r {} \; # cleanup older versions
      mv "${proto_dir_out}"/include "${PROTOC_INCLUDE_DIR}"
    fi

    # Cleanup stale resources.
    rm "${proto_zip_out}"
    rm -rf "${proto_dir_out}"
  fi
}

install_googleapis() {
  final_out_dir="${BUILD_ROOT}/bin/googleapis-${GOOGLEAPIS_SHA}/google"
  if [[ ! -f "${final_out_dir}/rpc/status.proto" ]]; then
    echo "info: Downloading googleapis@${GOOGLEAPIS_SHA} to build environment"
    googleapis_zip_out="/tmp/googleapis-${GOOGLEAPIS_SHA}.zip"
    curl -sSL -o "${googleapis_zip_out}" \
      "https://github.com/googleapis/googleapis/archive/${GOOGLEAPIS_SHA}.zip"

    googleapis_dir_out="/tmp/googleapis-${GOOGLEAPIS_SHA}"
    mkdir -p "${googleapis_dir_out}"
    unzip -q -o "${googleapis_zip_out}" -d "${googleapis_dir_out}"

    final_out_dir="${BUILD_ROOT}/bin/googleapis-${GOOGLEAPIS_SHA}/google"

    mkdir -p "${final_out_dir}/api" "${final_out_dir}/rpc"

    mv \
      "${googleapis_dir_out}/googleapis-${GOOGLEAPIS_SHA}/google/api/annotations.proto" \
      "${googleapis_dir_out}/googleapis-${GOOGLEAPIS_SHA}/google/api/field_behavior.proto" \
      "${googleapis_dir_out}/googleapis-${GOOGLEAPIS_SHA}/google/api/http.proto" \
      "${googleapis_dir_out}/googleapis-${GOOGLEAPIS_SHA}/google/api/httpbody.proto" \
      "${final_out_dir}/api"

    mv \
      "${googleapis_dir_out}/googleapis-${GOOGLEAPIS_SHA}/google/rpc/code.proto" \
      "${googleapis_dir_out}/googleapis-${GOOGLEAPIS_SHA}/google/rpc/error_details.proto" \
      "${googleapis_dir_out}/googleapis-${GOOGLEAPIS_SHA}/google/rpc/status.proto" \
      "${final_out_dir}/rpc"
  fi
}

install_clang_format() {
  CLANG_FORMAT_BIN="${GOBIN}/clang-format-${ANGULAR_CLANG_FORMAT_RELEASE}"
  if [[ ! -f "${CLANG_FORMAT_BIN}" ]]; then
    echo "info: Downloading clang-format to build environment"

    case "${OSTYPE}" in
    "darwin"*) clang_format_os="darwin" ;;
    "linux"*) clang_format_os="linux" ;;
    *) echo "error: Unsupported OS '${OSTYPE}' for clang-format install, please install manually" && exit 1 ;;
    esac

    clang_format_out="/tmp/clang-format"
    curl -sSL -o "${clang_format_out}" \
      "https://github.com/angular/clang-format/raw/v${ANGULAR_CLANG_FORMAT_RELEASE}/bin/${clang_format_os}_x64/clang-format"

    chmod a+x "${clang_format_out}"
    mv "${clang_format_out}" "${CLANG_FORMAT_BIN}"

    # Cleanup stale resources.
    rm -rf "${clang_format_out}"
  fi
}

prepare_build_environment() {
  export GOBIN="${BUILD_ROOT}/bin"
  mkdir -p "${GOBIN}"

  install_protoc
  install_googleapis

  if [[ "${ACTION}" == "lint" ]]; then
    install_clang_format
  fi
}

discover_protos() {
  while IFS= read -r -d '' proto; do
    PROTO_DOT_FILES+=("${proto}")
  done < <(find "${API_ROOT}" -name '*.proto' -print0 | sort -sdzu)

  while IFS= read -r -d '' proto_dirs; do
    PROTO_DIRS+=("${proto_dirs}")
  done < <(find "${API_ROOT}" -name '*.proto' -exec dirname {} \; | tr '\n' '\0' | sort -sdzu)

  while IFS= read -r -d '' proto; do
    PROTO_DOT_FILES+=("${proto}")
  done < <(find "${EXAMPLES_ROOT}" -name '*.proto' -print0 | sort -sdzu)

  while IFS= read -r -d '' proto_dirs; do
    PROTO_DIRS+=("${proto_dirs}")
  done < <(find "${EXAMPLES_ROOT}" -name '*.proto' -exec dirname {} \; | tr '\n' '\0' | sort -sdzu)
}

# get the go module stored dir and ensure that it's the correct version.
modpath() {
  set -e
  go mod download "${1}"
  go list -f "{{ .Dir }}" -m "${1}"
}

main() {
  REPO_ROOT="${SCRIPT_ROOT}"
  # Use alternate root if provided as command line argument.
  if [[ -n "${1-}" ]]; then
    REPO_ROOT="${1}"
  fi

  if [[ -z "${APP_API_ROOT}" ]]; then
    # if core is not provided, need to use a downloaded version.
    CORE_VERSION=$(cd "${REPO_ROOT}" && go list -f "{{ .Version }}" -m "${REPO_ADDR}")
    if [[ "${CORE_VERSION}" == *-*-* ]]; then
      # if a pseudo-version, figure out just the SHA
      CORE_VERSION=$(echo "${CORE_VERSION}" | awk -F"-" '{print $NF}')
    fi

    core_out="${REPO_ROOT}/build/bin/${CORE_VERSION}"
    if [[ ! -d "${core_out}" ]]; then
      echo "info: downloading core APIs ${CORE_VERSION} to build environment..."

      core_zip_out="/tmp/${CORE_VERSION}.tar.gz"
      core_tmp_out="/tmp/${CORE_VERSION}"
      wget -O "${core_zip_out}" \
        -q "https://${REPO_ADDR}/-/archive/${CORE_VERSION}/${NAME}-${CORE_VERSION}.zip"

      mkdir -p "${core_tmp_out}"
      unzip -q -o "${core_zip_out}" -d "${core_tmp_out}"

      mkdir -p "${core_out}"
      mv "${core_tmp_out}"/${NAME}-*/apis "${core_out}"

      # Cleanup stale resources.
      rm "${core_zip_out}"
      rm -rf "${core_tmp_out}"
    fi

    APP_API_ROOT="${core_out}/apis"
  fi

  API_ROOT="${REPO_ROOT}/apis"
  BUILD_ROOT="${REPO_ROOT}/build"
  EXAMPLES_ROOT="${REPO_ROOT}/examples/amiibo/apis"

  cd "${REPO_ROOT}/pkg"

  prepare_build_environment
  discover_protos

  googleapis_include_path="${BUILD_ROOT}/bin/googleapis-${GOOGLEAPIS_SHA}"
  pg_validate_include_path="$(modpath github.com/envoyproxy/protoc-gen-validate)"

  # Lint (fix) and exit if requested.
  if [[ "${ACTION}" == "lint" ]]; then
    cd "${API_ROOT}"

    # https://docs.buf.build/faq
    # https://docs.buf.build/tour-3
    buf_lint_config=$(cat "${API_ROOT}/buf.json")

    LINT_OK=true
    if [[ ${LINT_FIX} == true ]]; then
      for proto in "${PROTO_DOT_FILES[@]}"; do
        "${CLANG_FORMAT_BIN}" --style=file -i "${proto}"
      done
    else
      for proto in "${PROTO_DOT_FILES[@]}"; do
        if ! output=$("${CLANG_FORMAT_BIN}" --style=file "${proto}" | diff -u "${proto}" -); then
          echo "${output}"
          LINT_OK=false
        fi
      done

      for proto in "${PROTO_DOT_FILES[@]}"; do
        if ! output=$("${PROTOC_BIN}" \
          -I"${PROTOC_INCLUDE_DIR}" -I"${API_ROOT}" -I"${APP_API_ROOT}" \
          -I"${googleapis_include_path}" -I"${pg_validate_include_path}" \
          --buf-lint_out=. \
          "--buf-lint_opt={\"input_config\": ${buf_lint_config}}" \
          --plugin=protoc-gen-buf-lint="${GOBIN}/protoc-gen-buf-lint" \
          "${proto}" 2>&1); then
          printf "\nstart--- \n%s\n" "${proto}"
          echo "${output}"
          printf "\nsed --buf-check-lint_out:"
          echo "${output}" | sed 's/--buf-lint_out: //' | cut -d":" -f2-
          printf "end---\n"
          LINT_OK=false
        fi
      done
    fi
    ${LINT_OK} && exit 0 || exit 1
  fi

  # Compile.
  proto_out_dir="${REPO_ROOT}/pkg/api"
  mkdir -p "${proto_out_dir}"

  echo "info: compiling go(proto out dir: ${proto_out_dir})"
  for proto_dir in "${PROTO_DIRS[@]}"; do
    echo "${proto_dir}"
    "${PROTOC_BIN}" \
      -I"${PROTOC_INCLUDE_DIR}" -I"${API_ROOT}" -I"${EXAMPLES_ROOT}" -I"${APP_API_ROOT}" \
      -I"${googleapis_include_path}" -I"${pg_validate_include_path}" \
      --go_out "${proto_out_dir}" \
      --go_opt paths=source_relative \
      --go-grpc_out "${proto_out_dir}" \
      --go-grpc_opt require_unimplemented_servers=false,paths=source_relative \
      --validate_out paths=source_relative,lang=go:"${proto_out_dir}" \
      --grpc-gateway_out "${proto_out_dir}" \
      --grpc-gateway_opt warn_on_unbound_methods=true,paths=source_relative \
      --plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
      --plugin protoc-gen-go-grpc="${GOBIN}/protoc-gen-go-grpc" \
      --plugin protoc-gen-grpc-gateway="${GOBIN}/protoc-gen-grpc-gateway" \
      --plugin protoc-gen-validate="${GOBIN}/protoc-gen-validate" \
      "${proto_dir}"/*.proto
  done

  # Compile the test proto if we're in the core repository.
  if [[ "${SCRIPT_ROOT}" == "${REPO_ROOT}" ]]; then
    testpb="${REPO_ROOT}/pkg/internal/test/pb/test.proto"
    echo "${testpb}"
    "${PROTOC_BIN}" \
      -I"${PROTOC_INCLUDE_DIR}" -I"${API_ROOT}" -I"${googleapis_include_path}" \
      -I "${REPO_ROOT}"/pkg/internal/test/pb/ \
      --go_out "${REPO_ROOT}"/pkg/internal/test/pb \
      --go_opt paths=source_relative \
      --plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
      --plugin protoc-gen-go-grpc="${GOBIN}/protoc-gen-go-grpc" \
      --plugin protoc-gen-grpc-gateway="${GOBIN}/protoc-gen-grpc-gateway" \
      --plugin protoc-gen-validate="${GOBIN}/protoc-gen-validate" \
      "${testpb}"
  fi

  echo "compile-protos.sh OK"
}

main "$@"
