# PRD.md — DNS-Resolver Interactive CLI

## Product Name

**aqry**

## Product Type

Cross-platform interactive DNS lookup CLI/TUI built with **Go**, **Bubble Tea**, **Bubbles**, and **Lip Gloss**.

## One-Line Summary

`aqry` is a fast, modern terminal DNS resolver that returns an IPv4 address by default and provides an interactive Bubble Tea interface for selecting DNS record types without repeatedly typing terminal flags.

---

## 1. Product Vision

Most DNS lookup tools are powerful but not beginner-friendly or visually polished. Tools like `dig`, `nslookup`, and `host` are excellent for advanced users, but they require remembering flags, record types, and output formats.

`aqry` should feel like a modern terminal product: clean, fast, keyboard-driven, visually polished, and useful for both quick lookups and interactive DNS exploration.

The default behavior should remain simple:

```text
aqry google.com

142.250.182.46
```

But when opened interactively:

```text
aqry
```

The user gets a full-screen terminal UI for entering a domain, selecting record types, viewing DNS results, copying output, and switching views using keyboard navigation.

Bubble Tea is suitable for this because it is a Go TUI framework based on The Elm Architecture, designed for simple and complex terminal applications. Bubbles provides reusable Bubble Tea components such as lists, spinners, help views, and status messages. Lip Gloss provides declarative terminal styling and layout utilities.

---

## 2. Goals

## 2.1 Primary Goals

- Provide a **fast DNS lookup CLI** that works on:
  - Arch Linux
  - Ubuntu/Debian
  - Windows
  - macOS later, optional

- Return **IPv4 / A record by default**.
- Provide a **modern interactive TUI** when the user runs `aqry` without arguments.
- Allow users to select DNS record types without typing long flags.
- Make DNS results easy to read, scan, and copy.
- Keep the tool lightweight and dependency-conscious.
- Use keyboard-first navigation.
- Support polished terminal UI elements:
  - Keyboard navigation
  - Animations
  - Loading spinners
  - Progress bars
  - Help menu
  - Record type selector
  - Result panels
  - Error states
  - Empty states
  - Responsive layout

---

## 2.2 Secondary Goals

- Support custom DNS resolvers.
- Support multiple output modes:
  - Plain text
  - Pretty terminal output
  - JSON
  - Copy-friendly output

- Support future features like:
  - DNS propagation comparison
  - Lookup history
  - Favorite domains
  - Saved resolver profiles
  - Export results

---

## 2.3 Non-Goals

The first version should **not** try to replace full DNS tools like `dig`.

Initial non-goals:

- No full DNS packet inspector.
- No zone transfer functionality.
- No WHOIS lookup in v1.
- No DNSSEC validation in v1.
- No network monitoring dashboard in v1.
- No GUI desktop application.
- No web application.
- No account/login system.
- No cloud sync.

---

## 3. Target Users

## 3.1 Primary Users

| User Type            | Need                                           |
| -------------------- | ---------------------------------------------- |
| Developers           | Quickly check domain IPs and records           |
| Students             | Learn DNS record types visually                |
| DevOps beginners     | Inspect records without memorizing commands    |
| Terminal-first users | Fast keyboard-driven DNS workflow              |
| Web developers       | Verify domain setup, A records, CNAME, TXT, MX |

---

## 3.2 User Skill Level

The product should work well for:

- Beginner terminal users
- Intermediate developers
- Advanced users who want a faster workflow

---

## 4. Core Use Cases

## 4.1 Quick IPv4 Lookup

User runs:

```text
aqry google.com
```

Expected output:

```text
142.250.182.46
```

This should be the fastest and cleanest path.

---

## 4.2 Interactive Lookup

User runs:

```text
aqry
```

Expected behavior:

- Opens full-screen Bubble Tea UI.
- Focus starts on domain input.
- User types a domain.
- User presses `Enter`.
- App resolves IPv4 / A record by default.
- Results appear in a polished result panel.

---

## 4.3 Select Record Type Without Flags

User should be able to select DNS record type from a keyboard menu.

Supported record types:

