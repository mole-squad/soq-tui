package taskform

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	soqapi "github.com/mole-squad/soq-api/api"
	"github.com/mole-squad/soq-tui/pkg/api"
	"github.com/mole-squad/soq-tui/pkg/common"
	"github.com/mole-squad/soq-tui/pkg/forms"
	"github.com/mole-squad/soq-tui/pkg/utils"
)

const (
	taskFormID       = "taskform"
	summaryFieldID   = "summary"
	notesFieldID     = "notes"
	focusAreaFieldID = "focusarea"
)

type TaskFormModel struct {
	client *api.Client

	isNewTask bool
	task      soqapi.TaskDTO

	focusareas []soqapi.FocusAreaDTO

	form forms.Model
}

func NewTaskFormModel(client *api.Client) TaskFormModel {
	summary := forms.NewTextInput(summaryFieldID, "Summary")
	notes := forms.NewTextInput(notesFieldID, "Notes")
	focusArea := forms.NewSelectInput(focusAreaFieldID, "Focus Area")

	form := forms.NewFormModel(
		taskFormID,
		forms.WithField(summary),
		forms.WithField(notes),
		forms.WithField(focusArea),
	)

	return TaskFormModel{
		client: client,
		form:   form,
	}
}

func (m TaskFormModel) Init() tea.Cmd {
	return m.form.Focus()
}

func (m TaskFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case common.CreateTaskMsg:
		return m, tea.Sequence(
			m.onTaskCreate(),
			m.form.Focus(),
		)

	case common.SelectTaskMsg:
		return m, tea.Sequence(
			m.onTaskSelect(msg.Task),
			m.form.Focus(),
		)

	case forms.SubmitFormMsg:
		if msg.FormID == taskFormID {
			return m, tea.Sequence(
				m.submitTask(),
				m.form.Blur(),
			)
		}
	}

	m.form, cmd = utils.ApplyUpdate(m.form, msg)

	return m, cmd
}

func (m TaskFormModel) View() string {
	return m.form.View()
}

func (m *TaskFormModel) refreshFocusAreas() tea.Cmd {
	slog.Debug("Refreshing focus areas")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	focusAreas, err := m.client.ListFocusAreas(ctx)
	if err != nil {
		return common.NewErrorMsg(fmt.Errorf("error fetching focus areas: %w", err))
	}

	slog.Debug("Focus areas fetched", "count", len(focusAreas))

	var opts = make([]forms.SelectOption, len(focusAreas))
	for i, fa := range focusAreas {
		opts[i] = NewFocusAreaOption(fa)
	}

	m.focusareas = focusAreas

	return forms.NewSetSelectOptionsCmd(focusAreaFieldID, opts)
}

func (m *TaskFormModel) onTaskCreate() tea.Cmd {
	refreshCmd := m.refreshFocusAreas()

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

	return tea.Sequence(
		refreshCmd,
		tea.Batch(
			forms.NewSetFieldValueCmd(summaryFieldID, m.task.Summary),
			forms.NewSetFieldValueCmd(notesFieldID, m.task.Notes),
			forms.NewSetFieldValueCmd(focusAreaFieldID, strconv.FormatUint(uint64(focusArea.ID), 10)),
		),
	)
}

func (m *TaskFormModel) onTaskSelect(task soqapi.TaskDTO) tea.Cmd {
	refreshCmd := m.refreshFocusAreas()

	slog.Debug("Editing task", "task", task)

	if len(m.focusareas) == 0 {
		return common.NewErrorMsg(fmt.Errorf("no focus areas available"))
	}

	m.isNewTask = false
	m.task = task

	return tea.Sequence(
		refreshCmd,
		tea.Batch(
			forms.NewSetFieldValueCmd(summaryFieldID, m.task.Summary),
			forms.NewSetFieldValueCmd(notesFieldID, m.task.Notes),
			forms.NewSetFieldValueCmd(focusAreaFieldID, strconv.FormatUint(uint64(task.FocusArea.ID), 10)),
		),
	)
}

func (m *TaskFormModel) submitTask() tea.Cmd {
	values := m.form.Value()

	summary := values[summaryFieldID]
	notes := values[notesFieldID]
	focusAreaID, err := strconv.ParseUint(values[focusAreaFieldID], 10, 64)
	if err != nil {
		return common.NewErrorMsg(fmt.Errorf("error parsing focus area ID: %w", err))
	}

	// TODO validation

	if m.isNewTask {
		err = m.createTask(summary, notes, uint(focusAreaID))
	} else {
		err = m.updateTask(summary, notes, uint(focusAreaID))
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
