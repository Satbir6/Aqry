package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"aqry/internal/dns"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		m.height = message.Height
		m.resizeComponents()
		return m, nil

	case lookupProgressMsg:
		if m.state == StateLoading && message.id == m.lookupID {
			m.loadingMessage = message.message
			m.progress = message.percent
		}
		return m, nil

	case lookupResponseMsg:
		if message.id != m.lookupID || m.state != StateLoading {
			return m, nil
		}
		m.loadingMessage = "Parsing DNS response..."
		m.progress = 0.85
		return m, finalizeLookupCmd(message.id, message.result)

	case lookupSuccessMsg:
		if message.id != m.lookupID || m.state != StateLoading {
			return m, nil
		}
		m.finishLookup()
		m.state = StateSuccess
		m.result = message.result
		m.selectedRecord = 0
		m.err = nil
		m.progress = 1
		m.loadingMessage = "Resolved successfully"
		m.status = fmt.Sprintf("Success · %d record(s)", len(message.result.Records))
		m.focus = FocusResults
		m.input.Blur()
		return m, nil

	case lookupErrorMsg:
		if message.id != m.lookupID || m.state != StateLoading {
			return m, nil
		}
		m.finishLookup()
		m.state = StateError
		m.err = message.err
		m.progress = 0
		m.loadingMessage = "Lookup failed"
		m.status = "Error"
		return m, nil

	case copySuccessMsg:
		m.status = "Copied selected result"
		return m, nil

	case copyErrorMsg:
		m.status = "Copy failed: " + message.err.Error()
		return m, nil

	case spinner.TickMsg:
		if m.state == StateLoading {
			var command tea.Cmd
			m.spinner, command = m.spinner.Update(message)
			return m, command
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(message)
	}

	if m.modal == ModalNone && m.focus == FocusDomainInput && m.state != StateLoading {
		var command tea.Cmd
		m.input, command = m.input.Update(message)
		return m, command
	}

	return m, nil
}

func (m Model) handleKey(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.String() == "ctrl+c" {
		m.cancelLookup()
		return m, tea.Quit
	}
	if m.modal != ModalNone {
		return m.handleModalKey(key)
	}

	if m.focus == FocusDomainInput && m.state != StateLoading {
		switch key.String() {
		case "enter":
			return m.startLookup()
		case "tab":
			m.moveFocus(1)
			return m, nil
		case "shift+tab":
			m.moveFocus(-1)
			return m, nil
		case "?":
			m.modal = ModalHelp
			m.input.Blur()
			return m, nil
		}

		var command tea.Cmd
		m.input, command = m.input.Update(key)
		m.status = "Ready"
		return m, command
	}

	switch key.String() {
	case "q":
		m.cancelLookup()
		return m, tea.Quit
	case "?":
		m.modal = ModalHelp
		return m, nil
	case "esc":
		if m.state == StateLoading {
			m.cancelLookup()
			m.state = StateIdle
			m.status = "Lookup cancelled"
			m.progress = 0
			return m, nil
		}
	case "tab":
		m.moveFocus(1)
		return m, nil
	case "shift+tab":
		m.moveFocus(-1)
		return m, nil
	case "enter":
		if m.focus == FocusRecordType {
			m.syncPickerIndexes()
			m.modal = ModalRecordPicker
			return m, nil
		}
		if m.focus == FocusResolver {
			m.syncPickerIndexes()
			m.modal = ModalResolverPicker
			return m, nil
		}
		return m.startLookup()
	case "r":
		m.syncPickerIndexes()
		m.modal = ModalRecordPicker
		return m, nil
	case "s":
		m.syncPickerIndexes()
		m.modal = ModalResolverPicker
		return m, nil
	case "a":
		if m.state == StateSuccess && len(m.result.Records) > 1 {
			m.showAll = !m.showAll
			if m.showAll {
				m.status = "Showing all results"
			} else {
				m.status = "Showing selected result"
			}
		}
		return m, nil
	case "c":
		if record, ok := m.selectedResult(); ok {
			return m, copyCmd(m.clipboard, recordDisplayValue(record))
		}
	case "j", "down":
		m.moveResult(1)
		return m, nil
	case "k", "up":
		m.moveResult(-1)
		return m, nil
	}

	return m, nil
}

