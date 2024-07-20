package forms

import tea "github.com/charmbracelet/bubbletea"

type SetSelectOptionsMsg struct {
	InputID string
	Options []SelectOption
}

func NewSetSelectOptionsCmd(inputID string, options []SelectOption) tea.Cmd {
	return func() tea.Msg {
		return SetSelectOptionsMsg{
			Options: options,
			InputID: inputID,
		}
	}
}
