package forms

import tea "github.com/charmbracelet/bubbletea"

type SubmitFormMsg struct {
	FormID string
}

func NewSubmitFormCmd(formID string) tea.Cmd {
	return func() tea.Msg {
		return SubmitFormMsg{
			FormID: formID,
		}
	}
}
