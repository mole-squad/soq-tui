package forms

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mole-squad/soq-tui/pkg/common"
	"github.com/mole-squad/soq-tui/pkg/sidepanelview"
	"github.com/mole-squad/soq-tui/pkg/utils"
)

type Model struct {
	fields []FormField

	focusedIdx int

	panelView sidepanelview.Model

	keys formKeyMap
	help help.Model

	height int
	width  int
}

type FormModelOption func(*Model)

func NewFormModel(opts ...FormModelOption) Model {
	model := Model{
		fields:     make([]FormField, 0),
		focusedIdx: 0,
		keys:       newFormKeyMap(),
		help:       help.New(),
		panelView:  sidepanelview.New(),
	}

	for _, opt := range opts {
		opt(&model)
	}

	return model
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, field := range m.fields {
		cmds = utils.AppendIfNotNil(cmds, field.Init())
	}

	return utils.BatchIfExists(cmds)
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
	help := m.help.View(m.keys)
	availHeight := m.height - lipgloss.Height(help)

	renderedFields := make([]string, len(m.fields))

	for i, field := range m.fields {
		renderedFields[i] = field.View()
	}

	formContent := lipgloss.JoinVertical(lipgloss.Top, renderedFields...)
	content := m.panelView.Render(formContent, "TMP side panel ")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Height(availHeight).Render(content),
		lipgloss.NewStyle().Width(m.width).Render(help),
	)
}

func (m Model) Focus() tea.Cmd {
	if len(m.fields) == 0 {
		return nil
	}

	field := m.fields[m.focusedIdx]

	return field.Focus()
}

func (m Model) Blur() tea.Cmd {
	if len(m.fields) == 0 {
		return nil
	}

	field := m.fields[m.focusedIdx]

	return field.Blur()
}

func (m Model) Value() map[string]string {
	values := make(map[string]string)

	for _, field := range m.fields {
		values[field.GetID()] = field.GetValue()
	}

	return values
}

func (m Model) onWindowMsg(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.height = msg.Height
	m.width = msg.Width

	help := m.help.View(m.keys)
	availHeight := m.height - lipgloss.Height(help)

	var cmd tea.Cmd
	m.panelView, cmd = m.panelView.Update(
		tea.WindowSizeMsg{Width: m.width, Height: availHeight},
	)

	for _, field := range m.fields {
		field.SetWidth(m.panelView.GetContentWidth())
	}

	return m, cmd
}

func (m Model) onKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
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

func (m Model) next() (Model, tea.Cmd) {
	if len(m.fields) == 0 {
		return m, nil
	}

	m.fields[m.focusedIdx].Blur()

	m.focusedIdx = (m.focusedIdx + 1) % len(m.fields)
	m.fields[m.focusedIdx].Focus()

	m.panelView.SetIsOpen(
		m.fields[m.focusedIdx].HasPanelContent(),
	)

	slog.Info("panel width", "width", m.panelView.GetContentWidth())

	for _, field := range m.fields {
		field.SetWidth(m.panelView.GetContentWidth())
	}

	return m, nil
}

func WithField(field FormField) FormModelOption {
	return func(m *Model) {
		m.fields = append(m.fields, field)
	}
}
