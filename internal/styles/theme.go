package styles

import "github.com/charmbracelet/lipgloss"

type Set struct {
	NoColor bool
	Theme   Theme

	AppFrame         lipgloss.Style
	Header           lipgloss.Style
	HeaderTitle      lipgloss.Style
	HeaderMeta       lipgloss.Style
	Footer           lipgloss.Style
	Panel            lipgloss.Style
	FocusedPanel     lipgloss.Style
	Label            lipgloss.Style
	Muted            lipgloss.Style
	Primary          lipgloss.Style
	Success          lipgloss.Style
	Warning          lipgloss.Style
	Danger           lipgloss.Style
	Input            lipgloss.Style
	RecordPill       lipgloss.Style
	RecordPillActive lipgloss.Style
	FocusedValue     lipgloss.Style
	ResultValue      lipgloss.Style
	Modal            lipgloss.Style
	ModalTitle       lipgloss.Style
	SelectedItem     lipgloss.Style
}

func New(noColor bool) Set {
	theme := DefaultTheme()
	border := lipgloss.RoundedBorder()

	set := Set{
		NoColor:          noColor,
		Theme:            theme,
		AppFrame:         lipgloss.NewStyle().Border(border).Padding(0, 1),
		Header:           lipgloss.NewStyle().Bold(true),
		HeaderTitle:      lipgloss.NewStyle().Bold(true),
		HeaderMeta:       lipgloss.NewStyle(),
		Footer:           lipgloss.NewStyle().PaddingTop(1),
		Panel:            lipgloss.NewStyle().Border(border).Padding(0, 1),
		FocusedPanel:     lipgloss.NewStyle().Border(border).Padding(0, 1),
		Label:            lipgloss.NewStyle().Bold(true),
		Muted:            lipgloss.NewStyle(),
		Primary:          lipgloss.NewStyle().Bold(true),
		Success:          lipgloss.NewStyle().Bold(true),
		Warning:          lipgloss.NewStyle().Bold(true),
		Danger:           lipgloss.NewStyle().Bold(true),
		Input:            lipgloss.NewStyle(),
		RecordPill:       lipgloss.NewStyle().Padding(0, 1),
		RecordPillActive: lipgloss.NewStyle().Padding(0, 1).Bold(true),
		FocusedValue:     lipgloss.NewStyle().Padding(0, 1).Bold(true),
		ResultValue:      lipgloss.NewStyle().Bold(true),
		Modal:            lipgloss.NewStyle().Border(border).Padding(1, 2),
		ModalTitle:       lipgloss.NewStyle().Bold(true).MarginBottom(1),
		SelectedItem:     lipgloss.NewStyle().Padding(0, 1).Bold(true),
	}

	if noColor {
		set.FocusedPanel = set.FocusedPanel.BorderStyle(lipgloss.DoubleBorder())
		set.RecordPillActive = set.RecordPillActive.Reverse(true)
		set.FocusedValue = set.FocusedValue.Reverse(true)
		set.SelectedItem = set.SelectedItem.Reverse(true)
		return set
	}

	set.AppFrame = set.AppFrame.BorderForeground(lipgloss.Color(theme.Border))
	set.HeaderTitle = set.HeaderTitle.Foreground(lipgloss.Color(theme.Primary))
	set.HeaderMeta = set.HeaderMeta.Foreground(lipgloss.Color(theme.Muted))
	set.Footer = set.Footer.Foreground(lipgloss.Color(theme.Muted))
	set.Panel = set.Panel.BorderForeground(lipgloss.Color(theme.Border))
	set.FocusedPanel = set.FocusedPanel.BorderForeground(lipgloss.Color(theme.Primary))
	set.Label = set.Label.Foreground(lipgloss.Color(theme.Text))
	set.Muted = set.Muted.Foreground(lipgloss.Color(theme.Muted))
	set.Primary = set.Primary.Foreground(lipgloss.Color(theme.Primary))
	set.Success = set.Success.Foreground(lipgloss.Color(theme.Success))
	set.Warning = set.Warning.Foreground(lipgloss.Color(theme.Warning))
	set.Danger = set.Danger.Foreground(lipgloss.Color(theme.Danger))
	set.Input = set.Input.Foreground(lipgloss.Color(theme.Text))
	set.RecordPill = set.RecordPill.Foreground(lipgloss.Color(theme.Muted))
	set.RecordPillActive = set.RecordPillActive.
		Foreground(lipgloss.Color("#10131A")).Background(lipgloss.Color(theme.Primary))
	set.FocusedValue = set.FocusedValue.
		Foreground(lipgloss.Color("#10131A")).Background(lipgloss.Color(theme.Primary))
	set.ResultValue = set.ResultValue.Foreground(lipgloss.Color(theme.Text))
	set.Modal = set.Modal.BorderForeground(lipgloss.Color(theme.Primary))
	set.ModalTitle = set.ModalTitle.Foreground(lipgloss.Color(theme.Primary))
	set.SelectedItem = set.SelectedItem.
		Foreground(lipgloss.Color("#10131A")).Background(lipgloss.Color(theme.Primary))

	return set
}
