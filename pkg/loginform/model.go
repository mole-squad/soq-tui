package loginform

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mole-squad/soq-tui/pkg/api"
	"github.com/mole-squad/soq-tui/pkg/common"
	"github.com/mole-squad/soq-tui/pkg/forms"
	"github.com/mole-squad/soq-tui/pkg/logger"
	"github.com/mole-squad/soq-tui/pkg/utils"
)

const (
	usernameInputIdx = iota
	passwordInputIdx
)

type LoginFormModel struct {
	client *api.Client
	logger *logger.Logger
	form   forms.Model
}

const (
	loginFormId = "loginForm"
	usernameKey = "username"
	passwordKey = "password"
)

func NewLoginFormModel(logger *logger.Logger, client *api.Client) common.AppView {
	model := LoginFormModel{
		client: client,
		logger: logger,
	}

	username := forms.NewTextInput(
		usernameKey,
		"Username",
	)

	password := forms.NewTextInput(
		passwordKey,
		"Password",
		forms.WithHiddenTextInput(),
	)

	model.form = forms.NewFormModel(
		loginFormId,
		forms.WithField(username),
		forms.WithField(password),
	)

	return model
}

func (m LoginFormModel) Init() tea.Cmd {
	return m.form.Focus()
}

func (m LoginFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case forms.SubmitFormMsg:
		if msg.FormID == loginFormId {
			return m, m.onSubmit()
		}
	}

	m.form, cmd = utils.ApplyUpdate(m.form, msg)

	return m, cmd
}

func (m LoginFormModel) View() string {
	return m.form.View()
}

func (m LoginFormModel) Blur() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m LoginFormModel) Focus() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *LoginFormModel) onSubmit() tea.Cmd {
	values := m.form.Value()

	username := values[usernameKey]
	password := values[passwordKey]

	token, err := m.client.Login(context.Background(), username, password)
	if err != nil {
		return common.NewErrorMsg(err)
	}

	return common.NewAuthMsg(token)
}
