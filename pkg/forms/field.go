package forms

import tea "github.com/charmbracelet/bubbletea"

type FormField interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string

	Blur() tea.Cmd
	Focus() tea.Cmd

	GetID() string

	GetValue() string
	SetValue(string)

	SetWidth(int)
	SetHeight(int)
}
