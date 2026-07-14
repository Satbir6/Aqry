# TECHNICAL_DESIGN.md — `aqry`

## 1. Overview

`aqry` is a cross-platform DNS resolver CLI and interactive terminal UI built in **Go**.

The product has two primary modes:

1. **Fast CLI mode**

```bash
# Folder: .
# File: terminal

aqry google.com
```

Output:

```text
# Folder: .
# File: terminal-output

142.250.182.46
```

2. **Interactive TUI mode**

```bash
# Folder: .
# File: terminal

aqry
```

This launches a modern Bubble Tea interface where the user can:

- Enter a domain
- Resolve IPv4 / A records by default
- Select DNS record types without typing flags
- Choose a DNS resolver
- View results in a polished layout
- Navigate using keyboard shortcuts
- See loading spinners, progress bars, and help menus

---

## 2. Technical Goals

| Goal                   | Requirement                                                   |
| ---------------------- | ------------------------------------------------------------- |
| Cross-platform         | Linux, Windows, macOS-compatible binary                       |
| Fast default mode      | `aqry google.com` should return quickly with minimal overhead |
| Interactive mode       | Bubble Tea-based full-screen TUI                              |
| Modern UI              | Lip Gloss styling, panels, focus states, footer help          |
| DNS support            | A, AAAA, CNAME, MX, TXT, NS, SOA                              |
| Keyboard-first         | No mouse required                                             |
| Clean architecture     | CLI, DNS engine, TUI, and styles separated                    |
| Testable core          | DNS logic should be testable without TUI                      |
| Script-friendly output | Plain output by default                                       |
| Future-ready           | Easy to add history, configs, propagation checks              |

---

## 3. Recommended Tech Stack

| Layer          | Library / Tool                              | Purpose                                   |
| -------------- | ------------------------------------------- | ----------------------------------------- |
| Language       | Go                                          | Native cross-platform CLI                 |
| CLI framework  | Cobra                                       | Commands, flags, help, shell completions  |
| TUI framework  | Bubble Tea                                  | Interactive terminal app                  |
| TUI components | Bubbles                                     | Text input, spinner, progress, list, help |
| Styling        | Lip Gloss                                   | Borders, colors, spacing, layouts         |
| DNS            | Go `net` package + optional custom resolver | DNS lookups                               |
| Config         | Viper, later optional                       | Config files and defaults                 |
| Clipboard      | `atotto/clipboard` or similar               | Copy selected result                      |
| Testing        | Go testing package                          | Unit tests                                |
| CI             | GitHub Actions                              | Linux/Windows/macOS builds                |
| Release        | GoReleaser                                  | Binary releases                           |

---

## 4. High-Level Architecture

```text
# Folder: .
# File: architecture.txt

┌─────────────────────────────────────────────────────────────┐
│                         User                                │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                         Cobra CLI                           │
│  - Parses args                                               │
│  - Parses flags                                              │
│  - Chooses CLI mode or TUI mode                              │
└───────────────┬───────────────────────────────┬─────────────┘
                │                               │
                ▼                               ▼
┌─────────────────────────────┐   ┌───────────────────────────┐
│        Plain CLI Mode        │   │       Bubble Tea TUI      │
│  aqry google.com             │   │  aqry                     │
│  aqry google.com -t MX       │   │  aqry -i google.com       │
└───────────────┬─────────────┘   └──────────────┬────────────┘
                │                                │
                ▼                                ▼
┌─────────────────────────────────────────────────────────────┐
│                        DNS Engine                            │
│  - Domain normalization                                      │
│  - Domain validation                                         │
│  - Resolver selection                                        │
│  - Record lookup                                             │
│  - Timeout handling                                          │
│  - Result normalization                                      │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                         Output                              │
│  - Plain text                                                │
│  - Pretty terminal output                                    │
│  - JSON                                                      │
│  - TUI result panels                                         │
└─────────────────────────────────────────────────────────────┘
```

---

## 5. Repository Structure

