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

func NewSettingsModel(client *api.Client) SettingsModel {
	return SettingsModel{
		client: client,
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
	docFrameWidth, _ := styles.PageWrapperStyle.GetFrameSize()

	sections := make([]string, 0)

	helpContent := m.help.View(m.keys)

	sections = append(sections, lipgloss.NewStyle().Width(m.width-docFrameWidth).Render(helpContent))

	return styles.PageWrapperStyle.Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (m SettingsModel) onWindowMsg(msg tea.WindowSizeMsg) (SettingsModel, tea.Cmd) {
	docFrameWidth, _ := styles.PageWrapperStyle.GetFrameSize()

	m.width = msg.Width

	m.help.Width = msg.Width - docFrameWidth

	return m, nil
}

func (m SettingsModel) onKeyMsg(msg tea.KeyMsg) (SettingsModel, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		return m, common.AppStateCmd(common.AppStateTaskList)
	}

	return m, nil
}
