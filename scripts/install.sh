#!/bin/sh

set -eu

REPOSITORY=${AQRY_REPOSITORY:-Satbir6/Aqry}
REF=${AQRY_REF:-main}
MIN_GO_MAJOR=1
MIN_GO_MINOR=23
TEMP_DIR=
INSTALL_TMP=

if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
  COLOR_INFO=$(printf '\033[1;34m')
  COLOR_OK=$(printf '\033[1;32m')
  COLOR_WARN=$(printf '\033[1;33m')
  COLOR_ERROR=$(printf '\033[1;31m')
  COLOR_RESET=$(printf '\033[0m')
else
  COLOR_INFO=
  COLOR_OK=
  COLOR_WARN=
  COLOR_ERROR=
  COLOR_RESET=
fi

info() {
  printf '%s==>%s %s\n' "$COLOR_INFO" "$COLOR_RESET" "$*"
}

success() {
  printf '%s==>%s %s\n' "$COLOR_OK" "$COLOR_RESET" "$*"
}

warn() {
  printf '%s==>%s %s\n' "$COLOR_WARN" "$COLOR_RESET" "$*"
}

fail() {
  printf '%serror:%s %s\n' "$COLOR_ERROR" "$COLOR_RESET" "$*" >&2
  exit 1
}

cleanup() {
  if [ -n "$INSTALL_TMP" ]; then
    rm -f "$INSTALL_TMP"
  fi
  if [ -n "$TEMP_DIR" ]; then
    rm -rf "$TEMP_DIR"
  fi
}

trap cleanup EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM

require_command() {
  command -v "$1" >/dev/null 2>&1 || fail "required command not found: $1"
}

detect_platform() {
  SYSTEM_NAME=$(uname -s)
  case "$SYSTEM_NAME" in
    Linux)
      OS_NAME=linux
      ;;
    Darwin)
      OS_NAME=darwin
      ;;
    MINGW*|MSYS*|CYGWIN*|Windows_NT)
      fail "this installer supports Linux and macOS; build aqry from source with Go on Windows"
      ;;
    *)
      fail "unsupported operating system: $SYSTEM_NAME"
      ;;
  esac

  MACHINE_NAME=$(uname -m)
  case "$MACHINE_NAME" in
    x86_64|amd64)
      ARCH_NAME=amd64
      ;;
    arm64|aarch64)
      ARCH_NAME=arm64
      ;;
    *)
      fail "unsupported CPU architecture: $MACHINE_NAME (supported: amd64, arm64)"
      ;;
  esac
}

