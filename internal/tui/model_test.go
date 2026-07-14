package tui

import (
	"context"
	"errors"
	"strings"
	"testing"

	"aqry/internal/dns"

	tea "github.com/charmbracelet/bubbletea"
)

type tuiMockResolver struct {
	result dns.LookupResult
	err    error
}

func (r tuiMockResolver) Lookup(context.Context, dns.LookupRequest) (dns.LookupResult, error) {
	return r.result, r.err
}

type memoryClipboard struct {
	value string
	err   error
}

func (c *memoryClipboard) WriteAll(value string) error {
	c.value = value
	return c.err
}

func updateModel(t *testing.T, model Model, message tea.Msg) (Model, tea.Cmd) {
	t.Helper()
	updated, command := model.Update(message)
	result, ok := updated.(Model)
	if !ok {
		t.Fatalf("Update returned %T, want tui.Model", updated)
	}
	return result, command
}

func runeKey(value string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(value)}
}

func TestTypingUpdatesDomainInput(t *testing.T) {
	t.Parallel()

	model := NewModel(Options{NoColor: true})
	model, _ = updateModel(t, model, runeKey("query.example.com"))
	if model.input.Value() != "query.example.com" {
		t.Fatalf("input = %q", model.input.Value())
	}
}

func TestRecordPickerOpensNavigatesAndCloses(t *testing.T) {
	t.Parallel()

	model := NewModel(Options{NoColor: true})
	model, _ = updateModel(t, model, tea.KeyMsg{Type: tea.KeyTab})
	if model.focus != FocusRecordType {
		t.Fatalf("focus = %v", model.focus)
	}
	model, _ = updateModel(t, model, runeKey("r"))
	if model.modal != ModalRecordPicker {
		t.Fatalf("modal = %v", model.modal)
	}
	model, _ = updateModel(t, model, tea.KeyMsg{Type: tea.KeyDown})
	model, _ = updateModel(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	if model.modal != ModalNone || model.recordType != dns.RecordAAAA {
		t.Fatalf("modal = %v, record type = %s", model.modal, model.recordType)
	}

	model.modal = ModalRecordPicker
	model, _ = updateModel(t, model, tea.KeyMsg{Type: tea.KeyEsc})
	if model.modal != ModalNone {
		t.Fatalf("Esc left modal %v open", model.modal)
	}
}

func TestHelpModalToggles(t *testing.T) {
	t.Parallel()

	model := NewModel(Options{NoColor: true})
	model, _ = updateModel(t, model, runeKey("?"))
	if model.modal != ModalHelp {
		t.Fatalf("modal = %v", model.modal)
	}
	model, _ = updateModel(t, model, tea.KeyMsg{Type: tea.KeyEsc})
	if model.modal != ModalNone || model.focus != FocusDomainInput {
		t.Fatalf("modal = %v, focus = %v", model.modal, model.focus)
	}
}

func TestLookupStateTransitions(t *testing.T) {
	t.Parallel()

	model := NewModel(Options{NoColor: true, Resolver: tuiMockResolver{}})
	model.input.SetValue("https://Example.com/path")
	model, command := updateModel(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	if command == nil || model.state != StateLoading {
		t.Fatalf("state = %v, command nil = %v", model.state, command == nil)
	}
	if model.input.Value() != "example.com" {
		t.Fatalf("normalized input = %q", model.input.Value())
	}

	success := dns.LookupResult{
		Domain: "example.com", RecordType: dns.RecordA, Resolver: "system",
		Records: []dns.Record{{Type: dns.RecordA, Value: "192.0.2.1"}},
	}
	model, _ = updateModel(t, model, lookupSuccessMsg{id: model.lookupID, result: success})
	if model.state != StateSuccess || model.focus != FocusResults || len(model.result.Records) != 1 {
		t.Fatalf("success model = %#v", model)
	}
}

func TestLookupResponseUsesParsingPhaseBeforeSuccess(t *testing.T) {
	t.Parallel()

	model := NewModel(Options{NoColor: true})
	model.state = StateLoading
	model.lookupID = 3
	result := dns.LookupResult{
		Domain: "example.com", RecordType: dns.RecordA,
		Records: []dns.Record{{Type: dns.RecordA, Value: "192.0.2.1"}},
	}
	model, command := updateModel(t, model, lookupResponseMsg{id: 3, result: result})
	if model.state != StateLoading || model.progress != 0.85 || command == nil {
		t.Fatalf("state = %v, progress = %.2f, command nil = %v", model.state, model.progress, command == nil)
	}
	model, _ = updateModel(t, model, command())
	if model.state != StateSuccess || model.progress != 1 {
		t.Fatalf("state = %v, progress = %.2f", model.state, model.progress)
	}
}

func TestLookupErrorTransition(t *testing.T) {
	t.Parallel()

	model := NewModel(Options{NoColor: true})
	model.input.SetValue("example.com")
	model, _ = updateModel(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	expected := dns.NewLookupError(dns.ErrNoRecords, "example.com", dns.RecordA, nil)
	model, _ = updateModel(t, model, lookupErrorMsg{id: model.lookupID, err: expected})
	if model.state != StateError || !errors.Is(model.err, expected) {
		t.Fatalf("state = %v, error = %v", model.state, model.err)
	}
}

func TestInvalidDomainDoesNotStartLookup(t *testing.T) {
	t.Parallel()

	model := NewModel(Options{NoColor: true})
	model.input.SetValue("bad domain")
	model, command := updateModel(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	if command != nil || model.state != StateError || !dns.IsErrorKind(model.err, dns.ErrInvalidDomain) {
		t.Fatalf("state = %v, command = %v, error = %v", model.state, command, model.err)
	}
}

func TestResultNavigationAndShowAll(t *testing.T) {
	t.Parallel()

	model := NewModel(Options{NoColor: true})
	model.focus = FocusResults
	model.input.Blur()
	model.state = StateSuccess
	model.result = dns.LookupResult{Records: []dns.Record{
		{Type: dns.RecordA, Value: "192.0.2.1"},
		{Type: dns.RecordA, Value: "192.0.2.2"},
	}}
	model, _ = updateModel(t, model, runeKey("j"))
	if model.selectedRecord != 1 {
		t.Fatalf("selected record = %d", model.selectedRecord)
	}
	model, _ = updateModel(t, model, runeKey("a"))
	if !model.showAll {
		t.Fatal("show-all did not toggle")
	}
}

func TestCopySelectedResult(t *testing.T) {
	t.Parallel()

	clipboard := &memoryClipboard{}
	model := NewModel(Options{NoColor: true, Clipboard: clipboard})
	model.focus = FocusResults
	model.input.Blur()
	model.state = StateSuccess
	model.result = dns.LookupResult{Records: []dns.Record{
		{Type: dns.RecordMX, Value: "mail.example.com", Priority: 10},
	}}
	model, command := updateModel(t, model, runeKey("c"))
	if command == nil {
		t.Fatal("copy key returned no command")
	}
	message := command()
	model, _ = updateModel(t, model, message)
	if clipboard.value != "10 mail.example.com" || !strings.Contains(model.status, "Copied") {
		t.Fatalf("clipboard = %q, status = %q", clipboard.value, model.status)
	}
}

func TestCustomResolverPicker(t *testing.T) {
	t.Parallel()

	model := NewModel(Options{NoColor: true})
	model.focus = FocusResolver
	model.input.Blur()
	model.modal = ModalResolverPicker
	model.resolverPickerIndex = len(resolverChoices) - 1
	model, _ = updateModel(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	if model.modal != ModalCustomResolver {
		t.Fatalf("modal = %v", model.modal)
	}
	model, _ = updateModel(t, model, runeKey("1.1.1.1"))
	model, _ = updateModel(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	if model.modal != ModalNone || model.server != "1.1.1.1" {
		t.Fatalf("modal = %v, server = %q", model.modal, model.server)
	}
}

func TestResponsiveViews(t *testing.T) {
	t.Parallel()

	for _, width := range []int{48, 80, 120} {
		model := NewModel(Options{NoColor: true})
		model.width = width
		view := model.View()
		if !strings.Contains(view, "aqry") || !strings.Contains(view, "Domain") {
			t.Errorf("width %d view missing core content: %q", width, view)
		}
	}
}
