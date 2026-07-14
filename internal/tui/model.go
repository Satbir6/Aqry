package tui

import (
	"context"
	"strings"
	"time"

	"aqry/internal/dns"
	appstyles "aqry/internal/styles"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
)

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
	ModalCustomResolver
)

type Clipboard interface {
	WriteAll(string) error
}

type Options struct {
	Domain     string
	RecordType dns.RecordType
	Server     string
	Timeout    time.Duration
	NoColor    bool
	Resolver   dns.Resolver
	Clipboard  Clipboard
}

type resolverChoice struct {
	Name    string
	Address string
}

var resolverChoices = []resolverChoice{
	{Name: "System DNS", Address: "system"},
	{Name: "Cloudflare", Address: "1.1.1.1"},
	{Name: "Google", Address: "8.8.8.8"},
	{Name: "Quad9", Address: "9.9.9.9"},
	{Name: "Custom...", Address: ""},
}

type Model struct {
	width  int
	height int

	state AppState
	modal ModalType
	focus FocusArea

	recordType dns.RecordType
	server     string
	timeout    time.Duration

	result         dns.LookupResult
	selectedRecord int
	showAll        bool
	err            error
	status         string

	loadingMessage string
	progress       float64

	recordPickerIndex   int
	resolverPickerIndex int

	input       textinput.Model
	customInput textinput.Model
	spinner     spinner.Model
	progressBar progress.Model
	styles      appstyles.Set

	resolver  dns.Resolver
	clipboard Clipboard

	lookupID int
	cancel   context.CancelFunc
}

func NewModel(options Options) Model {
	recordType := options.RecordType
	if !recordType.Supported() {
		recordType = dns.RecordA
	}
	server := strings.TrimSpace(options.Server)
	if server == "" {
		server = "system"
	}
	timeout := options.Timeout
	if timeout <= 0 {
		timeout = dns.DefaultTimeout
	}
	resolver := options.Resolver
	if resolver == nil {
		resolver = dns.NewSystemResolver()
	}

	styleSet := appstyles.New(options.NoColor)
	domainInput := textinput.New()
	domainInput.Placeholder = "example.com"
	domainInput.CharLimit = 512
	domainInput.Width = 48
	domainInput.SetValue(options.Domain)
	domainInput.Focus()
	domainInput.Prompt = "› "
	domainInput.TextStyle = styleSet.Input
	domainInput.PromptStyle = styleSet.Primary

	customInput := textinput.New()
	customInput.Placeholder = "1.1.1.1"
	customInput.CharLimit = 64
	customInput.Width = 30
	customInput.Prompt = "› "
	customInput.TextStyle = styleSet.Input
	customInput.PromptStyle = styleSet.Primary

	loadingSpinner := spinner.New()
	loadingSpinner.Spinner = spinner.Dot
	loadingSpinner.Style = styleSet.Primary

	progressBar := progress.New(progress.WithDefaultGradient())
	progressBar.Width = 36
	if options.NoColor {
		progressBar.FullColor = ""
		progressBar.EmptyColor = ""
	}

	model := Model{
		width:          80,
		height:         24,
		state:          StateIdle,
		modal:          ModalNone,
		focus:          FocusDomainInput,
		recordType:     recordType,
		server:         server,
		timeout:        timeout,
		status:         "Ready",
		loadingMessage: "Ready",
		input:          domainInput,
		customInput:    customInput,
		spinner:        loadingSpinner,
		progressBar:    progressBar,
		styles:         styleSet,
		resolver:       resolver,
		clipboard:      options.Clipboard,
	}
	model.syncPickerIndexes()
	return model
}

func (m *Model) syncPickerIndexes() {
	for index, recordType := range dns.SupportedRecordTypes() {
		if recordType == m.recordType {
			m.recordPickerIndex = index
			break
		}
	}
	m.resolverPickerIndex = len(resolverChoices) - 1
	for index, choice := range resolverChoices {
		if strings.EqualFold(choice.Address, m.server) {
			m.resolverPickerIndex = index
			break
		}
	}
}
