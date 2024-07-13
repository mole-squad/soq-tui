package common

import (
	tea "github.com/charmbracelet/bubbletea"
	soqapi "github.com/mole-squad/soq-api/api"
)

type SelectTaskMsg struct {
	Task soqapi.TaskDTO
}

func NewSelectTaskMsg(task soqapi.TaskDTO) tea.Cmd {
	return func() tea.Msg {
		return SelectTaskMsg{Task: task}
	}
}
