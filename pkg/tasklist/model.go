package tasklist

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	soqapi "github.com/mole-squad/soq-api/api"
	"github.com/mole-squad/soq-tui/pkg/api"
	"github.com/mole-squad/soq-tui/pkg/common"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

type taskLoadMsg struct {
	tasks []soqapi.TaskDTO
}

type TaskListModel struct {
	client *api.Client
	tasks  []soqapi.TaskDTO
	keys   keyMap
	list   list.Model
}

func NewTaskListModel(client *api.Client) TaskListModel {
	listKeys := newKeyMap()

	list := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	list.Title = "My Tasks"

	list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.New,
			listKeys.Edit,
			listKeys.Delete,
			listKeys.Resolve,
		}
	}

	list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.New,
			listKeys.Edit,
			listKeys.Delete,
			listKeys.Resolve,
			listKeys.Settings,
		}
	}

	return TaskListModel{
		client: client,
		keys:   listKeys,
		list:   list,
	}
}

func (m TaskListModel) Init() tea.Cmd {
	return nil
}

func (m TaskListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case taskLoadMsg:
		m.tasks = msg.tasks

		newItems := make([]list.Item, len(m.tasks))

		for i, task := range m.tasks {
			newItems[i] = TaskListItem{task: task}
		}

		m.list.SetItems(newItems)

	case common.RefreshListMsg:
		return m, m.getTasks

	case tea.KeyMsg:
		return m.onKeyMsg(msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m TaskListModel) View() string {
	return m.list.View()
}

func (m TaskListModel) getTasks() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tasks, err := m.client.ListTasks(ctx)
	if err != nil {
		return common.ErrorMsg{Err: err}
	}

	return taskLoadMsg{tasks: tasks}
}

func (m TaskListModel) onKeyMsg(msg tea.KeyMsg) (TaskListModel, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.New):
		return m, tea.Sequence(
			common.NewCreateTaskMsg(),
			common.AppStateCmd(common.AppStateTaskForm),
		)

	case key.Matches(msg, m.keys.Edit):
		selected := m.list.SelectedItem()
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

	case key.Matches(msg, m.keys.Delete):
		selected := m.list.SelectedItem()
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

		return m, tea.Sequence(
			common.NewRefreshListMsg(),
		)

	case key.Matches(msg, m.keys.Resolve):
		selected := m.list.SelectedItem()
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

		return m, tea.Sequence(
			common.NewRefreshListMsg(),
		)

	case key.Matches(msg, m.keys.Settings):
		return m, common.AppStateCmd(common.AppStateSettings)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}
