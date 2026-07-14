package tui

import (
	"fmt"
	"strings"

	"aqry/internal/dns"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) header(width int) string {
	title := m.styles.HeaderTitle.Render("aqry")
	mode := m.styles.HeaderMeta.Render("DNS Resolver")
	gap := width - lipgloss.Width(title) - lipgloss.Width(mode)
	if gap < 1 {
		gap = 1
	}
	return title + strings.Repeat(" ", gap) + mode
}

func (m Model) footer() string {
	switch m.focus {
	case FocusDomainInput:
		return m.styles.Footer.Render("enter lookup  tab options  ? help  ctrl+c quit")
	case FocusRecordType:
		return m.styles.Footer.Render("enter choose record  tab next  ? help  q quit")
	case FocusResolver:
		return m.styles.Footer.Render("enter choose resolver  tab next  ? help  q quit")
	case FocusResults:
		if m.state == StateSuccess {
			return m.styles.Footer.Render("↑↓ select  c copy  a all  tab input  ? help  q quit")
		}
		return m.styles.Footer.Render("enter lookup  tab input  ? help  q quit")
	default:
		return m.styles.Footer.Render("tab navigate  ? help  q quit")
	}
}

func (m Model) focusLabel(area FocusArea, label string) string {
	if m.focus == area {
		return m.styles.Primary.Render("› " + label)
	}
	return m.styles.Label.Render("  " + label)
}

func (m Model) focusedValue(area FocusArea, value string) string {
	if m.focus == area {
		return m.styles.FocusedValue.Render(value)
	}
	return value
}

func (m Model) recordPills() string {
	items := make([]string, 0, len(dns.SupportedRecordTypes()))
	for _, recordType := range dns.SupportedRecordTypes() {
		if recordType == m.recordType {
			items = append(items, m.styles.RecordPillActive.Render(string(recordType)))
		} else {
			items = append(items, m.styles.RecordPill.Render(string(recordType)))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, items...)
}

func (m Model) resolverLabel() string {
	for _, choice := range resolverChoices {
		if choice.Address != "" && strings.EqualFold(choice.Address, m.server) {
			if choice.Address == "system" {
				return choice.Name
			}
			return fmt.Sprintf("%s · %s", choice.Name, choice.Address)
		}
	}
	return "Custom · " + m.server
}

func (m Model) inputContent() string {
	return strings.Join([]string{
		m.focusLabel(FocusDomainInput, "Domain"),
		m.input.View(),
		"",
		m.focusLabel(FocusRecordType, "Record Type"),
		m.recordPills(),
		"",
		m.focusLabel(FocusResolver, "Resolver"),
		m.focusedValue(FocusResolver, m.resolverLabel()),
		"",
		m.styles.Muted.Render("Status · " + m.status),
	}, "\n")
}

func (m Model) resultContent() string {
	switch m.state {
	case StateLoading:
		return strings.Join([]string{
			m.styles.Label.Render("Resolving"),
			m.spinner.View() + " " + m.loadingMessage,
			"",
			m.progressBar.ViewAs(m.progress),
		}, "\n")

	case StateSuccess:
		return m.successContent()

	case StateError:
		reason := "Enter a domain to continue"
		if m.err != nil {
			reason = m.err.Error()
		}
		return strings.Join([]string{
			m.styles.Danger.Render("✕ Could not resolve domain"),
			"",
			m.styles.Label.Render("Reason"),
			reason,
			"",
			m.styles.Muted.Render("Check spelling · try another record type · try another resolver"),
		}, "\n")

	default:
		return strings.Join([]string{
			m.styles.Label.Render("Ready for a lookup"),
			"",
			m.styles.Muted.Render("Enter a domain and press Enter."),
			m.styles.Muted.Render("A records are selected by default."),
		}, "\n")
	}
}

func (m Model) successContent() string {
	lines := []string{
		m.styles.Success.Render("✓ Resolved successfully"),
		m.styles.Muted.Render(fmt.Sprintf("%s · %s record · %s", m.result.Domain, m.result.RecordType, m.resolverLabel())),
		"",
	}
	if m.showAll {
		for index, record := range m.result.Records {
			marker := "  "
			if index == m.selectedRecord {
				marker = "› "
			}
			line := fmt.Sprintf("%s%d  %s", marker, index+1, recordDisplayValue(record))
			if index == m.selectedRecord {
				line = m.styles.Primary.Render(line)
			}
			lines = append(lines, line)
		}
	} else if record, ok := m.selectedResult(); ok {
		lines = append(lines, m.styles.ResultValue.Render(recordDisplayValue(record)))
		if len(m.result.Records) > 1 {
			lines = append(lines, "", m.styles.Muted.Render(
				fmt.Sprintf("Result %d of %d · use ↑/↓ to navigate", m.selectedRecord+1, len(m.result.Records)),
			))
		}
	}
	return strings.Join(lines, "\n")
}

func recordDisplayValue(record dns.Record) string {
	if record.Type == dns.RecordMX {
		return fmt.Sprintf("%d %s", record.Priority, record.Value)
	}
	return record.Value
}
