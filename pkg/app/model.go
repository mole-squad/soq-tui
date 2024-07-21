package app

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mole-squad/soq-tui/pkg/api"
	"github.com/mole-squad/soq-tui/pkg/common"
	"github.com/mole-squad/soq-tui/pkg/focusarealist"
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

	views map[common.AppState]common.AppView

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

	views := map[common.AppState]common.AppView{
		common.AppStateLoading:       NewLoadingModel(),
		common.AppStateLogin:         loginform.NewLoginFormModel(client),
		common.AppStateFocusAreaList: focusarealist.New(client),
		common.AppStateTaskList:      tasklist.NewTaskListModel(client),
		common.AppStateTaskForm:      taskform.NewTaskFormModel(client),
		common.AppStateSettings:      settings.NewSettingsModel(client),
	}

	return AppModel{
		appState: common.AppStateLoading,
		client:   client,
		keys:     newKeyMap(),
		views:    views,
	}
}

func (m AppModel) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, view := range m.views {
		cmds = utils.AppendIfNotNil(cmds, view.Init())
	}

	initCmd := utils.BatchIfNotNil(cmds...)

	navCmd := common.AppStateCmd(common.AppStateLogin)
	if m.client.IsAuthenticated() {
		navCmd = common.AppStateCmd(common.AppStateTaskList)
	}

	return tea.Sequence(initCmd, navCmd)
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

		return m, common.AppStateCmd(common.AppStateTaskList)

	case tea.WindowSizeMsg:
		return m.onWindowSizeMsg(msg)

	case tea.KeyMsg:
		return m.onKeyMsg(msg)

	case common.AppStateMsg:
		return m.onAppStateMsg(msg)
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

	return m.views[m.appState].View()
}

func (m AppModel) onAppStateMsg(msg common.AppStateMsg) (tea.Model, tea.Cmd) {
	var (
		blurCmd  tea.Cmd
		focusCmd tea.Cmd
	)

	blurredView, blurCmd := m.views[m.appState].Blur()
	m.views[m.appState] = blurredView.(common.AppView)

	m.appState = msg.NewState

	focusedView, focusCmd := m.views[m.appState].Focus()
	m.views[m.appState] = focusedView.(common.AppView)

	return m, utils.SequenceIfNotNil(blurCmd, focusCmd)
}

func (m AppModel) applyUpdates(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	for key, view := range m.views {
		m.views[key], cmds = utils.GatherUpdates(view, msg, cmds)
	}

	return m, utils.BatchIfNotNil(cmds...)
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
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	}

	updatedView, cmd := m.views[m.appState].Update(msg)
	m.views[m.appState] = updatedView.(common.AppView)

	return m, cmd
}
