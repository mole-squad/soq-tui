package taskform

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	soqapi "github.com/mole-squad/soq-api/api"
	"github.com/mole-squad/soq-tui/pkg/api"
	"github.com/mole-squad/soq-tui/pkg/common"
	"github.com/mole-squad/soq-tui/pkg/oldforms"
	"github.com/mole-squad/soq-tui/pkg/selectinput"
	"github.com/mole-squad/soq-tui/pkg/styles"
)

const (
	summaryInputIdx = iota
	notesInputIdx
	focusAreaInputIdx
)

const (
	sidePanelWidth = 30
)

type toggleSidePanelMsg struct {
	isOpen bool
}

func NewToggleSidePanelMsg(isOpen bool) tea.Cmd {
	return func() tea.Msg {
		return toggleSidePanelMsg{isOpen: isOpen}
	}
}

type TaskFormModel struct {
	client *api.Client

	isNewTask bool
	task      soqapi.TaskDTO

	focusareas []soqapi.FocusAreaDTO

	height int
	width  int

	isSidePanelVisible bool

	focused        int
	summary        textarea.Model
	notes          textarea.Model
	focusAreaInput selectinput.SelectInputModel

	keys keyMap
	help help.Model
}

func NewTaskFormModel(client *api.Client) TaskFormModel {
	summaryInput := oldforms.NewTextArea("Summary", 2)
	notesInput := oldforms.NewTextArea("Notes", 5)

	return TaskFormModel{
		client:         client,
		keys:           newKeyMap(),
		help:           help.New(),
		summary:        summaryInput,
		height:         0,
		notes:          notesInput,
		focusAreaInput: *selectinput.NewSelectInputModel("Focus Area"),
		focused:        summaryInputIdx,
		width:          0,
	}
}

func (m TaskFormModel) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
	)
}

func (m TaskFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		return m.onKeyMsg(msg)

	case tea.WindowSizeMsg:
		m.onWindowSize(msg.Width, msg.Height)

	case common.CreateTaskMsg:
		cmd = m.onTaskCreate()
		return m, cmd

	case common.SelectTaskMsg:
		cmd = m.onTaskSelect(msg.Task)
		return m, cmd

	case toggleSidePanelMsg:
		m.onSidePanelToggle(msg.isOpen)
	}

	m.summary, cmd = m.summary.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	m.notes, cmd = m.notes.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	m.focusAreaInput, cmd = m.focusAreaInput.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	cmd = nil
	if len(cmds) > 0 {
		cmd = tea.Batch(cmds...)
	}

	return m, cmd
}

func (m TaskFormModel) View() string {
	sectionFrameWidth, sectionFrameHeight := styles.BorderStyle.GetFrameSize()

	help := m.help.View(m.keys)
	availHeight := m.height - lipgloss.Height(help)

	summary := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.InputLabelStyle.Render("Summary"),
		styles.InputStyle.Render(m.summary.View()),
	)

	notes := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.InputLabelStyle.Render("Notes"),
		styles.InputStyle.Render(m.notes.View()),
	)

	focusArea := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.InputLabelStyle.Render("Focus Area"),
		styles.InputStyle.Render(m.focusAreaInput.View()),
	)

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.FormFieldWrapperStyle.Render(summary),
		styles.FormFieldWrapperStyle.Render(notes),
		styles.FormFieldWrapperStyle.Render(focusArea),
	)

	content := form
	if m.isSidePanelVisible {
		panelContent := m.focusAreaInput.ViewSelectPanel()
		formWidth := m.width - sidePanelWidth

		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Width(formWidth).Render(form),
			styles.BorderStyle.Width(sidePanelWidth-sectionFrameWidth).Height(availHeight-sectionFrameHeight).Render(panelContent),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Height(availHeight).Render(content),
		lipgloss.NewStyle().Width(m.width-sectionFrameWidth).Render(help),
	)
}

func (m *TaskFormModel) onWindowSize(width, height int) {
	m.height = height
	m.width = width

	m.setInputSizes()
}

func (m *TaskFormModel) onSidePanelToggle(isOpen bool) {
	m.isSidePanelVisible = isOpen
	m.setInputSizes()
}

func (m *TaskFormModel) setInputSizes() {
	sectionFrameWidth, sectionFrameHeight := styles.BorderStyle.GetFrameSize()
	formFieldWrapperWidth, _ := styles.FormFieldWrapperStyle.GetFrameSize()
	inputFrameWidth, _ := styles.InputStyle.GetFrameSize()

	help := m.help.View(m.keys)
	helpHeight := lipgloss.Height(help)

	availWidth := m.width
	if m.isSidePanelVisible {
		availWidth -= sidePanelWidth
	}

	inputWidth := availWidth - formFieldWrapperWidth - inputFrameWidth
	m.summary.SetWidth(inputWidth)
	m.notes.SetWidth(inputWidth)

	m.help.Width = m.width - sectionFrameWidth

	m.focusAreaInput.SetSize(sidePanelWidth-sectionFrameWidth, m.height-sectionFrameHeight-helpHeight)
}

