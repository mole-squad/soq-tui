package sidepanelview

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	isPanelOpen bool

	panelWidth int
	height     int
	width      int
}

type SidePanelViewOption func(*Model)

func New(opts ...SidePanelViewOption) Model {
	view := Model{
		isPanelOpen: false,
		panelWidth:  20,
	}

	for _, opt := range opts {
		opt(&view)
	}

	return view
}

func (v Model) Init() tea.Cmd {
	return nil
}

func (v Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return v.onWindowSize(msg.Width, msg.Height)
	}

	return v, nil
}

func (v Model) View() string {
	return "Don't call View - use Render instead"
}

func (v Model) Render(mainPanelContent string, sidePanelContent string) string {
	content := mainPanelContent
	if v.isPanelOpen {
		contentWidth := v.width - v.panelWidth

		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Width(contentWidth).Render(mainPanelContent),
			lipgloss.NewStyle().Width(v.panelWidth).Height(v.height).Render(sidePanelContent),
		)
	}

	return content
}

func (v *Model) SetIsOpen(isOpen bool) {
	v.isPanelOpen = isOpen
}

func (v Model) GetContentWidth() int {
	if v.isPanelOpen {
		return v.width - v.panelWidth
	}

	return v.width
}

func (v Model) onWindowSize(width, height int) (Model, tea.Cmd) {
	v.height = height
	v.width = width

	return v, nil
}
