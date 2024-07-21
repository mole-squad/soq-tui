package settings

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mole-squad/soq-tui/pkg/api"
	"github.com/mole-squad/soq-tui/pkg/common"
	"github.com/mole-squad/soq-tui/pkg/logger"
	"github.com/mole-squad/soq-tui/pkg/styles"
)

type Model struct {
	client *api.Client
	logger *logger.Logger

	keys keyMap
	help help.Model

	width int
}

func New(logger *logger.Logger, client *api.Client) common.AppView {
	return Model{
		client: client,
		logger: logger,
		help:   help.New(),
		keys:   newKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.onWindowMsg(msg)

	case tea.KeyMsg:
		return m.onKeyMsg(msg)
	}

	return m, nil
}

func (m Model) View() string {
	sections := make([]string, 0)

	helpContent := m.help.View(m.keys)

	sections = append(sections, lipgloss.NewStyle().Width(m.width).Render(helpContent))

	return styles.PageWrapperStyle.Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (m Model) Blur() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) Focus() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) onWindowMsg(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.width = msg.Width

	m.help.Width = msg.Width

	return m, nil
}

func (m Model) onKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		return m, common.AppStateCmd(common.AppStateTaskList)

	case key.Matches(msg, m.keys.FocusAreas):
		return m, common.AppStateCmd(common.AppStateFocusAreaList)
	}

	return m, nil
}
