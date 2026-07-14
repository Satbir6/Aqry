# aqry

`aqry` is a fast DNS lookup CLI with an interactive terminal interface. The default path stays script-friendly—one IPv4 address and no decoration—while the TUI makes record types and resolver choices easy to explore.

## Build

Requirements: Go 1.23 or newer.

```sh
go build -o bin/aqry ./cmd/aqry
./bin/aqry --version
```

You can also install the development build into your Go binary directory:

```sh
go install ./cmd/aqry
```

## Quick CLI usage

```sh
# First IPv4 address (the default)
aqry example.com

# Every matching address
aqry example.com --all

# Other record types
aqry example.com --type MX
aqry example.com -t AAAA

# Structured output
aqry example.com --json

# Query a specific resolver with a longer timeout
aqry example.com --server 1.1.1.1 --timeout 5
```

Plain output has no labels or ANSI styling. MX output includes its priority before the mail server. `--json` emits the shared structured result model.

## Interactive mode

Run without a domain to launch the Bubble Tea interface:

```sh
aqry
```

Prefill the domain with:

```sh
aqry --interactive example.com
```

The TUI adapts to compact, standard, and wide terminal widths. It includes phase-based loading progress, record and resolver pickers, a help modal, result navigation, lookup cancellation, and clipboard copy.

### TUI controls

| Key | Action |
| --- | --- |
| `Enter` | Run a lookup or confirm a picker |
| `Tab` / `Shift+Tab` | Move focus |
| `r` | Open the record picker when the domain input is not focused |
| `s` | Open the resolver picker when the domain input is not focused |
| `j` / `k`, `↑` / `↓` | Navigate picker items or results |
| `a` | Toggle the all-results view |
| `c` | Copy the selected result |
| `?` | Toggle help |
| `Esc` | Close a modal or cancel an active lookup |
| `q` | Quit when the domain input is not focused |
| `Ctrl+C` | Quit from anywhere |

Letter shortcuts are treated as text while the domain input is focused. Press `Tab` to move to the option areas before using them.

## Supported records

- `A` — IPv4 addresses
- `AAAA` — IPv6 addresses
- `CNAME` — canonical name
- `MX` — mail exchanges and priorities
- `TXT` — text records
- `NS` — name servers

`SOA`, `CAA`, and `SRV` deliberately return unsupported-record errors in this MVP.

## Resolver choices

The system DNS configuration is used by default. CLI and TUI modes also support:

- Cloudflare — `1.1.1.1`
- Google — `8.8.8.8`
- Quad9 — `9.9.9.9`
- A custom IPv4 or IPv6 resolver, optionally with a port

## Exit codes

| Code | Meaning |
| ---: | --- |
| `0` | Success |
| `1` | Invalid input |
| `2` | DNS lookup failure or no records |
| `3` | Timeout |
| `4` | Unsupported record type |
| `5` | Internal error |

## Test

```sh
go test ./...
go vet ./...
```

Most tests use injected resolvers and do not access the network. Enable the optional live smoke test with:

```sh
AQRY_LIVE_DNS_TEST=1 go test ./internal/dns -run TestLiveSystemResolver
```

## MVP limitations

- Go's standard resolver APIs do not expose accurate DNS TTL values, so TTL is omitted.
- `SOA`, `CAA`, DNSSEC validation, propagation comparison, history, and config files are deferred.
- Clipboard behavior depends on the host OS and terminal session providing a clipboard implementation.