```text
# Folder: .
# File: repository-tree.txt

aqry/
├── cmd/
│   └── aqry/
│       └── main.go
├── internal/
│   ├── cli/
│   │   ├── root.go
│   │   ├── flags.go
│   │   ├── output.go
│   │   └── completions.go
│   ├── dns/
│   │   ├── resolver.go
│   │   ├── records.go
│   │   ├── normalize.go
│   │   ├── validate.go
│   │   └── errors.go
│   ├── tui/
│   │   ├── app.go
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── view.go
│   │   ├── messages.go
│   │   ├── keys.go
│   │   ├── layout.go
│   │   ├── modals.go
│   │   └── components.go
│   ├── styles/
│   │   ├── theme.go
│   │   ├── borders.go
│   │   └── colors.go
│   ├── config/
│   │   └── config.go
│   └── version/
│       └── version.go
├── testdata/
│   └── dns-fixtures.json
├── scripts/
│   ├── build.sh
│   ├── build.ps1
│   └── release.sh
├── .github/
│   └── workflows/
│       ├── test.yml
│       └── release.yml
├── go.mod
├── go.sum
├── README.md
├── PRD.md
└── TECHNICAL_DESIGN.md
```

---

## 6. Package Responsibilities

| Package            | Responsibility                                |
| ------------------ | --------------------------------------------- |
| `cmd/aqry`         | Application entrypoint                        |
| `internal/cli`     | Cobra commands, flags, non-interactive output |
| `internal/dns`     | DNS resolution, validation, normalization     |
| `internal/tui`     | Bubble Tea app state, update loop, views      |
| `internal/styles`  | Lip Gloss theme and layout styles             |
| `internal/config`  | Future config file support                    |
| `internal/version` | Build version metadata                        |

---

## 7. CLI Design

## 7.1 Default Command

```bash
# Folder: .
# File: terminal

aqry <domain>
```

Behavior:

- Resolves `A` record by default.
- Prints the first IPv4 address only.
- No styling by default.
- Suitable for scripts.

Example:

```bash
# Folder: .
# File: terminal

aqry google.com
```

Output:

```text
# Folder: .
# File: terminal-output

142.250.182.46
```

---

## 7.2 Interactive Command

```bash
# Folder: .
# File: terminal

aqry
```

Behavior:

- Launches Bubble Tea TUI.
- Starts with an empty domain input.
- Default record type is `A`.

---

## 7.3 Prefilled Interactive Mode

```bash
# Folder: .
# File: terminal

aqry google.com --interactive
```

Behavior:

- Opens TUI with `google.com` already entered.
- Automatically focuses result area after lookup.
- Optional auto-lookup can be enabled later.

---

## 7.4 CLI Flags

| Flag            | Short | Type   | Default | Description        |
| --------------- | ----: | ------ | ------- | ------------------ |
| `--type`        |  `-t` | string | `A`     | DNS record type    |
| `--all`         |       | bool   | `false` | Print all records  |
| `--json`        |       | bool   | `false` | Print JSON         |
| `--server`      |  `-s` | string | system  | DNS resolver       |
| `--timeout`     |       | int    | `3`     | Timeout in seconds |
| `--interactive` |  `-i` | bool   | `false` | Force TUI mode     |
| `--no-color`    |       | bool   | `false` | Disable styling    |
| `--version`     |  `-v` | bool   | `false` | Print version      |

---

## 8. DNS Engine Design

The DNS engine must be independent from Cobra and Bubble Tea.

This allows:

- Fast CLI lookups
- TUI lookups
- Unit tests
- Future API usage

---

## 8.1 DNS Record Type Model

```go
// Folder: internal/dns
// File: records.go

package dns

type RecordType string

const (
	RecordA     RecordType = "A"
	RecordAAAA  RecordType = "AAAA"
	RecordCNAME RecordType = "CNAME"
	RecordMX    RecordType = "MX"
	RecordTXT   RecordType = "TXT"
	RecordNS    RecordType = "NS"
	RecordSOA   RecordType = "SOA"
)

func SupportedRecordTypes() []RecordType {
	return []RecordType{
		RecordA,
		RecordAAAA,
		RecordCNAME,
		RecordMX,
		RecordTXT,
		RecordNS,
		RecordSOA,
	}
}
```

