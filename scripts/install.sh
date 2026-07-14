#!/bin/sh

set -eu

REPOSITORY=${AQRY_REPOSITORY:-Satbir6/Aqry}
REF=${AQRY_REF:-main}
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
      OS_NAME=Linux
      ;;
    Darwin)
      OS_NAME=Darwin
      ;;
    MINGW*|MSYS*|CYGWIN*|Windows_NT)
      fail "use the Windows zip from https://github.com/$REPOSITORY/tree/$REF/dist"
      ;;
    *)
      fail "unsupported operating system: $SYSTEM_NAME"
      ;;
  esac

  MACHINE_NAME=$(uname -m)
  case "$MACHINE_NAME" in
    x86_64|amd64)
      ARCH_NAME=x86_64
      ;;
    arm64|aarch64)
      ARCH_NAME=arm64
      ;;
    *)
      fail "unsupported CPU architecture: $MACHINE_NAME (supported: amd64, arm64)"
      ;;
  esac

  ARCHIVE_NAME=aqry_${OS_NAME}_${ARCH_NAME}.tar.gz
}

download_file() {
  DOWNLOAD_URL=$1
  DOWNLOAD_PATH=$2

  if command -v curl >/dev/null 2>&1; then
    curl --fail --silent --show-error --location --retry 3 --proto '=https' --tlsv1.2 \
      "$DOWNLOAD_URL" --output "$DOWNLOAD_PATH"
  elif command -v wget >/dev/null 2>&1; then
    wget --quiet --output-document="$DOWNLOAD_PATH" "$DOWNLOAD_URL"
  else
    fail "curl or wget is required to download aqry"
  fi
}

acquire_artifact() {
  ARCHIVE_PATH=$TEMP_DIR/$ARCHIVE_NAME
  CHECKSUM_PATH=$TEMP_DIR/SHA256SUMS

  if [ -n "${AQRY_ARTIFACT_DIR:-}" ]; then
    [ -f "$AQRY_ARTIFACT_DIR/$ARCHIVE_NAME" ] || fail "artifact not found: $AQRY_ARTIFACT_DIR/$ARCHIVE_NAME"
    [ -f "$AQRY_ARTIFACT_DIR/SHA256SUMS" ] || fail "checksum file not found: $AQRY_ARTIFACT_DIR/SHA256SUMS"
    info "Using bundled artifact from $AQRY_ARTIFACT_DIR"
    cp "$AQRY_ARTIFACT_DIR/$ARCHIVE_NAME" "$ARCHIVE_PATH"
    cp "$AQRY_ARTIFACT_DIR/SHA256SUMS" "$CHECKSUM_PATH"
    return
  fi

  DOWNLOAD_BASE=${AQRY_DOWNLOAD_BASE:-https://raw.githubusercontent.com/${REPOSITORY}/${REF}/dist}
  info "Downloading $ARCHIVE_NAME from repository branch $REF"
  download_file "$DOWNLOAD_BASE/$ARCHIVE_NAME" "$ARCHIVE_PATH"
  download_file "$DOWNLOAD_BASE/SHA256SUMS" "$CHECKSUM_PATH"
}

file_checksum() {
  CHECKSUM_FILE=$1

  if command -v sha256sum >/dev/null 2>&1; then
    CHECKSUM_VALUE=$(sha256sum "$CHECKSUM_FILE")
    printf '%s\n' "${CHECKSUM_VALUE%% *}"
  elif command -v shasum >/dev/null 2>&1; then
    CHECKSUM_VALUE=$(shasum -a 256 "$CHECKSUM_FILE")
    printf '%s\n' "${CHECKSUM_VALUE%% *}"
  elif command -v openssl >/dev/null 2>&1; then
    CHECKSUM_VALUE=$(openssl dgst -sha256 "$CHECKSUM_FILE")
    printf '%s\n' "${CHECKSUM_VALUE##* }"
  else
    fail "sha256sum, shasum, or openssl is required to verify the download"
  fi
}

verify_artifact() {
  EXPECTED_CHECKSUM=
  while read -r CHECKSUM_VALUE CHECKSUM_FILE_NAME; do
    CHECKSUM_FILE_NAME=${CHECKSUM_FILE_NAME#\*}
    if [ "$CHECKSUM_FILE_NAME" = "$ARCHIVE_NAME" ]; then
      EXPECTED_CHECKSUM=$CHECKSUM_VALUE
      break
    fi
  done < "$CHECKSUM_PATH"

  [ -n "$EXPECTED_CHECKSUM" ] || fail "no checksum was published for $ARCHIVE_NAME"

  ACTUAL_CHECKSUM=$(file_checksum "$ARCHIVE_PATH")
  [ "$ACTUAL_CHECKSUM" = "$EXPECTED_CHECKSUM" ] || fail "checksum verification failed for $ARCHIVE_NAME"
  success "Checksum verified"
}

extract_binary() {
  EXTRACT_DIR=$TEMP_DIR/extracted
  mkdir -p "$EXTRACT_DIR"
  info "Extracting aqry binary"
  tar -xzf "$ARCHIVE_PATH" -C "$EXTRACT_DIR"
  BINARY_PATH=$EXTRACT_DIR/aqry
  [ -f "$BINARY_PATH" ] || fail "the downloaded archive does not contain aqry"
  chmod 0755 "$BINARY_PATH"
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
  info "Detected $OS_NAME/$ARCH_NAME"

  TEMP_DIR=$(mktemp -d "${TMPDIR:-/tmp}/aqry.XXXXXX")
  acquire_artifact
  verify_artifact
  extract_binary
  choose_install_directory
  install_binary
  verify_installation
  print_next_steps
}

main "$@"