| Record Type | Description                         | v1 Priority |
| ----------- | ----------------------------------- | ----------- |
| A           | IPv4 address                        | Required    |
| AAAA        | IPv6 address                        | Required    |
| CNAME       | Canonical name                      | Required    |
| MX          | Mail exchange                       | Required    |
| TXT         | Text records                        | Required    |
| NS          | Name servers                        | Required    |
| SOA         | Start of authority                  | Required    |
| SRV         | Service records                     | Optional    |
| CAA         | Certificate authority authorization | Optional    |

Default selected type:

```text
A
```

---

## 4.4 Help Menu

User presses:

```text
?
```

Expected result:

A help overlay appears showing keybindings.

Example:

```text
┌─ Help ─────────────────────────────┐
│ Enter       Run lookup              │
│ Tab         Switch section          │
│ ↑ / ↓       Move selection          │
│ r           Change record type      │
│ s           Change DNS server       │
│ c           Copy selected result    │
│ j / k       Navigate results        │
│ esc         Close modal             │
│ q           Quit                    │
│ ctrl+c      Force quit              │
└────────────────────────────────────┘
```

---

## 5. CLI Behavior

## 5.1 Basic Command

```text
aqry <domain>
```

Default behavior:

- Resolve A record.
- Print first IPv4 address only.
- No UI.
- No decorations.
- Script-friendly.

Example:

```text
aqry google.com
```

Output:

```text
142.250.182.46
```

---

## 5.2 Interactive Mode

```text
aqry
```

Behavior:

- Launches Bubble Tea TUI.
- Starts on domain input screen.
- Default record type is `A`.
- Default DNS resolver is system resolver.

---

## 5.3 Optional Flags

Flags should exist for power users but should not be required for normal use.

| Flag            | Short | Description        | Example                       |
| --------------- | ----: | ------------------ | ----------------------------- |
| `--type`        |  `-t` | DNS record type    | `aqry google.com -t MX`       |
| `--all`         |       | Show all results   | `aqry google.com --all`       |
| `--json`        |       | JSON output        | `aqry google.com --json`      |
| `--server`      |  `-s` | DNS server         | `aqry google.com -s 1.1.1.1`  |
| `--timeout`     |       | Timeout in seconds | `aqry google.com --timeout 3` |
| `--interactive` |  `-i` | Force TUI mode     | `aqry google.com -i`          |
| `--no-color`    |       | Disable colors     | `aqry google.com --no-color`  |
| `--version`     |  `-v` | Print version      | `aqry -v`                     |
| `--help`        |  `-h` | Print help         | `aqry -h`                     |

---

## 5.4 Command Modes

| Command                  | Mode      | Output                          |
| ------------------------ | --------- | ------------------------------- |
| `aqry google.com`        | Plain CLI | First IPv4 only                 |
| `aqry google.com --all`  | Plain CLI | All IPv4 records                |
| `aqry google.com -t MX`  | Plain CLI | MX records                      |
| `aqry google.com --json` | Plain CLI | JSON                            |
| `aqry`                   | TUI       | Interactive app                 |
| `aqry -i google.com`     | TUI       | Opens app with domain prefilled |

---

## 6. TUI Product Requirements

## 6.1 UI Personality

The UI should feel:

- Modern
- Calm
- Sharp
- Terminal-native
- Not overdecorated
- Not like a generic AI-generated box layout
- Inspired by high-quality tools like:
  - LazyGit
  - GitHub CLI interactive views
  - Charmbracelet tools
  - Modern terminal dashboards

---

## 6.2 Design Principles

| Principle                  | Requirement                                                         |
| -------------------------- | ------------------------------------------------------------------- |
| Fast first                 | UI should not slow down quick lookups                               |
| Keyboard-first             | Every major action must be keyboard-accessible                      |
| Minimal typing             | Record type selection should use menus, not flags                   |
| Visual hierarchy           | Domain, record type, result, and status should be clearly separated |
| Progressive disclosure     | Advanced options hidden until needed                                |
| Script-friendly            | Non-interactive mode must remain clean                              |
| Responsive terminal layout | Works from compact to wide terminals                                |
| Accessible colors          | Must work in dark and light terminals                               |

---

## 7. TUI Screens

