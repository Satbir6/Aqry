<div align="center">

# aqry

**Fast DNS answers in the shell. A focused DNS workspace when you need more.**

[![Test and build](https://github.com/Satbir6/Aqry/actions/workflows/test.yml/badge.svg)](https://github.com/Satbir6/Aqry/actions/workflows/test.yml)
[![Go 1.23+](https://img.shields.io/badge/Go-1.23%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Platforms](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-7C8CFF)](#platform-support)

</div>

`aqry` is a modern DNS resolver for the terminal. Use it as a clean, script-friendly CLI for quick answers, or open the interactive TUI to explore record types and DNS resolvers without memorizing flags.

```console
$ aqry example.com
104.20.34.220
```

```text
╭──────────────────────────────────────────────────────────────────────╮
│ aqry                                                    DNS Resolver │
│                                                                      │
│ › Domain                         ╭─ Results ───────────────────────╮ │
│   example.com                    │ ✓ Resolved successfully         │ │
│                                  │                                  │ │
│   Record Type                    │ 104.20.34.220                    │ │
│   [ A ] AAAA CNAME MX TXT NS     │                                  │ │
│                                  ╰──────────────────────────────────╯ │
│   Resolver · System DNS                                             │
│                                                                      │
│ enter lookup · tab navigate · ? help · ctrl+c quit                  │
╰──────────────────────────────────────────────────────────────────────╯
```

## Why aqry?

- **Fast by default** — `aqry domain.com` prints one IPv4 address and nothing else.
- **Interactive when useful** — launch the keyboard-driven TUI with `aqry`.
- **Multiple record types** — query A, AAAA, CNAME, MX, TXT, and NS records.
- **Resolver control** — use system DNS, Cloudflare, Google, Quad9, or a custom server.
- **Automation friendly** — plain output, JSON output, stable errors, and documented exit codes.
- **Responsive terminal UI** — compact, standard, and wide layouts with clear focus states.
- **Cross-platform Go binary** — designed for Linux, macOS, and Windows.

## Install

### One-command installation on Linux and macOS

The installer builds `aqry` directly from the repository source. **It does not download or use GitHub Release artifacts.**

```sh
curl -fsSL https://raw.githubusercontent.com/Satbir6/Aqry/main/scripts/install.sh | sh
```

The installer:

1. Detects Linux or macOS.
2. Detects `amd64` or `arm64`.
3. Downloads and extracts the repository source archive.
4. Builds the native binary with Go.
5. Installs it into a user-writable binary directory.
6. Verifies the installation with `aqry --version`.
7. Prints the installed path and the next command to run.

Go 1.23 or newer is required because installation happens from source. The installer also needs `curl` or `wget`, plus `tar`.

#### Arch Linux

```sh
sudo pacman -S --needed go
curl -fsSL https://raw.githubusercontent.com/Satbir6/Aqry/main/scripts/install.sh | sh
```

For clipboard support in the TUI, X11 sessions may also need `xclip`:

```sh
sudo pacman -S --needed xclip
```

#### Review the installer before running it

```sh
curl -fsSL https://raw.githubusercontent.com/Satbir6/Aqry/main/scripts/install.sh -o install-aqry.sh
less install-aqry.sh
sh install-aqry.sh
```

#### Choose another installation directory

```sh
curl -fsSL https://raw.githubusercontent.com/Satbir6/Aqry/main/scripts/install.sh | \
  AQRY_INSTALL_DIR="$HOME/bin" sh
```

By default, the installer uses `/usr/local/bin` only when that directory is already writable. Otherwise it uses `$HOME/.local/bin`; it never invokes `sudo`.

To install another branch, tag, or commit:

```sh
curl -fsSL https://raw.githubusercontent.com/Satbir6/Aqry/main/scripts/install.sh | \
  AQRY_REF=main sh
```

### Build manually

```sh
git clone https://github.com/Satbir6/Aqry.git
cd Aqry
go test ./...
mkdir -p bin
go build -o bin/aqry ./cmd/aqry
./bin/aqry --version
```

To install the binary for your user:

```sh
mkdir -p "$HOME/.local/bin"
install -m 0755 bin/aqry "$HOME/.local/bin/aqry"
```

## Usage

### Quick lookup

```sh
aqry example.com
```

The default output is the first IPv4 address with no labels or ANSI styling.

### Common commands

```sh
# Every A record
aqry example.com --all

# IPv6
aqry example.com --type AAAA

# Mail exchanges
aqry example.com -t MX --all

# Structured output
aqry example.com --json

# Query Cloudflare with a five-second timeout
aqry example.com --server 1.1.1.1 --timeout 5

# Open the TUI
aqry

# Open the TUI with a prefilled domain
aqry --interactive example.com
```

### Flags

| Flag | Short | Default | Description |
| --- | ---: | --- | --- |
| `--type` | `-t` | `A` | Record type: A, AAAA, CNAME, MX, TXT, or NS |
| `--all` | | `false` | Print every matching record |
| `--json` | | `false` | Print structured JSON |
| `--server` | `-s` | `system` | DNS resolver IP address |
| `--timeout` | | `3` | Lookup timeout in seconds |
| `--interactive` | `-i` | `false` | Force interactive mode |
| `--no-color` | | `false` | Disable semantic TUI colors |
| `--version` | `-v` | | Print the installed version |
| `--help` | `-h` | | Show command help |

### JSON output

```json
{
  "domain": "example.com",
  "type": "A",
  "resolver": "system",
  "records": [
    {
      "type": "A",
      "name": "example.com",
      "value": "104.20.34.220"
    }
  ]
}
```

## Interactive TUI

Run `aqry` without a domain. Focus starts in the domain input.

`Tab` cycles through:

```text
Domain → Record Type → Resolver → Results
```

The focused section has a `›` marker and an accent highlight. On Record Type or Resolver, press `Enter` or an arrow key to open its picker.

### Keyboard controls

| Key | Action |
| --- | --- |
| `Enter` | Run a lookup, open a focused picker, or confirm a choice |
| `Tab` / `Shift+Tab` | Move focus forward or backward |
| `r` | Open the record picker outside the domain input |
| `s` | Open the resolver picker outside the domain input |
| `j` / `k`, `↑` / `↓` | Navigate picker items or DNS results |
| `a` | Toggle the all-results view |
| `c` | Copy the selected result |
| `?` | Toggle help |
| `Esc` | Close a modal or cancel an active lookup |
| `q` | Quit outside the domain input |
| `Ctrl+C` | Quit from anywhere |

Letter shortcuts are treated as normal text while the domain input is focused. Press `Tab` before using `r`, `s`, `a`, `c`, or `q` as shortcuts.

## DNS support

### Record types

| Type | Result |
| --- | --- |
| `A` | IPv4 addresses |
| `AAAA` | IPv6 addresses |
| `CNAME` | Canonical name |
| `MX` | Mail servers with priority |
| `TXT` | Text records |
| `NS` | Authoritative name servers |

`SOA`, `CAA`, and `SRV` are intentionally deferred in the current version.

### Resolver presets

| Resolver | Address |
| --- | --- |
| System DNS | Operating-system default |
| Cloudflare | `1.1.1.1` |
| Google | `8.8.8.8` |
| Quad9 | `9.9.9.9` |
| Custom | IPv4 or IPv6, optionally with a port |

## Exit codes

| Code | Meaning |
| ---: | --- |
| `0` | Success |
| `1` | Invalid input |
| `2` | DNS lookup failed or no records were found |
| `3` | Lookup timed out |
| `4` | Unsupported record type |
| `5` | Internal error |

## Platform support

| Platform | Status |
| --- | --- |
| Linux amd64 / arm64 | Supported |
| macOS amd64 / arm64 | Supported by the source installer |
| Windows amd64 | Cross-build verified; build from source with Go |

Windows PowerShell build:

```powershell
git clone https://github.com/Satbir6/Aqry.git
Set-Location Aqry
go build -o aqry.exe ./cmd/aqry
./aqry.exe --version
```

## Development

```sh
go test ./...
go test -race ./...
go vet ./...
```

Most tests use injected resolvers and do not access the network. Enable the optional live smoke test with:

```sh
AQRY_LIVE_DNS_TEST=1 go test ./internal/dns -run TestLiveSystemResolver
```

## Current limitations

- Go's standard resolver APIs do not expose accurate TTL values, so TTL is omitted.
- SOA, CAA, DNSSEC validation, propagation comparison, history, and config files are not implemented yet.
- Clipboard integration depends on clipboard support being available in the host session.
- The one-command installer builds from source and therefore requires Go 1.23+.

## Uninstall

If installed into the default user directory:

```sh
rm "$HOME/.local/bin/aqry"
```

If the installer selected `/usr/local/bin`, remove `/usr/local/bin/aqry` instead.

## Contributing

Bug reports and focused pull requests are welcome. Before opening a pull request, run:

```sh
go test ./...
go vet ./...
```

<div align="center">

Built with [Go](https://go.dev/), [Cobra](https://github.com/spf13/cobra), [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles), and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

</div>
