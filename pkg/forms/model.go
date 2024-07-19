package forms

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mole-squad/soq-tui/pkg/common"
	"github.com/mole-squad/soq-tui/pkg/utils"
)

type FormModel struct {
	fields []FormField

	focusedIdx int

	keys formKeyMap
	help help.Model

	height int
	width  int
}

type FormModelOption func(*FormModel)

func NewFormModel(opts ...FormModelOption) FormModel {
	model := FormModel{
		fields: make([]FormField, 0),

		focusedIdx: 0,

		keys: newFormKeyMap(),
		help: help.New(),
	}

	for _, opt := range opts {
		opt(&model)
	}

	return model
}

func (m FormModel) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, field := range m.fields {
		cmds = utils.AppendIfNotNil(cmds, field.Init())
	}

	return utils.BatchIfExists(cmds)
}

func (m FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.onWindowMsg(msg)

	case tea.KeyMsg:
		return m.onKeyMsg(msg)
	}

	return m, nil
}

func (m FormModel) View() string {
	help := m.help.View(m.keys)
	availHeight := m.height - lipgloss.Height(help)

	renderedFields := make([]string, len(m.fields))

	for i, field := range m.fields {
		renderedFields[i] = field.View()
	}

	formContent := lipgloss.JoinVertical(lipgloss.Top, renderedFields...)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Height(availHeight).Render(formContent),
		lipgloss.NewStyle().Width(m.width).Render(help),
	)
}

func (m FormModel) Focus() tea.Cmd {
	if len(m.fields) == 0 {
		return nil
	}

	field := m.fields[m.focusedIdx]

	return field.Focus()
}

func (m FormModel) Blur() tea.Cmd {
	if len(m.fields) == 0 {
		return nil
	}

	field := m.fields[m.focusedIdx]

	return field.Blur()
}

func (m FormModel) Value() map[string]string {
	values := make(map[string]string)

	for _, field := range m.fields {
		values[field.GetID()] = field.GetValue()
	}

	return values
}

func (m FormModel) onWindowMsg(msg tea.WindowSizeMsg) (FormModel, tea.Cmd) {
	m.height = msg.Height
	m.width = msg.Width

	for _, field := range m.fields {
		field.SetWidth(msg.Width)
	}

	return m, nil
}

func (m FormModel) onKeyMsg(msg tea.KeyMsg) (FormModel, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keys.Next):
		return m.next()

	case key.Matches(msg, m.keys.Submit):
		return m, NewSubmitFormMsg
	}

	if len(m.fields) == 0 {
		return m, nil
	}

	if m.focusedIdx >= len(m.fields) {
		return m, common.NewErrorMsg(fmt.Errorf("focused index out of range"))
	}

	field := m.fields[m.focusedIdx]

	m.fields[m.focusedIdx], cmd = utils.ApplyUpdate(field, msg)

	return m, cmd
}

func (m FormModel) next() (FormModel, tea.Cmd) {
	if len(m.fields) == 0 {
		return m, nil
	}

	m.fields[m.focusedIdx].Blur()

	m.focusedIdx = (m.focusedIdx + 1) % len(m.fields)
	m.fields[m.focusedIdx].Focus()

	return m, nil
}

func WithField(field FormField) FormModelOption {
	return func(m *FormModel) {
		m.fields = append(m.fields, field)
	}
}
