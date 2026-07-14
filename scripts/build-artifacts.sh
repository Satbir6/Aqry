#!/bin/sh

set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT_DIR=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)
DIST_DIR=$ROOT_DIR/dist
TEMP_DIR=$(mktemp -d "${TMPDIR:-/tmp}/aqry-build.XXXXXX")

cleanup() {
  rm -rf "$TEMP_DIR"
}

trap cleanup EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM

command -v go >/dev/null 2>&1 || {
  printf 'error: Go is required to build repository artifacts\n' >&2
  exit 1
}
command -v tar >/dev/null 2>&1 || {
  printf 'error: tar is required to package repository artifacts\n' >&2
  exit 1
}
command -v zip >/dev/null 2>&1 || {
  printf 'error: zip is required to package the Windows artifact\n' >&2
  exit 1
}

VERSION=${AQRY_VERSION:-repo-main}
if [ -n "${AQRY_COMMIT:-}" ]; then
  COMMIT=$AQRY_COMMIT
elif command -v git >/dev/null 2>&1; then
  COMMIT=$(git -C "$ROOT_DIR" rev-parse --short HEAD 2>/dev/null || printf unknown)
else
  COMMIT=unknown
fi
BUILD_DATE=${AQRY_BUILD_DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}
LDFLAGS="-s -w -X aqry/internal/version.Version=$VERSION -X aqry/internal/version.Commit=$COMMIT -X aqry/internal/version.Date=$BUILD_DATE"

mkdir -p "$DIST_DIR"
rm -f "$DIST_DIR"/aqry_*.tar.gz "$DIST_DIR"/aqry_*.zip "$DIST_DIR"/SHA256SUMS

build_tar() {
  TARGET_OS=$1
  TARGET_ARCH=$2
  ARCHIVE_NAME=$3
  BINARY_PATH=$TEMP_DIR/aqry

  printf 'Building %s/%s...\n' "$TARGET_OS" "$TARGET_ARCH"
  (
    cd "$ROOT_DIR"
    CGO_ENABLED=0 GOOS="$TARGET_OS" GOARCH="$TARGET_ARCH" \
      go build -trimpath -buildvcs=false -ldflags "$LDFLAGS" \
        -o "$BINARY_PATH" ./cmd/aqry
  )
  tar -czf "$DIST_DIR/$ARCHIVE_NAME" -C "$TEMP_DIR" aqry
  rm -f "$BINARY_PATH"
}

build_zip() {
  TARGET_OS=$1
  TARGET_ARCH=$2
  ARCHIVE_NAME=$3
  BINARY_PATH=$TEMP_DIR/aqry.exe

  printf 'Building %s/%s...\n' "$TARGET_OS" "$TARGET_ARCH"
  (
    cd "$ROOT_DIR"
    CGO_ENABLED=0 GOOS="$TARGET_OS" GOARCH="$TARGET_ARCH" \
      go build -trimpath -buildvcs=false -ldflags "$LDFLAGS" \
        -o "$BINARY_PATH" ./cmd/aqry
  )
  zip -q -j "$DIST_DIR/$ARCHIVE_NAME" "$BINARY_PATH"
  rm -f "$BINARY_PATH"
}

build_tar linux amd64 aqry_Linux_x86_64.tar.gz
build_tar linux arm64 aqry_Linux_arm64.tar.gz
build_tar darwin amd64 aqry_Darwin_x86_64.tar.gz
build_tar darwin arm64 aqry_Darwin_arm64.tar.gz
build_zip windows amd64 aqry_Windows_x86_64.zip

(
  cd "$DIST_DIR"
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum aqry_*.tar.gz aqry_*.zip > SHA256SUMS
  elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 aqry_*.tar.gz aqry_*.zip > SHA256SUMS
  else
    printf 'error: sha256sum or shasum is required to publish checksums\n' >&2
    exit 1
  fi
)

printf '\nArtifacts written to %s\n' "$DIST_DIR"
sed 's/^/  /' "$DIST_DIR/SHA256SUMS"
