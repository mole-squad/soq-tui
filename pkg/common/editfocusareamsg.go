package common

import (
	tea "github.com/charmbracelet/bubbletea"
	soqapi "github.com/mole-squad/soq-api/api"
)

type EditFocusAreaMsg struct {
	FocusArea soqapi.FocusAreaDTO
}

func NewEditFocusAreaMsg(focusArea soqapi.FocusAreaDTO) tea.Cmd {
	return func() tea.Msg {
		return EditFocusAreaMsg{FocusArea: focusArea}
	}
}
