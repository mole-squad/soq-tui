package forms

import tea "github.com/charmbracelet/bubbletea"

type SetFieldValueMsg struct {
	FieldID string
	Value   string
}

func NewSetFieldValueCmd(fieldID string, value string) tea.Cmd {
	return func() tea.Msg {
		return SetFieldValueMsg{
			FieldID: fieldID,
			Value:   value,
		}
	}
}
