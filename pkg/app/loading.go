package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mole-squad/soq-tui/pkg/common"
)

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
)

type LoadingModel struct {
	spinner spinner.Model

	height int
	width  int
}

func NewLoadingModel() common.AppView {
	teaSpinner := spinner.New()
	teaSpinner.Style = spinnerStyle
	teaSpinner.Spinner = spinner.Points

	return LoadingModel{
		spinner: teaSpinner,
	}
}

func (m LoadingModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m LoadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m LoadingModel) View() string {
	content := fmt.Sprintf("%s", m.spinner.View())

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}

func (m LoadingModel) Blur() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m LoadingModel) Focus() (tea.Model, tea.Cmd) {
	return m, nil
}