---

## 8.2 Resolver Request Model

```go
// Folder: internal/dns
// File: resolver.go

package dns

import "time"

type LookupRequest struct {
	Domain     string
	RecordType RecordType
	Server     string
	Timeout    time.Duration
	All        bool
}
```

---

## 8.3 Resolver Result Model

```go
// Folder: internal/dns
// File: records.go

package dns

type Record struct {
	Type     RecordType `json:"type"`
	Name     string     `json:"name"`
	Value    string     `json:"value"`
	Priority int        `json:"priority,omitempty"`
	TTL      int        `json:"ttl,omitempty"`
}

type LookupResult struct {
	Domain     string       `json:"domain"`
	RecordType RecordType   `json:"type"`
	Resolver   string       `json:"resolver"`
	Records    []Record     `json:"records"`
}
```

---

## 8.4 Resolver Interface

```go
// Folder: internal/dns
// File: resolver.go

package dns

import "context"

type Resolver interface {
	Lookup(ctx context.Context, req LookupRequest) (LookupResult, error)
}
```

This allows future implementations:

| Resolver         | Purpose                             |
| ---------------- | ----------------------------------- |
| `SystemResolver` | Uses OS DNS                         |
| `CustomResolver` | Uses explicit server like `1.1.1.1` |
| `MockResolver`   | Tests                               |
| `MultiResolver`  | Future propagation comparison       |

---

## 8.5 Domain Normalization

Inputs should be cleaned before validation.

```text
# Folder: internal/dns
# File: normalization-rules.txt

https://google.com/search?q=x  -> google.com
http://example.com/about       -> example.com
google.com/                    -> google.com
 google.com                    -> google.com
```

Rules:

- Trim spaces.
- Remove `http://`.
- Remove `https://`.
- Remove paths.
- Remove query strings.
- Remove trailing slash.
- Lowercase the domain.
- Preserve subdomains.

---

## 8.6 Domain Validation

Validation should reject:

- Empty input
- Spaces
- Protocol-only values
- Domains with invalid characters
- Domains longer than DNS limits
- Labels longer than 63 characters

Validation should allow:

- `google.com`
- `www.google.com`
- `api.example.co.uk`
- Internationalized domains later, optional

---

## 8.7 DNS Lookup Strategy

For MVP, use Go standard library lookups:

| Record Type | Go Function                    |
| ----------- | ------------------------------ |
| A / AAAA    | `LookupIPAddr` or `LookupHost` |
| CNAME       | `LookupCNAME`                  |
| MX          | `LookupMX`                     |
| TXT         | `LookupTXT`                    |
| NS          | `LookupNS`                     |

For SOA, CAA, and richer TTL support, consider adding:

```text
# Folder: .
# File: dependency-note.txt

github.com/miekg/dns
```

Recommended approach:

| Version | DNS implementation                                            |
| ------- | ------------------------------------------------------------- |
| MVP     | Go `net` package                                              |
| v1.1    | Add `miekg/dns` for TTL, SOA, CAA, custom DNS server behavior |
| v1.2    | Add multi-resolver propagation checks                         |

---

## 8.8 Custom Resolver Design

Go supports custom resolver dialing through `net.Resolver`.

For a specific DNS server:

```go
// Folder: internal/dns
// File: custom_resolver_example.go

resolver := &net.Resolver{
	PreferGo: true,
	Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{
			Timeout: timeout,
		}

		return d.DialContext(ctx, "udp", server+":53")
	},
}
```

The custom resolver should support:

- `1.1.1.1`
- `8.8.8.8`
- `9.9.9.9`
- User-provided IP

---

## 9. Bubble Tea TUI Design

## 9.1 Core Bubble Tea Flow

Bubble Tea follows:

```text
# Folder: internal/tui
# File: bubble-tea-flow.txt

Model
  ↓
Update(msg)
  ↓
View()
  ↓
Terminal redraw
```

The app should treat DNS lookup as an async Bubble Tea command.

---

## 9.2 TUI State Model