func (m Model) handleModalKey(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.String() == "q" && m.modal != ModalCustomResolver {
		m.cancelLookup()
		return m, tea.Quit
	}

	switch m.modal {
	case ModalHelp:
		if key.String() == "esc" || key.String() == "?" {
			m.closeModal()
		}
		return m, nil

	case ModalRecordPicker:
		recordTypes := dns.SupportedRecordTypes()
		switch key.String() {
		case "esc":
			m.closeModal()
		case "j", "down":
			m.recordPickerIndex = wrapIndex(m.recordPickerIndex+1, len(recordTypes))
		case "k", "up":
			m.recordPickerIndex = wrapIndex(m.recordPickerIndex-1, len(recordTypes))
		case "enter":
			m.recordType = recordTypes[m.recordPickerIndex]
			m.status = "Record type: " + string(m.recordType)
			m.closeModal()
		}
		return m, nil

	case ModalResolverPicker:
		switch key.String() {
		case "esc":
			m.closeModal()
		case "j", "down":
			m.resolverPickerIndex = wrapIndex(m.resolverPickerIndex+1, len(resolverChoices))
		case "k", "up":
			m.resolverPickerIndex = wrapIndex(m.resolverPickerIndex-1, len(resolverChoices))
		case "enter":
			choice := resolverChoices[m.resolverPickerIndex]
			if choice.Address == "" {
				m.modal = ModalCustomResolver
				m.customInput.SetValue("")
				m.customInput.Focus()
				return m, textinput.Blink
			}
			m.server = choice.Address
			m.status = "Resolver: " + choice.Name
			m.closeModal()
		}
		return m, nil

	case ModalCustomResolver:
		switch key.String() {
		case "esc":
			m.customInput.Blur()
			m.modal = ModalResolverPicker
			return m, nil
		case "enter":
			server := strings.TrimSpace(m.customInput.Value())
			if err := dns.ValidateResolverServer(server); err != nil {
				m.status = "Invalid resolver: " + err.Error()
				return m, nil
			}
			m.server = server
			m.status = "Resolver: " + server
			m.customInput.Blur()
			m.closeModal()
			return m, nil
		}

		var command tea.Cmd
		m.customInput, command = m.customInput.Update(key)
		return m, command
	}

	return m, nil
}

func (m Model) startLookup() (tea.Model, tea.Cmd) {
	if m.state == StateLoading {
		return m, nil
	}

	domain, err := dns.NormalizeAndValidateDomain(m.input.Value())
	if err != nil {
		m.state = StateError
		m.err = err
		m.status = "Invalid domain"
		return m, nil
	}
	m.input.SetValue(domain)
	m.cancelLookup()
	m.lookupID++
	lookupContext, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	m.state = StateLoading
	m.err = nil
	m.status = "Resolving"
	m.loadingMessage = "Checking domain..."
	m.progress = 0.20
	m.input.Blur()

	request := dns.LookupRequest{
		Domain: domain, RecordType: m.recordType, Server: m.server,
		Timeout: m.timeout, All: true,
	}
	return m, tea.Batch(
		m.spinner.Tick,
		progressCmd(m.lookupID, 20*time.Millisecond, "Preparing resolver...", 0.35),
		progressCmd(m.lookupID, 60*time.Millisecond, "Querying DNS resolver...", 0.60),
		lookupCmd(lookupContext, m.lookupID, m.resolver, request),
	)
}

func (m *Model) finishLookup() {
	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}
}

func (m *Model) cancelLookup() {
	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}
}

func (m *Model) moveFocus(delta int) {
	m.focus = FocusArea(wrapIndex(int(m.focus)+delta, 4))
	if m.focus == FocusDomainInput && m.state != StateLoading {
		m.input.Focus()
	} else {
		m.input.Blur()
	}
}

func (m *Model) moveResult(delta int) {
	if m.state != StateSuccess || len(m.result.Records) == 0 {
		return
	}
	m.selectedRecord = wrapIndex(m.selectedRecord+delta, len(m.result.Records))
	m.status = fmt.Sprintf("Result %d of %d", m.selectedRecord+1, len(m.result.Records))
}

func (m Model) selectedResult() (dns.Record, bool) {
	if m.state != StateSuccess || m.selectedRecord < 0 || m.selectedRecord >= len(m.result.Records) {
		return dns.Record{}, false
	}
	return m.result.Records[m.selectedRecord], true
}

func (m *Model) closeModal() {
	m.modal = ModalNone
	if m.focus == FocusDomainInput && m.state != StateLoading {
		m.input.Focus()
	}
}

func wrapIndex(index, length int) int {
	if length <= 0 {
		return 0
	}
	return (index%length + length) % length
}

func (m *Model) resizeComponents() {
	componentWidth := m.width - 12
	if m.width > 100 {
		componentWidth = (m.width / 2) - 10
	}
	if componentWidth < 12 {
		componentWidth = 12
	}
	m.input.Width = componentWidth
	m.progressBar.Width = componentWidth
}
