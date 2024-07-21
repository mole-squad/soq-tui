package settings

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mole-squad/soq-tui/pkg/api"
	"github.com/mole-squad/soq-tui/pkg/common"
	"github.com/mole-squad/soq-tui/pkg/styles"
)

type SettingsModel struct {
	client *api.Client

	keys keyMap
	help help.Model

	width int
}

func NewSettingsModel(client *api.Client) common.AppView {
	return SettingsModel{
		client: client,
		help:   help.New(),
		keys:   newKeyMap(),
	}
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.onWindowMsg(msg)

	case tea.KeyMsg:
		return m.onKeyMsg(msg)
	}

	return m, nil
}

func (m SettingsModel) View() string {
	sections := make([]string, 0)

	helpContent := m.help.View(m.keys)

	sections = append(sections, lipgloss.NewStyle().Width(m.width).Render(helpContent))

	return styles.PageWrapperStyle.Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (m SettingsModel) Blur() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m SettingsModel) Focus() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m SettingsModel) onWindowMsg(msg tea.WindowSizeMsg) (SettingsModel, tea.Cmd) {
	m.width = msg.Width

	m.help.Width = msg.Width

	return m, nil
}

func (m SettingsModel) onKeyMsg(msg tea.KeyMsg) (SettingsModel, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		return m, common.AppStateCmd(common.AppStateTaskList)

	case key.Matches(msg, m.keys.FocusAreas):
		return m, common.AppStateCmd(common.AppStateFocusAreaList)
	}

	return m, nil
}
