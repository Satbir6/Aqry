package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) compactView() string {
	lines := []string{
		m.styles.HeaderTitle.Render("aqry"),
		m.styles.Muted.Render("DNS Resolver"),
		"",
		m.styles.Label.Render("Domain"),
		m.input.View(),
		"",
		m.styles.Label.Render("Type · " + string(m.recordType)),
		m.styles.Muted.Render(m.resolverLabel()),
		"",
		m.resultContent(),
		"",
		m.footer(),
	}
	return strings.Join(lines, "\n")
}

func (m Model) standardView() string {
	frameWidth := m.width - 6
	if frameWidth < 42 {
		frameWidth = 42
	}
	panelWidth := frameWidth - 4
	inputStyle := m.styles.Panel
	if m.focus == FocusDomainInput || m.focus == FocusRecordType || m.focus == FocusResolver {
		inputStyle = m.styles.FocusedPanel
	}
	resultStyle := m.styles.Panel
	if m.focus == FocusResults {
		resultStyle = m.styles.FocusedPanel
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		m.header(panelWidth),
		"",
		inputStyle.Width(panelWidth).Render(m.inputContent()),
		"",
		resultStyle.Width(panelWidth).Render(m.resultContent()),
		m.footer(),
	)
	return m.styles.AppFrame.Width(frameWidth).Render(content)
}

func (m Model) wideView() string {
	frameWidth := m.width - 6
	if frameWidth < 96 {
		frameWidth = 96
	}
	columnWidth := (frameWidth - 7) / 2
	inputStyle := m.styles.Panel
	if m.focus != FocusResults {
		inputStyle = m.styles.FocusedPanel
	}
	resultStyle := m.styles.Panel
	if m.focus == FocusResults {
		resultStyle = m.styles.FocusedPanel
	}
	columns := lipgloss.JoinHorizontal(lipgloss.Top,
		inputStyle.Width(columnWidth).Render(m.inputContent()),
		" ",
		resultStyle.Width(columnWidth).Render(m.resultContent()),
	)
	content := lipgloss.JoinVertical(lipgloss.Left,
		m.header(frameWidth-3),
		"",
		columns,
		m.footer(),
	)
	return m.styles.AppFrame.Width(frameWidth).Render(content)
}
