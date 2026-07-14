# Phased Plan For Building `aqry`

## Summary

Build `aqry` as a Go CLI/TUI DNS resolver from the two docs, starting with the script-friendly CLI path and then layering in the interactive Bubble Tea experience. The first complete release should support `A`, `AAAA`, `CNAME`, `MX`, `TXT`, and `NS` records, with `SOA`, `CAA`, history, config files, and propagation checks left for later.

## Phased Implementation

### Phase 1 — Project Bootstrap

- Create Go module using `module aqry` unless a remote repo path is later added.
- Add `cmd/aqry/main.go` and internal packages:
  - `internal/cli`
  - `internal/dns`
  - `internal/tui`
  - `internal/styles`
  - `internal/version`
- Use Go `1.23` per the technical design.
- Add dependencies:
  - `github.com/spf13/cobra`
  - `github.com/charmbracelet/bubbletea`
  - `github.com/charmbracelet/bubbles`
  - `github.com/charmbracelet/lipgloss`
  - Clipboard package only when copy support is added.

### Phase 2 — DNS Engine

- Implement `internal/dns` independently from CLI and TUI.
- Add:
  - `RecordType`
  - `LookupRequest`
  - `Record`
  - `LookupResult`
  - `Resolver` interface
  - `SystemResolver`
- Support `A`, `AAAA`, `CNAME`, `MX`, `TXT`, and `NS` using Go standard library DNS APIs.
- Normalize inputs like `https://google.com/search?q=x` into `google.com`.
- Validate empty domains, spaces, invalid characters, DNS length limits, and label length limits.
- Implement timeout handling with `context.Context`.
- Map errors to stable kinds:
  - invalid domain
  - no records
  - timeout
  - resolver failed
  - unsupported type
- Treat `SOA` and `CAA` as unsupported for this first implementation unless `miekg/dns` is added later.

### Phase 3 — Plain CLI Mode

- Implement Cobra root command.
- Mode selection:
  - `aqry <domain>` runs plain CLI lookup.
  - `aqry` opens TUI.
  - `aqry -i <domain>` opens TUI with the domain prefilled.
- Add flags:
  - `--type`, `-t`, default `A`
  - `--all`
  - `--json`
  - `--server`, `-s`
  - `--timeout`, default `3`
  - `--interactive`, `-i`
  - `--no-color`
  - `--version`, `-v`
- Plain output rules:
  - Default `aqry google.com` prints only the first IPv4.
  - `--all` prints one record value per line.
  - `--json` prints `LookupResult` as JSON.
  - No styling in default CLI output.
- Exit codes:
  - `0` success
  - `1` validation error
  - `2` DNS lookup failure
  - `3` timeout
  - `4` unsupported record type
  - `5` internal error

### Phase 4 — Bubble Tea TUI MVP

- Implement TUI model, messages, commands, update loop, and view rendering.
- Initial screen:
  - domain input focused
  - record type defaults to `A`
  - resolver defaults to system DNS
- Pressing `Enter` validates the domain and starts async lookup.
- Render states:
  - idle
  - loading
  - success
  - error
  - empty input
- Add spinner during lookup.
- Keep DNS lookup as a `tea.Cmd` so the UI does not block.
- Add base keyboard support:
  - `q` quit
  - `ctrl+c` force quit
  - `enter` lookup/confirm
  - `tab` and `shift+tab` cycle focus
  - `esc` closes modal or cancels transient UI state

### Phase 5 — Modern UI Layer

- Implement Lip Gloss theme tokens:
  - primary
  - muted
  - success
  - warning
  - danger
  - border
  - surface
  - text
- Build terminal layouts:
  - compact below 60 columns
  - standard from 60 to 100 columns
  - wide split layout above 100 columns
- Add:
  - header
  - main content area
  - footer help bar
  - focused panel styling
  - result panel
  - error panel
  - loading panel
- Support `--no-color` by disabling semantic color styling while preserving labels and symbols.

### Phase 6 — Interactive Features

- Add record type picker opened with `r`.
- Add resolver picker opened with `s`.
- Built-in resolver choices:
  - System DNS
  - Cloudflare `1.1.1.1`
  - Google `8.8.8.8`
  - Quad9 `9.9.9.9`
  - Custom
- Add help modal opened with `?`.
- Add result navigation:
  - `j` / `k`
  - up / down arrows
- Add `a` to toggle show-all results in TUI.
- Add `c` to copy selected result once clipboard dependency is included.
- Add phase-based progress bar:
  - validate domain
  - prepare resolver
  - query DNS
  - parse response
  - done
- Do not delay plain CLI output for animation.

### Phase 7 — Testing, Docs, And Release Prep

- Add unit tests for:
  - domain normalization
  - domain validation
  - record type parsing
  - unsupported records
  - output formatting
  - error-to-exit-code mapping
- Add TUI model tests for:
  - typing domain
  - opening/closing modals
  - starting lookup
  - success transition
  - error transition
  - record navigation
- Add mock resolver tests so most tests avoid live DNS.
- Add a small live DNS smoke test only if guarded/skippable.
- Add README with:
  - install/build instructions
  - CLI examples
  - TUI controls
  - supported records
  - known MVP limitations
- Add GitHub Actions later when the folder is an actual Git repo:
  - `go test ./...`
  - Linux build
  - Windows build

## Public Interfaces And Types

- `internal/dns.LookupRequest`:
  - `Domain string`
  - `RecordType dns.RecordType`
  - `Server string`
  - `Timeout time.Duration`
  - `All bool`
- `internal/dns.LookupResult`:
  - `Domain string`
  - `RecordType dns.RecordType`
  - `Resolver string`
  - `Records []dns.Record`
- `internal/dns.Record`:
  - `Type dns.RecordType`
  - `Name string`
  - `Value string`
  - `Priority int,omitempty`
  - `TTL int,omitempty`
- CLI contract:
  - default output remains plain and script-friendly
  - JSON output uses the same structured DNS result shape
  - unsupported record types return exit code `4`

## Test Plan

- Run `go test ./...` after each implementation phase.
- Manually verify:
  - `aqry google.com`
  - `aqry google.com --all`
  - `aqry google.com -t MX`
  - `aqry google.com --json`
  - `aqry google.com --server 1.1.1.1`
  - `aqry`
  - `aqry -i google.com`
- Manually check TUI at:
  - compact width below 60 columns
  - normal 80x24 terminal
  - wide layout above 100 columns
- Verify invalid domains fail cleanly without panic.

## Assumptions And Defaults

- The first implementation targets MVP plus the documented CLI flags.
- `SOA`, `CAA`, DNS TTL accuracy, propagation comparison, config files, and history are deferred.
- Go standard library DNS is used first; `miekg/dns` is introduced later only if accurate TTL/SOA/CAA support becomes required.
- Module path defaults to `aqry` because this folder is not currently a usable Git repository and no remote module path is discoverable.