```go
// Folder: internal/tui
// File: model.go

package tui

type FocusArea int

const (
	FocusDomainInput FocusArea = iota
	FocusRecordType
	FocusResolver
	FocusResults
)

type AppState int

const (
	StateIdle AppState = iota
	StateLoading
	StateSuccess
	StateError
)

type ModalType int

const (
	ModalNone ModalType = iota
	ModalHelp
	ModalRecordPicker
	ModalResolverPicker
)

type Model struct {
	width  int
	height int

	state AppState
	modal ModalType
	focus FocusArea

	domain string

	recordType string
	resolver   string
	timeoutSec int

	records       []string
	selectedRecord int

	err error

	loadingMessage string
	progress       float64

	showAll bool

	input    textinput.Model
	spinner  spinner.Model
	progressBar progress.Model
	help     help.Model
}
```

---

## 9.3 Bubble Tea Messages

```go
// Folder: internal/tui
// File: messages.go

package tui

import dnsengine "github.com/yourname/aqry/internal/dns"

type lookupStartedMsg struct{}

type lookupProgressMsg struct {
	message  string
	percent float64
}

type lookupSuccessMsg struct {
	result dnsengine.LookupResult
}

type lookupErrorMsg struct {
	err error
}

type copySuccessMsg struct{}

type copyErrorMsg struct {
	err error
}
```

---

## 9.4 Lookup Command

DNS lookup should run as a Bubble Tea command.

```go
// Folder: internal/tui
// File: commands.go

package tui

func lookupCmd(req dns.LookupRequest) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			req.Timeout,
		)
		defer cancel()

		result, err := resolver.Lookup(ctx, req)
		if err != nil {
			return lookupErrorMsg{err: err}
		}

		return lookupSuccessMsg{result: result}
	}
}
```

Requirements:

- UI must not block during DNS lookup.
- Spinner should continue while lookup is running.
- Keyboard events should still be handled.
- `Esc` should cancel lookup later if cancellation is implemented.

---

## 9.5 Update Function Responsibilities

The `Update` function handles:

| Message / Key       | Responsibility                   |
| ------------------- | -------------------------------- |
| `tea.WindowSizeMsg` | Update layout dimensions         |
| `tea.KeyMsg`        | Keyboard navigation              |
| `spinner.TickMsg`   | Spinner animation                |
| `lookupSuccessMsg`  | Store records, set success state |
| `lookupErrorMsg`    | Store error, set error state     |
| `lookupProgressMsg` | Update progress bar              |
| `copySuccessMsg`    | Show status message              |

---

## 10. Keyboard Controls

## 10.1 Global Controls

| Key          | Action                      |
| ------------ | --------------------------- |
| `q`          | Quit                        |
| `ctrl+c`     | Force quit                  |
| `?`          | Toggle help modal           |
| `esc`        | Close modal / cancel        |
| `tab`        | Next focus area             |
| `shift+tab`  | Previous focus area         |
| `enter`      | Run lookup or confirm modal |
| `r`          | Open record type picker     |
| `s`          | Open resolver picker        |
| `c`          | Copy selected result        |
| `a`          | Toggle all results          |
| `j` / `down` | Move down                   |
| `k` / `up`   | Move up                     |

---

## 10.2 Focus Behavior

| Current Focus | `Tab` Moves To |
| ------------- | -------------- |
| Domain input  | Record type    |
| Record type   | Resolver       |
| Resolver      | Results        |
| Results       | Domain input   |

---

## 10.3 Modal Behavior

| Modal           | `Enter`            | `Esc` |
| --------------- | ------------------ | ----- |
| Help            | No-op              | Close |
| Record picker   | Select record type | Close |
| Resolver picker | Select resolver    | Close |

---

## 11. UI Layout Design

## 11.1 Layout Regions

```text
# Folder: internal/tui
# File: layout-regions.txt

┌──────────────────────────────────────────────┐
│ Header                                       │
├──────────────────────────────────────────────┤
│ Main content                                 │
│ - Domain input                               │
│ - Record type selector                       │
│ - Resolver selector                          │
│ - Results / loading / error panel            │
├──────────────────────────────────────────────┤
│ Footer help bar                              │
└──────────────────────────────────────────────┘
```

