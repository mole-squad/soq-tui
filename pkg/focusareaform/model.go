package focusareaform

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	soqapi "github.com/mole-squad/soq-api/api"
	"github.com/mole-squad/soq-tui/pkg/api"
	"github.com/mole-squad/soq-tui/pkg/common"
	"github.com/mole-squad/soq-tui/pkg/forms"
	"github.com/mole-squad/soq-tui/pkg/logger"
	"github.com/mole-squad/soq-tui/pkg/utils"
)

const (
	focusAreaFormID = "focusareaform"
	nameFieldID     = "name"
)

type Model struct {
	client *api.Client
	logger *logger.Logger

	isNew     bool
	focusArea soqapi.FocusAreaDTO

	form forms.Model
}

func New(logger *logger.Logger, client *api.Client) common.AppView {
	name := forms.NewTextInput(nameFieldID, "Name")

	form := forms.NewFormModel(
		focusAreaFormID,
		forms.WithField(name),
	)

	return Model{
		client: client,
		logger: logger,
		form:   form,
	}
}

func (m Model) Init() tea.Cmd {
	return m.form.Focus()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case common.CreateFocusAreaMsg:
		return m.onCreate()

	case common.EditFocusAreaMsg:
		return m.onEdit(msg.FocusArea)

	case forms.SubmitFormMsg:
		if msg.FormID == focusAreaFormID {
			return m.onSubmit()
		}
	}

	m.form, cmd = utils.ApplyUpdate(m.form, msg)

	return m, cmd
}

func (m Model) View() string {
	return m.form.View()
}

func (m Model) Blur() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) Focus() (tea.Model, tea.Cmd) {
	return m, m.form.Focus()
}

func (m Model) onCreate() (tea.Model, tea.Cmd) {
	m.isNew = true
	m.focusArea = soqapi.FocusAreaDTO{
		Name: "",
	}

	return m, tea.Batch(
		forms.NewSetFieldValueCmd(focusAreaFormID, nameFieldID, m.focusArea.Name),
	)
}

func (m Model) onEdit(focusArea soqapi.FocusAreaDTO) (tea.Model, tea.Cmd) {
	m.isNew = false
	m.focusArea = focusArea

	return m, tea.Batch(
		forms.NewSetFieldValueCmd(focusAreaFormID, nameFieldID, m.focusArea.Name),
	)
}

func (m Model) onSubmit() (tea.Model, tea.Cmd) {
	values := m.form.Value()

	name := values[nameFieldID]

	var err error
	if m.isNew {
		err = m.createFocusArea(name)
	} else {
		err = m.updateFocusArea(name)
	}

	if err != nil {
		return m, common.NewErrorMsg(err)
	}

	return m, common.AppStateCmd(common.AppStateFocusAreaList)
}

func (m Model) createFocusArea(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultRequestTimeout)
	defer cancel()

	dto := soqapi.CreateFocusAreaRequestDTO{
		Name: name,
	}

	_, err := m.client.CreateFocusArea(ctx, &dto)
	if err != nil {
		return fmt.Errorf("error creating focus area: %w", err)
	}

	return nil
}

func (m Model) updateFocusArea(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultRequestTimeout)
	defer cancel()

	dto := soqapi.UpdateFocusAreaRequestDTO{
		Name: name,
	}

	_, err := m.client.UpdateFocusArea(ctx, m.focusArea.ID, &dto)
	if err != nil {
		return fmt.Errorf("error updating focus area: %w", err)
	}

	return nil
}
