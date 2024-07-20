package common

import tea "github.com/charmbracelet/bubbletea"

type AppView interface {
	tea.Model

	Blur() (tea.Model, tea.Cmd)
	Focus() (tea.Model, tea.Cmd)
}