---

## 11.2 Desktop/Wide Layout

For width greater than `100`:

```text
# Folder: internal/tui
# File: wide-layout.txt

╭────────────────────────────────────────────────────────────────────╮
│ aqry                                                   DNS Resolver │
├───────────────────────────────┬────────────────────────────────────┤
│ Input                         │ Results                            │
│                               │                                    │
│ Domain                        │ A Record                           │
│ google.com                    │ 142.250.182.46                     │
│                               │                                    │
│ Record Type                   │ Metadata                           │
│ [ A ] AAAA CNAME MX TXT NS    │ Resolver: System DNS               │
│                               │ Status: Success                    │
│ Resolver                      │                                    │
│ System DNS                    │                                    │
├───────────────────────────────┴────────────────────────────────────┤
│ Enter lookup  r records  s resolver  ? help  q quit                │
╰────────────────────────────────────────────────────────────────────╯
```

---

## 11.3 Standard Layout

For width `60–100`:

```text
# Folder: internal/tui
# File: standard-layout.txt

╭────────────────────────────────────────────────────────────╮
│ aqry                                          DNS Resolver │
├────────────────────────────────────────────────────────────┤
│ Domain                                                     │
│ ╭──────────────────────────────────────────────────────╮   │
│ │ google.com                                           │   │
│ ╰──────────────────────────────────────────────────────╯   │
│                                                            │
│ Record Type                                                │
│ [ A ]  AAAA  CNAME  MX  TXT  NS                           │
│                                                            │
│ Result                                                     │
│ ╭─ IPv4 ───────────────────────────────────────────────╮   │
│ │ 142.250.182.46                                      │   │
│ ╰──────────────────────────────────────────────────────╯   │
├────────────────────────────────────────────────────────────┤
│ Enter lookup  r records  s resolver  ? help  q quit        │
╰────────────────────────────────────────────────────────────╯
```

---

## 11.4 Compact Layout

For width below `60`:

```text
# Folder: internal/tui
# File: compact-layout.txt

aqry

Domain
google.com

Type
A

IPv4
142.250.182.46

? help · q quit
```

Compact mode should prioritize usability over borders.

---

## 12. Lip Gloss Styling System

## 12.1 Theme Tokens

```go
// Folder: internal/styles
// File: theme.go

package styles

type Theme struct {
	Primary string
	Muted   string
	Success string
	Warning string
	Danger  string
	Border  string
	Text    string
	Surface string
}
```

---

## 12.2 Style Categories

| Style              | Purpose                |
| ------------------ | ---------------------- |
| `AppFrame`         | Outer app border       |
| `Header`           | Top title bar          |
| `Footer`           | Bottom keybinding bar  |
| `Panel`            | Normal content panel   |
| `FocusedPanel`     | Active section         |
| `Input`            | Domain text input      |
| `RecordPill`       | Inactive record type   |
| `RecordPillActive` | Selected record type   |
| `ResultValue`      | Main IP address        |
| `ErrorBox`         | Error display          |
| `Modal`            | Help / picker overlays |

---

## 12.3 Visual Rules

- Use rounded borders where terminal supports them.
- Use focus borders for active areas.
- Keep primary result large and clear.
- Do not overuse color.
- Do not put every element inside a box.
- Use whitespace for visual hierarchy.
- Footer should always show available actions.

---

## 13. Bubbles Components

| Component  | Package             | Usage                            |
| ---------- | ------------------- | -------------------------------- |
| Text input | `bubbles/textinput` | Domain input                     |
| Spinner    | `bubbles/spinner`   | DNS lookup loading               |
| Progress   | `bubbles/progress`  | Phase-based lookup progress      |
| Help       | `bubbles/help`      | Keybinding help                  |
| Key        | `bubbles/key`       | Keyboard shortcut definitions    |
| List       | `bubbles/list`      | Record type and resolver pickers |

---

## 14. Progress Bar Design

The progress bar should represent lookup phases, not real network progress.