## 7.1 Main Lookup Screen

Default screen when running:

```text
aqry
```

Wireframe:

```text
╭────────────────────────────────────────────────────────────╮
│ aqry                                        DNS Lookup │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  Domain                                                    │
│  ╭──────────────────────────────────────────────────────╮  │
│  │ google.com                                           │  │
│  ╰──────────────────────────────────────────────────────╯  │
│                                                            │
│  Record Type                                               │
│  [ A ]  AAAA  CNAME  MX  TXT  NS  SOA                     │
│                                                            │
│  Resolver                                                  │
│  System DNS                                                │
│                                                            │
│  Status                                                    │
│  Ready                                                     │
│                                                            │
├────────────────────────────────────────────────────────────┤
│ Enter lookup  Tab focus  ? help  q quit                    │
╰────────────────────────────────────────────────────────────╯
```

---

## 7.2 Loading State

When user presses `Enter`:

```text
╭────────────────────────────────────────────────────────────╮
│ aqry                                        DNS Lookup │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  Domain                                                    │
│  google.com                                                │
│                                                            │
│  Record Type                                               │
│  A                                                         │
│                                                            │
│  Resolving                                                 │
│  ⠋ Querying DNS resolver...                                │
│                                                            │
│  Progress                                                  │
│  ████████████░░░░░░░░░░░░  48%                             │
│                                                            │
├────────────────────────────────────────────────────────────┤
│ esc cancel  ? help  q quit                                 │
╰────────────────────────────────────────────────────────────╯
```

Requirements:

- Spinner animates while lookup is running.
- Progress bar shows lookup phase, not fake exact network progress.
- Progress can represent:
  - Input validation
  - Resolver selection
  - DNS query running
  - Response parsing
  - Rendering result

---

## 7.3 Result State

Default result for A record:

```text
╭────────────────────────────────────────────────────────────╮
│ aqry                                        DNS Lookup │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  google.com                                                │
│  A Record                                                  │
│                                                            │
│  ╭─ IPv4 Result ────────────────────────────────────────╮  │
│  │ 142.250.182.46                                      │  │
│  ╰──────────────────────────────────────────────────────╯  │
│                                                            │
│  TTL                                                       │
│  300s                                                      │
│                                                            │
│  Resolver                                                  │
│  System DNS                                                │
│                                                            │
├────────────────────────────────────────────────────────────┤
│ c copy  r record type  enter lookup again  ? help  q quit  │
╰────────────────────────────────────────────────────────────╯
```

---

## 7.4 Multiple Results State

When more than one record exists:

```text
╭────────────────────────────────────────────────────────────╮
│ aqry                                        DNS Lookup │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  google.com                                                │
│  A Records                                                 │
│                                                            │
│  ┌────┬─────────────────┬──────┐                           │
│  │ #  │ IPv4 Address    │ TTL  │                           │
│  ├────┼─────────────────┼──────┤                           │
│  │ 1  │ 142.250.182.46  │ 300s │                           │
│  │ 2  │ 142.250.193.78  │ 300s │                           │
│  └────┴─────────────────┴──────┘                           │
│                                                            │
├────────────────────────────────────────────────────────────┤
│ ↑↓ select  c copy  a show all  ? help  q quit              │
╰────────────────────────────────────────────────────────────╯
```

Default plain CLI output should still print only the first IPv4 unless `--all` is used.

---

## 7.5 Record Type Selector

When user presses `r`:

```text
╭─ Select Record Type ─────────────────────────────╮
│                                                   │
│  › A      IPv4 address                            │
│    AAAA   IPv6 address                            │
│    CNAME  Canonical name                          │
│    MX     Mail server                             │
│    TXT    Verification / policy text              │
│    NS     Name servers                            │
│    SOA    Zone authority                          │
│                                                   │
│  Enter select  Esc cancel                         │
╰───────────────────────────────────────────────────╯
```

Requirements:

- `↑` / `↓` navigate.
- `j` / `k` navigate.
- `Enter` selects.
- `Esc` closes.
- Selected record type becomes active in the main screen.
- No repeated terminal flag typing required.

---

## 7.6 Resolver Selector