check_go_version() {
  command -v go >/dev/null 2>&1 || fail "Go 1.23 or newer is required because aqry is installed from source: https://go.dev/dl/"

  GO_VERSION=$(go env GOVERSION 2>/dev/null || true)
  GO_NUMBER=${GO_VERSION#go}
  GO_MAJOR=${GO_NUMBER%%.*}
  GO_REMAINDER=${GO_NUMBER#*.}
  GO_MINOR=${GO_REMAINDER%%.*}

  case "$GO_MAJOR:$GO_MINOR" in
    *[!0-9:]*)
      fail "could not parse the installed Go version: $GO_VERSION"
      ;;
  esac

  if [ "$GO_MAJOR" -lt "$MIN_GO_MAJOR" ] || { [ "$GO_MAJOR" -eq "$MIN_GO_MAJOR" ] && [ "$GO_MINOR" -lt "$MIN_GO_MINOR" ]; }; then
    fail "Go ${MIN_GO_MAJOR}.${MIN_GO_MINOR} or newer is required; found $GO_VERSION"
  fi
}

download_source() {
  ARCHIVE_PATH=$TEMP_DIR/aqry-source.tar.gz

  if [ -n "${AQRY_SOURCE_ARCHIVE:-}" ]; then
    [ -f "$AQRY_SOURCE_ARCHIVE" ] || fail "source archive not found: $AQRY_SOURCE_ARCHIVE"
    info "Using local source archive $AQRY_SOURCE_ARCHIVE"
    cp "$AQRY_SOURCE_ARCHIVE" "$ARCHIVE_PATH"
    return
  fi

  SOURCE_URL=${AQRY_SOURCE_URL:-https://codeload.github.com/${REPOSITORY}/tar.gz/${REF}}
  info "Downloading source from $REPOSITORY at $REF"

  if command -v curl >/dev/null 2>&1; then
    curl --fail --silent --show-error --location --retry 3 --proto '=https' --tlsv1.2 \
      "$SOURCE_URL" --output "$ARCHIVE_PATH"
  elif command -v wget >/dev/null 2>&1; then
    wget --quiet --output-document="$ARCHIVE_PATH" "$SOURCE_URL"
  else
    fail "curl or wget is required to download the aqry source archive"
  fi
}

extract_source() {
  SOURCE_DIR=$TEMP_DIR/source
  mkdir -p "$SOURCE_DIR"
  info "Extracting source archive"
  tar -xzf "$ARCHIVE_PATH" -C "$SOURCE_DIR" --strip-components=1

  [ -f "$SOURCE_DIR/go.mod" ] || fail "the downloaded archive does not contain aqry source code"
  [ -d "$SOURCE_DIR/cmd/aqry" ] || fail "the downloaded archive is missing cmd/aqry"
}

build_binary() {
  BINARY_PATH=$TEMP_DIR/aqry
  VERSION_LABEL=source-$REF
  info "Building aqry for $OS_NAME/$ARCH_NAME with $GO_VERSION"

  (
    cd "$SOURCE_DIR"
    CGO_ENABLED=0 GOOS="$OS_NAME" GOARCH="$ARCH_NAME" \
      go build -trimpath -buildvcs=false \
        -ldflags "-s -w -X aqry/internal/version.Version=${VERSION_LABEL}" \
        -o "$BINARY_PATH" ./cmd/aqry
  )

  [ -x "$BINARY_PATH" ] || fail "the aqry binary was not created"
}

choose_install_directory() {
  if [ -n "${AQRY_INSTALL_DIR:-}" ]; then
    INSTALL_DIR=$AQRY_INSTALL_DIR
  elif [ -d /usr/local/bin ] && [ -w /usr/local/bin ]; then
    INSTALL_DIR=/usr/local/bin
  else
    INSTALL_DIR=${HOME:?HOME is required}/.local/bin
  fi

  mkdir -p "$INSTALL_DIR"
  [ -d "$INSTALL_DIR" ] || fail "could not create installation directory: $INSTALL_DIR"
  [ -w "$INSTALL_DIR" ] || fail "installation directory is not writable: $INSTALL_DIR"
}

install_binary() {
  INSTALL_PATH=$INSTALL_DIR/aqry
  INSTALL_TMP=$INSTALL_DIR/.aqry-install-$$
  info "Installing aqry to $INSTALL_PATH"

  if command -v install >/dev/null 2>&1; then
    install -m 0755 "$BINARY_PATH" "$INSTALL_TMP"
  else
    cp "$BINARY_PATH" "$INSTALL_TMP"
    chmod 0755 "$INSTALL_TMP"
  fi

  mv -f "$INSTALL_TMP" "$INSTALL_PATH"
  INSTALL_TMP=
}

verify_installation() {
  info "Verifying installation"
  VERSION_OUTPUT=$("$INSTALL_PATH" --version 2>&1) || fail "installed binary could not be executed"
  success "$VERSION_OUTPUT"
  success "Installed at $INSTALL_PATH"
}

print_next_steps() {
  case ":${PATH:-}:" in
    *":$INSTALL_DIR:"*)
      RUN_COMMAND=aqry
      ;;
    *)
      RUN_COMMAND=$INSTALL_PATH
      warn "$INSTALL_DIR is not currently in PATH"
      printf '\nAdd it for this shell with:\n\n'
      printf '  export PATH="%s:$PATH"\n' "$INSTALL_DIR"
      ;;
  esac

  printf '\nTry aqry now:\n\n'
  printf '  "%s" example.com\n' "$RUN_COMMAND"
  printf '  "%s"\n\n' "$RUN_COMMAND"
}

main() {
  require_command uname
  require_command tar
  require_command mktemp
  require_command mkdir
  require_command rm
  require_command cp
  require_command mv
  require_command chmod

  detect_platform
  check_go_version
  info "Detected $OS_NAME/$ARCH_NAME"

  TEMP_DIR=$(mktemp -d "${TMPDIR:-/tmp}/aqry.XXXXXX")
  download_source
  extract_source
  build_binary
  choose_install_directory
  install_binary
  verify_installation
  print_next_steps
}

main "$@"
