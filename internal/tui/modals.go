package tui

import (
	"fmt"
	"strings"

	"aqry/internal/dns"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) modalView() string {
	var title string
	var lines []string

	switch m.modal {
	case ModalHelp:
		title = "Help"
		lines = []string{
			"Enter       Run lookup / confirm",
			"Tab         Switch focus",
			"↑ / ↓, j/k  Navigate",
			"r           Select record type",
			"s           Select DNS resolver",
			"a           Toggle all results",
			"c           Copy selected result",
			"Esc         Close / cancel lookup",
			"q / Ctrl+C  Quit",
			"",
			m.styles.Muted.Render("Esc or ? close"),
		}

	case ModalRecordPicker:
		title = "Select Record Type"
		for index, recordType := range dns.SupportedRecordTypes() {
			line := fmt.Sprintf("  %-6s %s", recordType, recordType.Description())
			if index == m.recordPickerIndex {
				line = m.styles.SelectedItem.Render("› " + strings.TrimLeft(line, " "))
			}
			lines = append(lines, line)
		}
		lines = append(lines, "", m.styles.Muted.Render("↑/↓ navigate · Enter select · Esc cancel"))

	case ModalResolverPicker:
		title = "DNS Resolver"
		for index, choice := range resolverChoices {
			address := choice.Address
			if address == "system" {
				address = "OS default"
			}
			line := fmt.Sprintf("  %-14s %s", choice.Name, address)
			if index == m.resolverPickerIndex {
				line = m.styles.SelectedItem.Render("› " + strings.TrimLeft(line, " "))
			}
			lines = append(lines, line)
		}
		lines = append(lines, "", m.styles.Muted.Render("↑/↓ navigate · Enter select · Esc cancel"))

	case ModalCustomResolver:
		title = "Custom DNS Resolver"
		lines = []string{
			"Enter an IPv4 or IPv6 resolver address:",
			"",
			m.customInput.View(),
			"",
			m.styles.Muted.Render(m.status),
			m.styles.Muted.Render("Enter save · Esc back"),
		}
	}

	content := m.styles.ModalTitle.Render(title) + "\n" + strings.Join(lines, "\n")
	modalWidth := 54
	if m.width-6 < modalWidth {
		modalWidth = m.width - 6
	}
	if modalWidth < 24 {
		modalWidth = 24
	}
	modal := m.styles.Modal.Width(modalWidth).Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
}