When user presses `s`:

```text
╭─ DNS Resolver ────────────────────────────────────╮
│                                                   │
│  › System DNS                                     │
│    Cloudflare       1.1.1.1                       │
│    Google           8.8.8.8                       │
│    Quad9            9.9.9.9                       │
│    Custom...                                      │
│                                                   │
│  Enter select  Esc cancel                         │
╰───────────────────────────────────────────────────╯
```

Requirements:

- System DNS should be default.
- Custom resolver opens input mode.
- Resolver name and IP should be displayed clearly.

---

## 7.7 Error State

Example invalid domain:

```text
╭────────────────────────────────────────────────────────────╮
│ aqry                                        DNS Lookup │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  Could not resolve domain                                  │
│                                                            │
│  Domain                                                    │
│  gogle.invalid                                             │
│                                                            │
│  Reason                                                    │
│  No A record found                                         │
│                                                            │
│  Suggestions                                               │
│  • Check spelling                                          │
│  • Try another record type                                 │
│  • Try a different DNS resolver                            │
│                                                            │
├────────────────────────────────────────────────────────────┤
│ r record type  s resolver  enter retry  ? help  q quit     │
╰────────────────────────────────────────────────────────────╯
```

---

## 8. Keyboard Navigation

## 8.1 Global Keybindings

| Key         | Action                             |
| ----------- | ---------------------------------- |
| `q`         | Quit                               |
| `ctrl+c`    | Force quit                         |
| `?`         | Toggle help menu                   |
| `esc`       | Close modal / cancel current state |
| `tab`       | Move focus forward                 |
| `shift+tab` | Move focus backward                |
| `enter`     | Confirm / run lookup               |
| `r`         | Open record type selector          |
| `s`         | Open DNS resolver selector         |
| `c`         | Copy selected result               |
| `a`         | Toggle show all results            |
| `j`         | Move down                          |
| `k`         | Move up                            |
| `↑`         | Move up                            |
| `↓`         | Move down                          |

---

## 8.2 Focus Areas

Main screen should support focus movement between:

1. Domain input
2. Record type selector
3. Resolver selector
4. Result panel
5. Footer/help area

---

## 9. Animation Requirements

## 9.1 Spinner

Use spinner during DNS lookup.

Requirements:

- Spinner starts immediately after lookup begins.
- Spinner stops when result or error arrives.
- Spinner should not block keyboard events.
- `esc` should cancel lookup if possible.

Recommended states:

| State      | Message                    |
| ---------- | -------------------------- |
| Validating | `Checking domain...`       |
| Resolving  | `Querying DNS resolver...` |
| Parsing    | `Parsing DNS response...`  |
| Done       | `Resolved successfully`    |
| Error      | `Lookup failed`            |

---

## 9.2 Progress Bar

Progress bar should represent lookup phases.

Example phases:

| Phase            | Progress |
| ---------------- | -------: |
| Validate domain  |      20% |
| Select resolver  |      35% |
| Send query       |      55% |
| Receive response |      80% |
| Render result    |     100% |

Important:

- Do not fake exact network progress.
- Progress should be phase-based.
- If lookup finishes quickly, progress should animate smoothly but not delay output too much.

---

## 9.3 Micro-Interactions

Recommended UI polish:

- Highlight active focus panel.
- Subtle status line updates.
- Smooth result replacement.
- Dim inactive sections.
- Show selected record type as a pill.
- Use success/error indicators:
  - Success: `✓`
  - Warning: `!`
  - Error: `✕`

- Do not overuse icons.

---

## 10. Visual Design System

## 10.1 Layout

Use a structured layout:

```text
Header
Main content
Context panel
Footer help bar
```

---

## 10.2 Recommended Visual Hierarchy

| Element              | Visual Treatment          |
| -------------------- | ------------------------- |
| App name             | Bold, accent color        |
| Current mode         | Muted label               |
| Domain input         | Bordered focused box      |
| Selected record type | Pill-style highlight      |
| Result               | Large, high-contrast text |
| Metadata             | Muted text                |
| Errors               | Clear but not aggressive  |
| Help footer          | Dim, compact              |

---

## 10.3 Color System

