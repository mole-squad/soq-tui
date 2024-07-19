package forms

import tea "github.com/charmbracelet/bubbletea"

type SubmitFormMsg struct{}

func NewSubmitFormMsg() tea.Msg {
	return SubmitFormMsg{}
}
