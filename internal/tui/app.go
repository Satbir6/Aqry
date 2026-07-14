package tui

import (
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

type systemClipboard struct{}

func (systemClipboard) WriteAll(value string) error {
	return clipboard.WriteAll(value)
}

func Run(options Options) error {
	if options.Clipboard == nil {
		options.Clipboard = systemClipboard{}
	}
	program := tea.NewProgram(NewModel(options), tea.WithAltScreen())
	_, err := program.Run()
	return err
}