Use semantic colors through Lip Gloss styles.

| Token     | Purpose                    |
| --------- | -------------------------- |
| `primary` | App title, selected states |
| `muted`   | Metadata, helper text      |
| `success` | Successful lookup          |
| `warning` | Partial result             |
| `danger`  | Errors                     |
| `border`  | Panels                     |
| `surface` | Modal-like sections        |
| `text`    | Primary output             |

Requirements:

- Must work on dark terminals.
- Should degrade gracefully on limited-color terminals.
- Support `--no-color`.

---

## 10.4 Layout Breakpoints

|   Terminal Width | Layout                       |
| ---------------: | ---------------------------- |
|   `< 60 columns` | Compact single-column        |
| `60–100 columns` | Standard layout              |
|  `> 100 columns` | Split layout with side panel |

---

## 11. DNS Behavior

## 11.1 Default Lookup

Default record:

```text
A
```

Default output:

- First IPv4 address only.
- Plain text in CLI mode.
- Pretty panel in TUI mode.

---

## 11.2 Record Type Behavior

| Record Type | Expected Output                         |
| ----------- | --------------------------------------- |
| A           | IPv4 addresses                          |
| AAAA        | IPv6 addresses                          |
| CNAME       | Canonical name                          |
| MX          | Priority + mail server                  |
| TXT         | TXT values                              |
| NS          | Name servers                            |
| SOA         | Primary NS, admin, serial, TTL metadata |
| CAA         | Tag, value, flags                       |

---

## 11.3 Resolver Options

Default resolver:

```text
System DNS
```

Built-in resolver presets:

| Name       | Address    |
| ---------- | ---------- |
| System DNS | OS default |
| Cloudflare | `1.1.1.1`  |
| Google     | `8.8.8.8`  |
| Quad9      | `9.9.9.9`  |

---

## 11.4 Timeout

Default timeout:

```text
3 seconds
```

Configurable via:

```text
--timeout 5
```

Interactive mode should show a timeout error clearly.

---

## 12. Output Formats

## 12.1 Plain Default

Command:

```text
aqry google.com
```

Output:

```text
142.250.182.46
```

---

## 12.2 All Records

Command:

```text
aqry google.com --all
```

Output:

```text
142.250.182.46
142.250.193.78
```

---

## 12.3 JSON

Command:

```text
aqry google.com --json
```

Output:

```json
{
  "domain": "google.com",
  "type": "A",
  "resolver": "system",
  "records": [
    {
      "value": "142.250.182.46",
      "ttl": 300
    }
  ]
}
```

---

## 12.4 Pretty Output

Optional future flag:

```text
aqry google.com --pretty
```

Example:

```text
✓ google.com

A Record
142.250.182.46

Resolver
System DNS
```

---

## 13. Technical Architecture

## 13.1 Recommended Stack

| Layer              | Library                                    |
| ------------------ | ------------------------------------------ |
| Language           | Go                                         |
| CLI commands/flags | Cobra                                      |
| TUI framework      | Bubble Tea                                 |
| TUI components     | Bubbles                                    |
| Styling/layout     | Lip Gloss                                  |
| DNS lookup         | Go standard library + optional DNS package |
| Clipboard          | Cross-platform clipboard package           |
| Config             | Viper, optional                            |
| Testing            | Go testing package                         |
| CI                 | GitHub Actions                             |

---

## 13.2 Why This Stack

| Tool       | Reason                                               |
| ---------- | ---------------------------------------------------- |
| Go         | Fast native binaries, easy cross-compilation         |
| Cobra      | Reliable command and flag structure                  |
| Bubble Tea | Full interactive terminal UI                         |
| Bubbles    | Ready-made text input, spinner, list, help, progress |
| Lip Gloss  | Modern styling, layout, borders, colors              |
| Viper      | Future config support                                |

---

## 13.3 App Architecture

Recommended internal architecture:

