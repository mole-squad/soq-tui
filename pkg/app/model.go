package app

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mole-squad/soq-tui/pkg/api"
	"github.com/mole-squad/soq-tui/pkg/common"
	"github.com/mole-squad/soq-tui/pkg/loginform"
	"github.com/mole-squad/soq-tui/pkg/settings"
	"github.com/mole-squad/soq-tui/pkg/styles"
	"github.com/mole-squad/soq-tui/pkg/taskform"
	"github.com/mole-squad/soq-tui/pkg/tasklist"
	"github.com/mole-squad/soq-tui/pkg/utils"
)

var (
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type AppModel struct {
	client *api.Client

	appState common.AppState
	error

	loginForm loginform.LoginFormModel
	taskForm  taskform.TaskFormModel
	taskList  tasklist.TaskListModel
	settings  settings.SettingsModel

	keys keyMap

	quitting bool
	width    int
}

func NewAppModel() AppModel {
	client := api.NewClient()

	err := client.LoadToken()
	if err != nil {
		slog.Error("failed to load token", "error", err)
	}

	return AppModel{
		appState:  common.AppStateLogin,
		client:    client,
		keys:      newKeyMap(),
		loginForm: loginform.NewLoginFormModel(client),
		taskForm:  taskform.NewTaskFormModel(client),
		taskList:  tasklist.NewTaskListModel(client),
		settings:  settings.NewSettingsModel(client),
	}
}

func (m AppModel) Init() tea.Cmd {
	initCmd := tea.Batch(
		m.loginForm.Init(),
		m.taskList.Init(),
		m.taskForm.Init(),
	)

	if !m.client.IsAuthenticated() {
		return initCmd
	}

	return tea.Sequence(initCmd, tea.Batch(
		common.NewRefreshListMsg(),
		common.AppStateCmd(common.AppStateTaskList),
	))
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug(fmt.Sprintf("AppModel.Update: %T", msg))

	switch msg := msg.(type) {

	case common.QuitMsg:
		m.quitting = true
		return m, tea.Quit

	case common.ErrorMsg:
		m.error = msg.Err

	case common.AuthMsg:
		err := m.client.SetToken(msg.Token)
		if err != nil {
			return m, common.NewErrorMsg(err)
		}

		return m, tea.Sequence(
			common.NewRefreshListMsg(),
			common.AppStateCmd(common.AppStateTaskList),
		)

	case tea.WindowSizeMsg:
		return m.onWindowSizeMsg(msg)

	case tea.KeyMsg:
		return m.onKeyMsg(msg)

	case common.AppStateMsg:
		m.appState = msg.NewState
	}

	return m.applyUpdates(msg)
}

func (m AppModel) View() string {
	return styles.PageWrapperStyle.Render(m.renderContent())
}

func (m AppModel) renderContent() string {
	if m.quitting {
		return "Bye!\n"
	}

	if m.error != nil {
		return errorStyle.Width(m.width).Render(fmt.Sprintf("Error: %s\n", m.error))
	}

	switch m.appState {
	case common.AppStateLogin:
		return m.loginForm.View()

	case common.AppStateTaskList:
		return m.taskList.View()

	case common.AppStateTaskForm:
		return m.taskForm.View()

	case common.AppStateSettings:
		return m.settings.View()
	}

	return "No state"
}

func (m AppModel) applyUpdates(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	m.loginForm, cmds = utils.GatherUpdates(m.loginForm, msg, cmds)
	m.taskList, cmds = utils.GatherUpdates(m.taskList, msg, cmds)
	m.taskForm, cmds = utils.GatherUpdates(m.taskForm, msg, cmds)
	m.settings, cmds = utils.GatherUpdates(m.settings, msg, cmds)

	return m, utils.BatchIfExists(cmds)
}

func (m AppModel) onWindowSizeMsg(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	docFrameWidth, docFrameHeight := styles.PageWrapperStyle.GetFrameSize()

	m.width = msg.Width

	wrappedMsg := tea.WindowSizeMsg{
		Width:  msg.Width - docFrameWidth,
		Height: msg.Height - docFrameHeight,
	}

	return m.applyUpdates(wrappedMsg)
}

func (m AppModel) onKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	}

	switch m.appState {

	case common.AppStateLogin:
		m.loginForm, cmd = utils.ApplyUpdate(m.loginForm, msg)
		return m, cmd

	case common.AppStateTaskList:
		m.taskList, cmd = utils.ApplyUpdate(m.taskList, msg)
		return m, cmd

	case common.AppStateTaskForm:
		m.taskForm, cmd = utils.ApplyUpdate(m.taskForm, msg)
		return m, cmd

	case common.AppStateSettings:
		m.settings, cmd = utils.ApplyUpdate(m.settings, msg)
		return m, cmd
	}

	return m, nil
}
