<div align="center">

# aqry

**Fast DNS answers in the shell. A focused DNS workspace when you need more.**

[![Test and build](https://github.com/Satbir6/Aqry/actions/workflows/test.yml/badge.svg)](https://github.com/Satbir6/Aqry/actions/workflows/test.yml)
[![Go 1.23+](https://img.shields.io/badge/Go-1.23%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Platforms](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-7C8CFF)](#platform-support)
[![Prebuilt binaries](https://img.shields.io/badge/prebuilt%20binaries-repository-63D297)](#install)
[![License: MIT](https://img.shields.io/badge/License-MIT-63D297.svg)](./LICENSE)
[![Last commit](https://img.shields.io/github/last-commit/Satbir6/Aqry?color=7C8CFF)](https://github.com/Satbir6/Aqry/commits/main)
[![Repository size](https://img.shields.io/github/repo-size/Satbir6/Aqry?color=63D297)](https://github.com/Satbir6/Aqry)

</div>

`aqry` is a fast DNS resolver for the terminal. It gives clean one-line answers for scripts and a polished keyboard-driven TUI when you want to inspect records interactively.

> Status: Early public release. Core DNS lookup and TUI features are usable, but packaging and advanced DNS features are still evolving.

## Contents

- [Why aqry?](#why-aqry)
- [Install](#install)
- [Usage](#usage)
- [Interactive TUI](#interactive-tui)
- [DNS support](#dns-support)
- [Platform support](#platform-support)
- [Development](#development)
- [Roadmap](#roadmap)
- [Uninstall](#uninstall)

## Preview

<p align="center">
  <img src="assets/aqry-demo.gif" alt="Animated terminal demo of aqry command-line lookup and interactive DNS resolver" width="900">
</p>

<p align="center"><em>Fast answers for scripts, with a keyboard-driven workspace for deeper inspection.</em></p>

## Quick start

Install a prebuilt binary on Linux or macOS—Go is not required:

```sh
curl -fsSL https://raw.githubusercontent.com/Satbir6/Aqry/main/scripts/install.sh | sh
```

Then run a quick lookup or open the interactive workspace:

```sh
aqry example.com
aqry
```

The installer prints the exact installed path and command to use if the destination is not already in your `PATH`.

## Why aqry?

Traditional DNS tools are excellent at deep diagnostics, but their output can be noisy for everyday checks. `aqry` keeps the common path short: pass a domain to get a script-friendly answer, or run it without arguments for a focused interactive workspace.

## Features

| Feature | Status |
| --- | --- |
| Fast IPv4 lookup | Supported |
| Interactive TUI | Supported |
| A / AAAA / CNAME / MX / TXT / NS | Supported |
| Resolver picker | Supported |
| JSON output | Supported |
| Clipboard copy | Supported |
| SOA / CAA / SRV | Planned |
| DNSSEC validation | Planned |
| GitHub Release binaries | Planned |

## How aqry compares

| Tool | Best for |
| --- | --- |
| `dig` | Advanced DNS debugging |
| `nslookup` | Basic built-in DNS checks |
| `host` | Simple DNS lookups |
| `aqry` | Fast answers plus an interactive DNS workspace |

## Install

| Platform | Recommended install | Requires Go |
| --- | --- | --- |
| Linux | `curl ... | sh` | No |
| macOS | `curl ... | sh` | No |
| Windows | `irm ... | iex` | No |
| From source | `go build` | Yes |

### Requirements

| Platform | Required tools |
| --- | --- |
| Linux/macOS | `curl` or `wget`, `tar`, checksum tool |
| Windows | PowerShell 5.1+ or PowerShell 7 |
| Source build | Go 1.23+ |

### Linux and macOS

The installer downloads a matching prebuilt `aqry` binary stored directly in this repository. It detects the operating system and CPU architecture, verifies the published SHA-256 checksum, extracts the binary, and installs it into a user-writable directory.

```sh
curl -fsSL https://raw.githubusercontent.com/Satbir6/Aqry/main/scripts/install.sh | sh
```

The installer supports Linux and macOS on `amd64` and `arm64`. It needs `curl` or `wget`, `tar`, and one of `sha256sum`, `shasum`, or `openssl`. It never invokes `sudo`, and Go is not required.

By default, it uses `/usr/local/bin` when that directory is already writable. Otherwise it installs to `$HOME/.local/bin`.

#### Arch Linux

Use the same one-command installer:

```sh
curl -fsSL https://raw.githubusercontent.com/Satbir6/Aqry/main/scripts/install.sh | sh
```

For clipboard support in an X11 TUI session, install `xclip`:

```sh
sudo pacman -S --needed xclip
```

#### Review the installer first

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

To install from another branch, tag, or commit containing matching files under `dist/`:

```sh
curl -fsSL https://raw.githubusercontent.com/Satbir6/Aqry/main/scripts/install.sh | \
  AQRY_REF=main sh
```

### Windows

Run this command in Windows PowerShell 5.1 or PowerShell 7:

```powershell
irm https://raw.githubusercontent.com/Satbir6/Aqry/main/scripts/install.ps1 | iex
```

The installer detects Windows `amd64`, downloads the committed zip and checksum file, verifies SHA-256, and installs `aqry.exe` to `%LOCALAPPDATA%\Programs\aqry`. It adds that directory to your user `PATH`, verifies the binary with `aqry.exe --version`, and prints the next commands to run. Administrator access and Go are not required.

To review the installer before running it:

```powershell
irm https://raw.githubusercontent.com/Satbir6/Aqry/main/scripts/install.ps1 -OutFile install-aqry.ps1
Get-Content .\install-aqry.ps1
& .\install-aqry.ps1
```

Choose another installation directory with:

```powershell
$env:AQRY_INSTALL_DIR = "$HOME\bin"
irm https://raw.githubusercontent.com/Satbir6/Aqry/main/scripts/install.ps1 | iex
```

### Direct binary downloads

These archives are committed under `dist/`; they are not GitHub Release assets.

| Platform | Architecture | Download |
| --- | --- | --- |
| Linux | x86_64 / amd64 | [aqry_Linux_x86_64.tar.gz](https://raw.githubusercontent.com/Satbir6/Aqry/main/dist/aqry_Linux_x86_64.tar.gz) |
| Linux | arm64 | [aqry_Linux_arm64.tar.gz](https://raw.githubusercontent.com/Satbir6/Aqry/main/dist/aqry_Linux_arm64.tar.gz) |
| macOS | Intel / x86_64 | [aqry_Darwin_x86_64.tar.gz](https://raw.githubusercontent.com/Satbir6/Aqry/main/dist/aqry_Darwin_x86_64.tar.gz) |
| macOS | Apple Silicon / arm64 | [aqry_Darwin_arm64.tar.gz](https://raw.githubusercontent.com/Satbir6/Aqry/main/dist/aqry_Darwin_arm64.tar.gz) |
| Windows | x86_64 / amd64 | [aqry_Windows_x86_64.zip](https://raw.githubusercontent.com/Satbir6/Aqry/main/dist/aqry_Windows_x86_64.zip) |

Published hashes are available in [SHA256SUMS](https://raw.githubusercontent.com/Satbir6/Aqry/main/dist/SHA256SUMS).

### Build from source

Building from source is optional and intended for contributors. It requires Go 1.23 or newer.

```sh
git clone https://github.com/Satbir6/Aqry.git
cd Aqry
go test ./...
mkdir -p bin
go build -o bin/aqry ./cmd/aqry
./bin/aqry --version
```

To install the source-built binary for your user:

```sh
mkdir -p "$HOME/.local/bin"
install -m 0755 bin/aqry "$HOME/.local/bin/aqry"
```

## Usage

Run `aqry` with a domain for plain output. By default it prints the first IPv4 address with no labels or ANSI styling.

```sh
aqry example.com
```

## Examples

| Task | Command |
| --- | --- |
| First IPv4 address | `aqry example.com` |
| All IPv4 addresses | `aqry example.com --all` |
| IPv6 records | `aqry example.com -t AAAA` |
| Mail records | `aqry example.com -t MX --all` |
| TXT records | `aqry example.com -t TXT --all` |
| Name servers | `aqry example.com -t NS --all` |
| JSON output | `aqry example.com --json` |
| Use Cloudflare DNS | `aqry example.com -s 1.1.1.1` |
| Open TUI | `aqry` |

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

Run `aqry` without a domain. Focus starts in the domain input, and `Tab` cycles through the workspace:

```text
Domain → Record Type → Resolver → Results
```

The focused section has a `›` marker and accent border. On Record Type or Resolver, press `Enter` or an arrow key to open its picker.

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

`SOA`, `CAA`, and `SRV` are planned but are not implemented in the current version.

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

| Platform | Installation status |
| --- | --- |
| Linux amd64 / arm64 | Supported by the repository binary installer |
| macOS amd64 / arm64 | Supported by the repository binary installer |
| Windows amd64 | Supported by the repository binary installer |

### Windows

Install the prebuilt binary without Go from PowerShell:

```powershell
irm https://raw.githubusercontent.com/Satbir6/Aqry/main/scripts/install.ps1 | iex
```

The installer adds aqry to your user `PATH`. Open a new terminal afterward if the current session does not recognize `aqry` immediately.

## Roadmap

### Core

- [x] Fast A-record lookup
- [x] JSON output
- [x] Resolver picker
- [x] Interactive Bubble Tea TUI

### Packaging

- [x] Linux/macOS one-command installer
- [x] Windows one-command installer
- [ ] GitHub Release binaries
- [ ] Homebrew tap
- [ ] AUR package
- [ ] Scoop package

### DNS features

- [ ] SOA, CAA, and SRV records
- [ ] DNS propagation comparison
- [ ] DNSSEC validation
- [ ] Config file support

## Development

Run the test suite and static checks:

```sh
go test ./...
go test -race ./...
go vet ./...
```

Maintainers can regenerate every committed platform archive and `SHA256SUMS` with:

```sh
scripts/build-artifacts.sh
```

The artifact builder requires Go, `tar`, and `zip`. End users installing a prebuilt binary do not need those build tools.

Most tests use injected resolvers and do not access the network. Enable the optional live smoke test with:

```sh
AQRY_LIVE_DNS_TEST=1 go test ./internal/dns -run TestLiveSystemResolver
```

Regenerate the README demo with [VHS](https://github.com/charmbracelet/vhs):

```sh
scripts/generate-demo.sh
```

The reproducible recording source is stored in `assets/aqry-demo.tape`. Rendering it performs live DNS lookups.

## Current limitations

- Go's standard resolver APIs do not expose accurate TTL values, so TTL is omitted.
- SOA, CAA, SRV, DNSSEC validation, propagation comparison, history, and config files are not implemented yet.
- Clipboard integration depends on clipboard support being available in the host session.
- Prebuilt binaries currently live in the repository rather than GitHub Releases.

## Uninstall

### Linux and macOS

If installed into the default user directory:

```sh
rm "$HOME/.local/bin/aqry"
```

If the installer selected `/usr/local/bin`, remove `/usr/local/bin/aqry` instead. For a custom location, remove the path printed by the installer.

### Windows

Remove the installation directory and its user `PATH` entry from PowerShell:

```powershell
$installDir = "$env:LOCALAPPDATA\Programs\aqry"
Remove-Item $installDir -Recurse -Force

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$updatedPath = (($userPath -split ";") | Where-Object {
  $_ -and -not [string]::Equals(
    $_.TrimEnd([IO.Path]::DirectorySeparatorChar),
    $installDir.TrimEnd([IO.Path]::DirectorySeparatorChar),
    [StringComparison]::OrdinalIgnoreCase
  )
}) -join ";"
[Environment]::SetEnvironmentVariable("Path", $updatedPath, "User")
```

If you selected a custom installation directory, change `$installDir` before running these commands.

## Contributing

Bug reports, documentation improvements, and focused pull requests are welcome. Before opening a pull request, run:

```sh
go test ./...
go vet ./...
```

## License

MIT License. See [LICENSE](./LICENSE).

## Built with

`aqry` is built with [Go](https://go.dev/), [Cobra](https://github.com/spf13/cobra), [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles), and [Lip Gloss](https://github.com/charmbracelet/lipgloss). The README demo is recorded with [VHS](https://github.com/charmbracelet/vhs).