```text
aqry/
├── cmd/
│   └── aqry/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── view.go
│   │   └── keys.go
│   ├── dns/
│   │   ├── resolver.go
│   │   ├── records.go
│   │   └── errors.go
│   ├── cli/
│   │   ├── root.go
│   │   └── output.go
│   ├── styles/
│   │   └── theme.go
│   └── config/
│       └── config.go
├── testdata/
├── go.mod
├── go.sum
├── README.md
├── TECHNICAL_DESIGN.md
└── PRD.md
```

---

## 13.4 Core Modules

| Module            | Responsibility                      |
| ----------------- | ----------------------------------- |
| `cmd/aqry`        | Binary entrypoint                   |
| `internal/cli`    | Cobra commands, flags, output modes |
| `internal/app`    | Bubble Tea model/update/view        |
| `internal/dns`    | DNS resolving logic                 |
| `internal/styles` | Lip Gloss theme                     |
| `internal/config` | Optional config loading             |
| `testdata`        | Test fixtures                       |

---

## 14. Bubble Tea State Model

## 14.1 Main Model Fields

The TUI model should track:

| Field                | Purpose                            |
| -------------------- | ---------------------------------- |
| `domain`             | Current domain input               |
| `recordType`         | Selected DNS record type           |
| `resolver`           | Selected DNS resolver              |
| `records`            | Current lookup results             |
| `selectedRecord`     | Selected result index              |
| `loading`            | Whether lookup is running          |
| `progress`           | Current phase progress             |
| `err`                | Current error                      |
| `focusedPanel`       | Current active UI section          |
| `showHelp`           | Help overlay visibility            |
| `showRecordPicker`   | Record selector modal visibility   |
| `showResolverPicker` | Resolver selector modal visibility |
| `width`              | Terminal width                     |
| `height`             | Terminal height                    |

---

## 14.2 Bubble Tea Flow

```text
Init
 ↓
View initial input screen
 ↓
User types domain
 ↓
Update model
 ↓
User presses Enter
 ↓
Start DNS command
 ↓
Spinner/progress messages update UI
 ↓
DNS result message received
 ↓
Render result state
```

---

## 15. Validation Rules

## 15.1 Domain Validation

Input should reject:

- Empty string
- Spaces
- Protocol prefixes like `https://`
- Paths like `/about`
- Invalid domain characters

Input should normalize:

| Input                     | Normalized        |
| ------------------------- | ----------------- |
| `https://google.com`      | `google.com`      |
| `http://example.com/path` | `example.com`     |
| `www.example.com/`        | `www.example.com` |
| `google.com`              | `google.com`      |

---

## 15.2 Error Messages

| Case             | Message                                  |
| ---------------- | ---------------------------------------- |
| Empty input      | `Enter a domain to continue`             |
| Invalid domain   | `This does not look like a valid domain` |
| No record found  | `No A record found for this domain`      |
| Timeout          | `DNS lookup timed out`                   |
| Network error    | `Could not reach DNS resolver`           |
| Unsupported type | `Record type is not supported yet`       |

---

## 16. Accessibility Requirements

- Support `--no-color`.
- Do not rely only on color to communicate state.
- Use labels and symbols together:
  - `✓ Success`
  - `✕ Error`
  - `! Warning`

- Keyboard-only operation required.
- Help menu required.
- UI should remain usable in compact terminals.
- Avoid excessive animation speed.

---

## 17. Performance Requirements

| Requirement           |                        Target |
| --------------------- | ----------------------------: |
| CLI startup time      |   Under 100 ms where possible |
| Plain lookup overhead |                       Minimal |
| TUI initial render    |                  Under 200 ms |
| DNS timeout default   |                     3 seconds |
| Memory usage          |    Lightweight for simple TUI |
| Binary distribution   | Single executable per OS/arch |

The DNS network request will usually take longer than Cobra, Bubble Tea, or Lip Gloss overhead. The app should avoid unnecessary startup work before resolving.

---

## 18. Cross-Platform Requirements

## 18.1 Supported Platforms

| Platform      | Required |
| ------------- | -------- |
| Arch Linux    | Yes      |
| Ubuntu/Debian | Yes      |
| Windows 10/11 | Yes      |
| macOS         | Optional |

---

## 18.2 Shell Support

| Shell      | Support       |
| ---------- | ------------- |
| Bash       | Required      |
| Zsh        | Required      |
| Fish       | Optional      |
| PowerShell | Required      |
| CMD        | Basic support |

