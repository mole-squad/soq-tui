package forms

import tea "github.com/charmbracelet/bubbletea"

type FormField interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)

	View() string
	ViewSidePanel() string

	Blur() tea.Cmd
	Focus() tea.Cmd

	HasPanelContent() bool

	GetID() string

	GetValue() string
	SetValue(string)

	SetSize(width int, height int)
	SetPanelSize(width int, height int)
}
