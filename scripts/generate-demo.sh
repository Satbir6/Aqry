#!/bin/sh

set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
cd "$ROOT_DIR"

command -v go >/dev/null 2>&1 || {
  printf 'error: go is required to build the demo binary\n' >&2
  exit 1
}
command -v vhs >/dev/null 2>&1 || {
  printf 'error: vhs is required to render the demo\n' >&2
  exit 1
}

mkdir -p bin
printf 'Building aqry for the recording...\n'
go build -o bin/aqry ./cmd/aqry

printf 'Rendering assets/aqry-demo.gif...\n'
vhs assets/aqry-demo.tape

printf 'Demo written to %s/assets/aqry-demo.gif\n' "$ROOT_DIR"