---

## 18.3 Distribution Targets

| Platform | Primary Distribution               | One-Command Install               | Future Package Manager      |
| -------- | ---------------------------------- | --------------------------------- | --------------------------- |
| Linux    | GitHub Releases `.tar.gz`          | `curl -fsSL .../install.sh \| sh` | `.deb`, AUR, Homebrew Linux |
| Windows  | GitHub Releases `.zip` with `.exe` | `irm .../install.ps1 \| iex`      | Scoop, Winget               |
| macOS    | GitHub Releases `.tar.gz`          | `curl -fsSL .../install.sh \| sh` | Homebrew                    |

Installer Responsibilities

The installer script should:

Detect operating system.
Detect CPU architecture.
Download the correct release artifact from GitHub Releases.
Extract the binary.
Install it into a user-writable binary directory.
Verify installation with aqry --version.
Print the installed path.
Show next command to run.
Linux/macOS Install Path

Default install location:

/usr/local/bin/aqry

If permission is denied, fallback to:

~/.local/bin/aqry
Windows Install Path

Default install location:

$env:LOCALAPPDATA\Programs\aqry\aqry.exe

The installer should add this directory to the user PATH when possible.

Release Asset Naming

GitHub Releases should publish predictable artifact names:

aqry_Linux_x86_64.tar.gz
aqry_Linux_arm64.tar.gz
aqry_Darwin_x86_64.tar.gz
aqry_Darwin_arm64.tar.gz
aqry_Windows_x86_64.zip

---

## 19. MVP Scope

## 19.1 MVP Features

Required for first release:

- `aqry <domain>` returns first IPv4.
- `aqry` opens interactive TUI.
- Domain input.
- A record lookup.
- Record type selector for:
  - A
  - AAAA
  - CNAME
  - MX
  - TXT
  - NS

- Loading spinner.
- Phase-based progress bar.
- Help menu.
- Error states.
- Keyboard navigation.
- Polished Lip Gloss styling.
- Cross-platform build.

---

## 19.2 MVP Exclusions

Not required in MVP:

- Persistent history
- Config file
- DNSSEC
- WHOIS
- Export to file
- Saved favorite domains
- Plugin system
- Auto-update system

---

## 20. Future Features

## 20.1 v1.1

- Lookup history
- Copy selected result
- JSON output
- Resolver selector
- Custom resolver input
- `--all`
- `--type`

---

## 20.2 v1.2

- DNS propagation comparison
- Query multiple DNS providers at once
- Result diff view
- Save favorite domains
- Config file support

---

## 20.3 v2.0

- DNSSEC validation
- WHOIS lookup
- HTTP status lookup
- Certificate inspection
- Export report
- Plugin architecture

---

## 21. Success Metrics

| Metric                                    | Target        |
| ----------------------------------------- | ------------- |
| Plain lookup works                        | 100% required |
| TUI starts without error                  | 100% required |
| A record lookup success                   | Required      |
| User can select record type without flags | Required      |
| Help menu accessible                      | Required      |
| Works on Linux and Windows                | Required      |
| UI usable at 80x24 terminal size          | Required      |
| No crash on invalid domain                | Required      |

---

## 22. Acceptance Criteria

## 22.1 CLI Acceptance

- Running `aqry google.com` prints only one IPv4 by default.
- Running `aqry google.com --all` prints all IPv4 records.
- Running `aqry google.com -t MX` prints MX records.
- Invalid domains return clear errors.
- Exit codes are correct:
  - `0` success
  - `1` validation error
  - `2` DNS lookup error
  - `3` timeout
  - `4` unsupported record type

---

## 22.2 TUI Acceptance

- Running `aqry` launches TUI.
- User can type a domain.
- User can press `Enter` to lookup.
- Spinner appears during lookup.
- Progress bar updates through phases.
- Result panel shows IPv4 by default.
- User can press `r` to choose record type.
- User can press `?` to open help.
- User can press `q` to quit.
- UI handles terminal resize.
- UI does not break in compact terminal widths.

---

## 22.3 UI Quality Acceptance