func (m TaskFormModel) onKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {

	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = !m.help.ShowAll
		return m, nil

	case key.Matches(msg, m.keys.Exit):
		return m, common.AppStateCmd(common.AppStateTaskList)

	case key.Matches(msg, m.keys.Save):
		return m, m.submitTask()

	case key.Matches(msg, m.keys.Next):
		return m, m.onNextField()
	}

	switch m.focused {
	case summaryInputIdx:
		m.summary, cmd = m.summary.Update(msg)
		return m, cmd

	case notesInputIdx:
		m.notes, cmd = m.notes.Update(msg)
		return m, cmd

	case focusAreaInputIdx:
		m.focusAreaInput, cmd = m.focusAreaInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *TaskFormModel) onFocusAreaRefresh() error {
	slog.Debug("Refreshing focus areas")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	focusAreas, err := m.client.ListFocusAreas(ctx)
	if err != nil {
		return fmt.Errorf("error fetching focus areas: %w", err)
	}

	slog.Debug("Focus areas fetched", "count", len(focusAreas))

	var opts = make([]selectinput.SelectOption, len(focusAreas))
	for i, fa := range focusAreas {
		opts[i] = NewFocusAreaOption(fa)
	}

	m.focusAreaInput.SetOptions(opts)
	m.focusareas = focusAreas

	return nil
}

func (m *TaskFormModel) onTaskCreate() tea.Cmd {
	err := m.onFocusAreaRefresh()
	if err != nil {
		return common.NewErrorMsg(fmt.Errorf("failed to refresh focus areas: %w", err))
	}

	slog.Debug("Creating new task")

	if len(m.focusareas) == 0 {
		return common.NewErrorMsg(fmt.Errorf("no focus areas available"))
	}

	focusArea := m.focusareas[0]

	m.isNewTask = true
	m.task = soqapi.TaskDTO{
		Summary:   "",
		Notes:     "",
		FocusArea: focusArea,
	}

	m.setFormStateFromModel()

	return nil
}

func (m *TaskFormModel) onTaskSelect(task soqapi.TaskDTO) tea.Cmd {
	err := m.onFocusAreaRefresh()
	if err != nil {
		return common.NewErrorMsg(fmt.Errorf("failed to refresh focus areas: %w", err))
	}

	slog.Debug("Editing task", "task", task)

	if len(m.focusareas) == 0 {
		slog.Error("No focus areas available")
		return nil
	}

	m.isNewTask = false
	m.task = task
	m.focused = summaryInputIdx
	m.isSidePanelVisible = false

	m.setFormStateFromModel()

	return nil
}

func (m *TaskFormModel) onNextField() tea.Cmd {
	if m.focused == summaryInputIdx {
		m.summary.Blur()
		m.notes.Focus()
		m.focused = notesInputIdx
	} else if m.focused == notesInputIdx {
		m.notes.Blur()
		m.focusAreaInput.Focus()

		m.focused = focusAreaInputIdx

		return NewToggleSidePanelMsg(true)
	} else if m.focused == focusAreaInputIdx {
		m.focusAreaInput.Blur()
		m.summary.Focus()

		m.focused = summaryInputIdx

		selectedFocusArea := m.focusAreaInput.SelectedItem().(*focusAreaOption)
		if selectedFocusArea != nil {
			m.task.FocusArea = selectedFocusArea.focusArea
		}

		return NewToggleSidePanelMsg(false)
	}

	return nil
}

func (m *TaskFormModel) setFormStateFromModel() {
	m.summary.SetValue(m.task.Summary)
	m.notes.SetValue(m.task.Notes)
	m.focusAreaInput.SetSelected(m.task.FocusArea.ID)

	m.focusAreaInput.Blur()
	m.notes.Blur()
	m.summary.Focus()
}

func (m *TaskFormModel) submitTask() tea.Cmd {
	summary := m.summary.Value()
	notes := m.notes.Value()
	focusAreaID := m.focusAreaInput.Value().(uint)

	// TODO validation

	var err error
	if m.isNewTask {
		err = m.createTask(summary, notes, focusAreaID)
	} else {
		err = m.updateTask(summary, notes, focusAreaID)
	}

	if err != nil {
		slog.Error("Error submitting task", "error", err)
		return common.NewErrorMsg(err)
	}

	return tea.Sequence(
		common.NewRefreshListMsg(),
		common.AppStateCmd(common.AppStateTaskList),
	)
}

func (m *TaskFormModel) createTask(summary, notes string, focusAreaID uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dto := soqapi.CreateTaskRequestDto{
		Summary:     summary,
		Notes:       notes,
		FocusAreaID: focusAreaID,
	}

	_, err := m.client.CreateTask(ctx, &dto)
	if err != nil {
		return fmt.Errorf("error creating task: %w", err)
	}

	return nil
}

func (m *TaskFormModel) updateTask(summary, notes string, focusAreaID uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dto := soqapi.UpdateTaskRequestDto{
		Summary:     summary,
		Notes:       notes,
		FocusAreaID: focusAreaID,
	}

	_, err := m.client.UpdateTask(ctx, m.task.ID, &dto)
	if err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	return nil
}