| Phase            | Progress | Message                 |
| ---------------- | -------: | ----------------------- |
| Validate domain  |   `0.20` | `Checking domain...`    |
| Prepare resolver |   `0.35` | `Preparing resolver...` |
| Query DNS        |   `0.60` | `Querying DNS...`       |
| Parse response   |   `0.85` | `Parsing response...`   |
| Done             |   `1.00` | `Resolved successfully` |

Important:

- If lookup completes quickly, progress should finish immediately.
- Do not artificially delay CLI output.
- TUI mode may animate briefly for polish, but never feel slow.

---

## 15. Error Handling

## 15.1 Error Types

```go
// Folder: internal/dns
// File: errors.go

package dns

type ErrorKind string

const (
	ErrInvalidDomain    ErrorKind = "invalid_domain"
	ErrNoRecords        ErrorKind = "no_records"
	ErrTimeout          ErrorKind = "timeout"
	ErrResolverFailed   ErrorKind = "resolver_failed"
	ErrUnsupportedType  ErrorKind = "unsupported_type"
)
```

---

## 15.2 Exit Codes

| Exit Code | Meaning                 |
| --------: | ----------------------- |
|       `0` | Success                 |
|       `1` | Invalid input           |
|       `2` | DNS lookup failed       |
|       `3` | Timeout                 |
|       `4` | Unsupported record type |
|       `5` | Internal error          |

---

## 15.3 TUI Error Display

Errors should be human-readable.

Example:

```text
# Folder: internal/tui
# File: error-view-example.txt

Could not resolve domain

Domain
gogle.invalid

Reason
No A record found

Suggestions
• Check spelling
• Try another record type
• Try another DNS resolver
```

---

## 16. Output Design

## 16.1 Plain Output

Default:

```text
# Folder: .
# File: terminal-output

142.250.182.46
```

Rules:

- No labels
- No colors
- First result only
- Newline at end

---

## 16.2 All Records Output

```text
# Folder: .
# File: terminal-output

142.250.182.46
142.250.193.78
```

---

## 16.3 JSON Output

```json
{
  "domain": "google.com",
  "type": "A",
  "resolver": "system",
  "records": [
    {
      "type": "A",
      "name": "google.com",
      "value": "142.250.182.46",
      "ttl": 300
    }
  ]
}
```

---

## 17. Config Design

Config is optional for v1.1.

Potential config path:

| OS      | Path                                             |
| ------- | ------------------------------------------------ |
| Linux   | `~/.config/aqry/config.toml`                     |
| Windows | `%APPDATA%\aqry\config.toml`                     |
| macOS   | `~/Library/Application Support/aqry/config.toml` |

Example:

```toml
# Folder: ~/.config/aqry
# File: config.toml

default_record = "A"
default_resolver = "system"
timeout = 3
show_all = false
theme = "default"
```

---

## 18. Testing Strategy

## 18.1 Unit Tests

| Package           | Tests                                       |
| ----------------- | ------------------------------------------- |
| `internal/dns`    | Validation, normalization, lookup parsing   |
| `internal/cli`    | Flag parsing, output formatting             |
| `internal/tui`    | Key handling, model transitions             |
| `internal/styles` | Snapshot-style string tests where practical |

---

## 18.2 DNS Engine Tests

Test without relying heavily on live DNS.

Required tests:

- Valid domain passes validation.
- Invalid domain fails validation.
- URL input normalizes to domain.
- Unsupported record type fails.
- Empty response becomes `ErrNoRecords`.
- Timeout maps to `ErrTimeout`.

---

## 18.3 TUI Model Tests

Test state transitions:

| Initial State | Action          | Expected State         |
| ------------- | --------------- | ---------------------- |
| Idle          | Type domain     | Domain updated         |
| Idle          | Press `r`       | Record picker modal    |
| Modal open    | Press `esc`     | Modal closes           |
| Idle          | Press `enter`   | Loading                |
| Loading       | Success message | Success                |
| Loading       | Error message   | Error                  |
| Success       | Press `c`       | Copy command triggered |

---

## 19. Build and Release

## 19.1 Local Build

