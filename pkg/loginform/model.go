package loginform

import (
	"context"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mole-squad/soq-tui/pkg/api"
	"github.com/mole-squad/soq-tui/pkg/common"
	"github.com/mole-squad/soq-tui/pkg/forms"
	"github.com/mole-squad/soq-tui/pkg/styles"
)

const (
	usernameInputIdx = iota
	passwordInputIdx
)

type LoginFormModel struct {
	client *api.Client

	keys keyMap
	help help.Model

	focused int

	height int
	width  int

	username textinput.Model
	password textinput.Model
}

func NewLoginFormModel(client *api.Client) LoginFormModel {
	username := forms.NewTextInput("Username", 1)
	password := forms.NewTextInput("Password", 1)

	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '●'

	username.Focus()

	return LoginFormModel{
		client:   client,
		keys:     keys,
		help:     help.New(),
		username: username,
		password: password,
	}
}

func (m LoginFormModel) Init() tea.Cmd {
	return nil
}

func (m LoginFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.onWindowMsg(msg)

	case tea.KeyMsg:
		return m.onKeyMsg(msg)
	}

	return m, nil
}

func (m LoginFormModel) View() string {
	help := m.help.View(m.keys)
	availHeight := m.height - lipgloss.Height(help)

	username := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.InputLabelStyle.Render("Username"),
		styles.InputStyle.Render(m.username.View()),
	)

	password := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.InputLabelStyle.Render("Password"),
		styles.InputStyle.Render(m.password.View()),
	)

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.FormFieldWrapperStyle.Render(username),
		styles.FormFieldWrapperStyle.Render(password),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Height(availHeight).Render(form),
		lipgloss.NewStyle().Width(m.width).Render(help),
	)
}

func (m *LoginFormModel) onKeyMsg(msg tea.KeyMsg) (LoginFormModel, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keys.Quit):
		return *m, common.NewQuitMsg

	case key.Matches(msg, m.keys.Next):
		return m.onNext()

	case key.Matches(msg, m.keys.Submit):
		return *m, m.onSubmit()
	}

	switch m.focused {
	case usernameInputIdx:
		m.username, cmd = m.username.Update(msg)

	case passwordInputIdx:
		m.password, cmd = m.password.Update(msg)

	}

	return *m, cmd
}

func (m *LoginFormModel) onWindowMsg(msg tea.WindowSizeMsg) (LoginFormModel, tea.Cmd) {
	formFieldWrapperWidth, _ := styles.FormFieldWrapperStyle.GetFrameSize()
	inputFrameWidth, _ := styles.InputStyle.GetFrameSize()

	m.height = msg.Height
	m.width = msg.Width

	m.help.Width = msg.Width

	m.username.Width = msg.Width - formFieldWrapperWidth - inputFrameWidth
	m.password.Width = msg.Width - formFieldWrapperWidth - inputFrameWidth

	return *m, nil
}

func (m *LoginFormModel) onNext() (LoginFormModel, tea.Cmd) {
	switch m.focused {
	case usernameInputIdx:
		m.username.Blur()
		m.password.Focus()
		m.focused = passwordInputIdx

	case passwordInputIdx:
		m.password.Blur()
		m.username.Focus()
		m.focused = usernameInputIdx
	}

	return *m, nil
}

func (m *LoginFormModel) onSubmit() tea.Cmd {
	username := m.username.Value()
	password := m.password.Value()

	token, err := m.client.Login(context.Background(), username, password)
	if err != nil {
		return common.NewErrorMsg(err)
	}

	return common.NewAuthMsg(token)
}
