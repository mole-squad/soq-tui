package app

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mole-squad/soq-tui/pkg/api"
	"github.com/mole-squad/soq-tui/pkg/common"
	"github.com/mole-squad/soq-tui/pkg/loginform"
	"github.com/mole-squad/soq-tui/pkg/taskform"
	"github.com/mole-squad/soq-tui/pkg/tasklist"
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
	quitting  bool
	width     int
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
		loginForm: loginform.NewLoginFormModel(client),
		taskForm:  taskform.NewTaskFormModel(client),
		taskList:  tasklist.NewTaskListModel(client),
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
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

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
		m.width = msg.Width

	case tea.KeyMsg:
		switch m.appState {

		case common.AppStateLogin:
			m.loginForm, cmd = m.loginForm.Update(msg)
			return m, cmd

		case common.AppStateTaskList:
			m.taskList, cmd = m.taskList.Update(msg)
			return m, cmd

		case common.AppStateTaskForm:
			m.taskForm, cmd = m.taskForm.Update(msg)
			return m, cmd
		}

	case common.AppStateMsg:
		m.appState = msg.NewState
	}

	m.loginForm, cmd = m.loginForm.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	m.taskList, cmd = m.taskList.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	m.taskForm, cmd = m.taskForm.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	cmd = nil
	if len(cmds) > 0 {
		cmd = tea.Batch(cmds...)
	}

	return m, cmd
}

func (m AppModel) View() string {
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
	}

	return "No state"
}