```bash
# Folder: .
# File: terminal

go build -o bin/aqry ./cmd/aqry
```

---

## 19.2 Linux Build

```bash
# Folder: .
# File: terminal

GOOS=linux GOARCH=amd64 go build -o dist/aqry-linux-amd64 ./cmd/aqry
```

---

## 19.3 Windows Build

```bash
# Folder: .
# File: terminal

GOOS=windows GOARCH=amd64 go build -o dist/aqry-windows-amd64.exe ./cmd/aqry
```

---

## 19.4 Arch Linux Install During Development

```bash
# Folder: .
# File: terminal

go install ./cmd/aqry
```

---

## 19.5 Release Targets

| Platform      | Artifact                   |
| ------------- | -------------------------- |
| Linux amd64   | `aqry-linux-amd64.tar.gz`  |
| Linux arm64   | `aqry-linux-arm64.tar.gz`  |
| Windows amd64 | `aqry-windows-amd64.zip`   |
| macOS amd64   | `aqry-darwin-amd64.tar.gz` |
| macOS arm64   | `aqry-darwin-arm64.tar.gz` |

---

## 20. CI Pipeline

## 20.1 Test Workflow

```yaml
# Folder: .github/workflows
# File: test.yml

name: Test

on:
  push:
    branches: [main]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.23"

      - run: go test ./...
```

---

## 20.2 Release Workflow

Use GoReleaser later.

```yaml
# Folder: .github/workflows
# File: release.yml

name: Release

on:
  push:
    tags:
      - "v*"

jobs:
  goreleaser:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.23"

      - uses: goreleaser/goreleaser-action@v6
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
```

---

## 21. Implementation Milestones

## Milestone 1 — CLI Foundation

Deliverables:

- Go module
- Cobra root command
- `aqry <domain>`
- A record lookup
- Plain output
- Basic errors

Acceptance:

```bash
# Folder: .
# File: terminal

aqry google.com
```

Returns:

```text
# Folder: .
# File: terminal-output

142.250.182.46
```

---

## Milestone 2 — DNS Engine

Deliverables:

- `internal/dns`
- Record type model
- Lookup request/result model
- Domain normalization
- Domain validation
- A, AAAA, CNAME, MX, TXT, NS support

---

## Milestone 3 — Bubble Tea MVP

Deliverables:

- TUI launches with `aqry`
- Text input works
- Enter triggers lookup
- Result appears
- Error appears
- Quit works

---

## Milestone 4 — Modern UI

Deliverables:

- Lip Gloss theme
- Header
- Footer
- Result panel
- Focus states
- Responsive layout

---

## Milestone 5 — Interactive Features

Deliverables:

- Record type picker
- Resolver picker
- Help modal
- Keyboard navigation
- Loading spinner
- Progress bar

---

## Milestone 6 — Polish and Release

Deliverables:

- Tests
- CI
- Cross-platform builds
- README
- Install instructions
- GitHub release artifacts

---

## 22. Main Risks

| Risk                                  | Impact | Mitigation                                |
| ------------------------------------- | ------ | ----------------------------------------- |
| Terminal rendering differs on Windows | Medium | Test with Windows Terminal and PowerShell |
| DNS TTL unavailable in stdlib         | Medium | Use `miekg/dns` in v1.1                   |
| UI becomes too complex                | Medium | Keep MVP focused                          |
| Lookup feels slow                     | Low    | Avoid artificial delays                   |
| Too much styling noise                | Medium | Use restrained theme tokens               |

---

## 23. Final Technical Direction

The best architecture is:

```text
# Folder: .
# File: final-architecture.txt

Cobra
  → chooses CLI mode or TUI mode

DNS engine
  → reusable resolver logic

Bubble Tea
  → interactive event loop

Bubbles
  → input, spinner, progress, list, help

Lip Gloss
  → modern layout and visual system
```

`aqry` should remain fast and minimal in plain CLI mode while offering a polished interactive experience when launched without arguments.

The most important technical rule:

> Keep DNS logic independent from the TUI.

That separation makes the app easier to test, faster to maintain, and simpler to expand later.
