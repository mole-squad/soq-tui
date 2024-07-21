package forms

import tea "github.com/charmbracelet/bubbletea"

type SetFieldValueMsg struct {
	FieldID string
	FormID  string
	Value   string
}

func NewSetFieldValueCmd(formID string, fieldID string, value string) tea.Cmd {
	return func() tea.Msg {
		return SetFieldValueMsg{
			FormID:  formID,
			FieldID: fieldID,
			Value:   value,
		}
	}
}
