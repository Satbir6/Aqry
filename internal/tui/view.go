package tui

func (m Model) View() string {
	if m.modal != ModalNone {
		return m.modalView()
	}
	if m.width < 60 {
		return m.compactView()
	}
	if m.width > 100 {
		return m.wideView()
	}
	return m.standardView()
}
