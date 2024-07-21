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
)

type taskLoadMsg struct {
	tasks []soqapi.TaskDTO
}

type TaskListModel struct {
	client  *api.Client
	tasks   []soqapi.TaskDTO
	keys    keyMap
	teaList list.Model
}

func NewTaskListModel(client *api.Client) common.AppView {
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

	return TaskListModel{
		client:  client,
		keys:    listKeys,
		teaList: teaList,
	}
}

func (m TaskListModel) Init() tea.Cmd {
	return nil
}

func (m TaskListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m TaskListModel) View() string {
	return m.teaList.View()
}

func (m TaskListModel) Blur() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m TaskListModel) Focus() (tea.Model, tea.Cmd) {
	return m.refreshTasks()
}

func (m TaskListModel) refreshTasks() (TaskListModel, tea.Cmd) {
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

func (m TaskListModel) getTasks() ([]soqapi.TaskDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tasks, err := m.client.ListTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load task list: %w", err)
	}

	return tasks, nil
}

func (m TaskListModel) onKeyMsg(msg tea.KeyMsg) (TaskListModel, tea.Cmd) {
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

func (m TaskListModel) onEditTask() (TaskListModel, tea.Cmd) {
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

func (m TaskListModel) onDeleteTask() (TaskListModel, tea.Cmd) {
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

func (m TaskListModel) onResolveTask() (TaskListModel, tea.Cmd) {
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
