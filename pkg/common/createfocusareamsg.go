package common

import tea "github.com/charmbracelet/bubbletea"

type CreateFocusAreaMsg struct {
}

func NewCreateFocusAreaMsg() tea.Cmd {
	return func() tea.Msg {
		return CreateFocusAreaMsg{}
	}
}