The UI should not look like a basic centered box only.

It must include:

- Header
- Main content area
- Focus styling
- Footer help bar
- Result panel
- Modal selector
- Loading state
- Empty state
- Error state
- Responsive layout
- Consistent theme tokens

---

## 23. Example Final UX

## 23.1 Quick Mode

```text
$ aqry google.com
142.250.182.46
```

---

## 23.2 Interactive Mode

```text
$ aqry
```

```text
╭────────────────────────────────────────────────────────────╮
│ aqry                                        DNS Lookup │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  Domain                                                    │
│  ╭──────────────────────────────────────────────────────╮  │
│  │ google.com                                           │  │
│  ╰──────────────────────────────────────────────────────╯  │
│                                                            │
│  Record Type                                               │
│  [ A ]  AAAA  CNAME  MX  TXT  NS  SOA                     │
│                                                            │
│  Resolver                                                  │
│  System DNS                                                │
│                                                            │
│  Result                                                    │
│  ╭─ IPv4 ───────────────────────────────────────────────╮  │
│  │ 142.250.182.46                                      │  │
│  ╰──────────────────────────────────────────────────────╯  │
│                                                            │
├────────────────────────────────────────────────────────────┤
│ Enter lookup  r records  s resolver  ? help  q quit        │
╰────────────────────────────────────────────────────────────╯
```

---

## 24. Recommended Development Milestones

## Milestone 1 — CLI Core

- Create Go module.
- Add Cobra root command.
- Implement `aqry <domain>`.
- Resolve A record.
- Print first IPv4 only.
- Add error handling.

---

## Milestone 2 — DNS Engine

- Add support for:
  - A
  - AAAA
  - CNAME
  - MX
  - TXT
  - NS

- Add structured record result model.
- Add timeout handling.
- Add resolver abstraction.

---

## Milestone 3 — Bubble Tea MVP

- Create model/update/view.
- Add text input.
- Add domain validation.
- Add lookup command.
- Add loading spinner.
- Add result view.

---

## Milestone 4 — Modern UI Layer

- Add Lip Gloss theme.
- Add header/footer layout.
- Add result panels.
- Add responsive terminal layout.
- Add selected/focused styling.

---

## Milestone 5 — Navigation and Modals

- Add record type picker.
- Add resolver picker.
- Add help menu.
- Add keyboard shortcuts.
- Add compact layout behavior.

---

## Milestone 6 — Polish

- Add progress bar.
- Add animations.
- Add better empty/error states.
- Add copy result.
- Add tests.
- Add CI builds.
- Add release packaging.

---

## 25. Technical Risks

| Risk                                   | Impact | Mitigation                              |
| -------------------------------------- | ------ | --------------------------------------- |
| Windows terminal rendering differences | Medium | Test in Windows Terminal and PowerShell |
| DNS results vary by network            | Medium | Allow resolver selection                |
| Terminal too small                     | Medium | Add compact layout                      |
| Too much UI complexity                 | Medium | Keep quick CLI mode simple              |
| Slow DNS response                      | Low    | Use timeout and loading states          |

---

## 26. Product Positioning

`aqry` should sit between basic lookup tools and advanced DNS debugging tools.

| Tool       | Position                                 |
| ---------- | ---------------------------------------- |
| `nslookup` | Basic DNS lookup                         |
| `dig`      | Advanced DNS debugging                   |
| `host`     | Simple DNS query                         |
| `aqry`     | Modern interactive DNS lookup for humans |

The key product advantage is not raw DNS power. The advantage is a better terminal workflow.

---

## 27. Final Product Requirement Summary

`aqry` must be:

- Fast by default.
- IPv4-first.
- Keyboard-driven.
- Interactive when needed.
- Beautiful but not distracting.
- Useful without memorizing flags.
- Cross-platform.
- Script-friendly.
- Designed like a real terminal product, not a basic demo UI.

The ideal experience:

```text
aqry google.com
```

Returns:

```text
142.250.182.46
```

And:

```text
aqry
```

Opens a polished interactive DNS resolver where users can search domains, choose record types, see loading states, navigate with the keyboard, open help, and inspect results cleanly.
