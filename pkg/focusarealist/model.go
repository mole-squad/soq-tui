package focusarealist

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mole-squad/soq-tui/pkg/api"
	"github.com/mole-squad/soq-tui/pkg/common"
)

type Model struct {
	client *api.Client

	keys    keyMap
	teaList list.Model
}

func New(client *api.Client) common.AppView {
	listKeys := newKeyMap()

	teaList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	teaList.Title = "Focus Areas"

	teaList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.Back,
			listKeys.New,
			listKeys.Edit,
			listKeys.Delete,
		}
	}

	teaList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.Back,
			listKeys.New,
			listKeys.Edit,
			listKeys.Delete,
		}
	}

	return Model{
		client:  client,
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
	return m.refreshFocusAreas()
}

func (m Model) onKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		return m, common.AppStateCmd(common.AppStateSettings)

	case key.Matches(msg, m.keys.New):
		return m, tea.Sequence(
			common.NewCreateFocusAreaMsg(),
			common.AppStateCmd(common.AppStateFocusAreaForm),
		)

	case key.Matches(msg, m.keys.Edit):
		return m.onEdit()

	case key.Matches(msg, m.keys.Delete):
		return m.onDelete()
	}

	var cmd tea.Cmd
	m.teaList, cmd = m.teaList.Update(msg)

	return m, cmd
}

func (m Model) refreshFocusAreas() (Model, tea.Cmd) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	focusAreas, err := m.client.ListFocusAreas(ctx)
	if err != nil {
		return m, common.NewErrorMsg(fmt.Errorf("error fetching focus areas: %w", err))
	}

	newItems := make([]list.Item, len(focusAreas))
	for i, fa := range focusAreas {
		newItems[i] = FocusAreaListItem{focusArea: fa}
	}

	return m, m.teaList.SetItems(newItems)
}

func (m Model) onEdit() (Model, tea.Cmd) {
	selected := m.teaList.SelectedItem()
	if selected == nil {
		return m, common.NewErrorMsg(fmt.Errorf("no focus area selected"))
	}

	focusAreaItem, ok := selected.(FocusAreaListItem)
	if !ok {
		return m, common.NewErrorMsg(fmt.Errorf("unexpected focus area item type"))
	}

	return m, tea.Sequence(
		common.NewEditFocusAreaMsg(focusAreaItem.focusArea),
		common.AppStateCmd(common.AppStateFocusAreaForm),
	)
}

func (m Model) onDelete() (Model, tea.Cmd) {
	selected := m.teaList.SelectedItem()
	if selected == nil {
		return m, common.NewErrorMsg(fmt.Errorf("no focus area selected"))
	}

	focusAreaItem, ok := selected.(FocusAreaListItem)
	if !ok {
		return m, common.NewErrorMsg(fmt.Errorf("unexpected focus area item type"))
	}

	err := m.client.DeleteFocusArea(context.Background(), focusAreaItem.focusArea.ID)
	if err != nil {
		return m, common.NewErrorMsg(fmt.Errorf("failed to delete focus area: %w", err))
	}

	return m.refreshFocusAreas()
}
