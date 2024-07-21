package tasklist

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	soqapi "github.com/mole-squad/soq-api/api"
	"github.com/mole-squad/soq-tui/pkg/api"
	"github.com/mole-squad/soq-tui/pkg/common"
	"github.com/mole-squad/soq-tui/pkg/logger"
)

type Model struct {
	client  *api.Client
	logger  *logger.Logger
	tasks   []soqapi.TaskDTO
	keys    keyMap
	teaList list.Model
}

func New(logger *logger.Logger, client *api.Client) common.AppView {
	listKeys := newKeyMap()

	teaList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	teaList.Title = "Tasks"

	teaList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.New,
			listKeys.Edit,
			listKeys.Delete,
			listKeys.Resolve,
		}
	}

	teaList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.New,
			listKeys.Edit,
			listKeys.Delete,
			listKeys.Resolve,
			listKeys.Settings,
		}
	}

	return Model{
		client:  client,
		logger:  logger,
		keys:    listKeys,
		teaList: teaList,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.teaList.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		return m.onKeyMsg(msg)
	}

	var cmd tea.Cmd
	m.teaList, cmd = m.teaList.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return m.teaList.View()
}

func (m Model) Blur() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) Focus() (tea.Model, tea.Cmd) {
	return m.refreshTasks()
}

func (m Model) refreshTasks() (Model, tea.Cmd) {
	tasks, err := m.getTasks()
	if err != nil {
		return m, common.NewErrorMsg(err)
	}

	m.tasks = tasks

	newItems := make([]list.Item, len(m.tasks))

	for i, task := range m.tasks {
		newItems[i] = TaskListItem{task: task}
	}

	return m, m.teaList.SetItems(newItems)
}

func (m Model) getTasks() ([]soqapi.TaskDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tasks, err := m.client.ListTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load task list: %w", err)
	}

	return tasks, nil
}

func (m Model) onKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.New):
		return m, tea.Sequence(
			common.NewCreateTaskMsg(),
			common.AppStateCmd(common.AppStateTaskForm),
		)

	case key.Matches(msg, m.keys.Edit):
		return m.onEditTask()

	case key.Matches(msg, m.keys.Delete):
		return m.onDeleteTask()

	case key.Matches(msg, m.keys.Resolve):
		return m.onResolveTask()

	case key.Matches(msg, m.keys.Settings):
		return m, common.AppStateCmd(common.AppStateSettings)
	}

	var cmd tea.Cmd
	m.teaList, cmd = m.teaList.Update(msg)

	return m, cmd
}

func (m Model) onEditTask() (Model, tea.Cmd) {
	selected := m.teaList.SelectedItem()
	if selected == nil {
		return m, common.NewErrorMsg(fmt.Errorf("no task selected"))
	}

	taskItem, ok := selected.(TaskListItem)
	if !ok {
		return m, common.NewErrorMsg(fmt.Errorf("unexpected task item type"))
	}

	return m, tea.Sequence(
		common.NewSelectTaskMsg(taskItem.task),
		common.AppStateCmd(common.AppStateTaskForm),
	)
}

func (m Model) onDeleteTask() (Model, tea.Cmd) {
	selected := m.teaList.SelectedItem()
	if selected == nil {
		return m, common.NewErrorMsg(fmt.Errorf("no task selected"))
	}

	taskItem, ok := selected.(TaskListItem)
	if !ok {
		return m, common.NewErrorMsg(fmt.Errorf("unexpected task item type"))
	}

	err := m.client.DeleteTask(context.Background(), taskItem.task.ID)
	if err != nil {
		return m, common.NewErrorMsg(fmt.Errorf("failed to delete task: %w", err))
	}

	return m.refreshTasks()
}

func (m Model) onResolveTask() (Model, tea.Cmd) {
	selected := m.teaList.SelectedItem()
	if selected == nil {
		return m, common.NewErrorMsg(fmt.Errorf("no task selected"))
	}

	taskItem, ok := selected.(TaskListItem)
	if !ok {
		return m, common.NewErrorMsg(fmt.Errorf("unexpected task item type"))
	}

	_, err := m.client.ResolveTask(context.Background(), taskItem.task.ID)
	if err != nil {
		return m, common.NewErrorMsg(fmt.Errorf("failed to resolve task: %w", err))
	}

	return m.refreshTasks()
}
