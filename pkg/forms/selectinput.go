package forms

import (
	tealist "github.com/charmbracelet/bubbles/list"
	teatextinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mole-squad/soq-tui/pkg/styles"
	"github.com/mole-squad/soq-tui/pkg/utils"
)

type SelectOption interface {
	Label() string
	Value() string
}

type SelectInput struct {
	id    string
	label string

	inputModel teatextinput.Model
	listModel  tealist.Model

	width int
}

type SelectInputOption func(*SelectInput)

func NewSelectInput(id string, label string, opts ...SelectInputOption) FormField {
	inputModel := teatextinput.New()
	inputModel.Prompt = ""

	listModel := tealist.New([]tealist.Item{}, selectInputDelegate{}, 0, 0)
	listModel.Title = label

	listModel.SetShowHelp(false)
	listModel.SetShowPagination(false)
	listModel.SetShowFilter(false)
	listModel.SetShowStatusBar(false)

	input := &SelectInput{
		id:         id,
		label:      label,
		inputModel: inputModel,
		listModel:  listModel,
	}

	for _, opt := range opts {
		opt(input)
	}

	return input
}

func (s SelectInput) Init() tea.Cmd {
	return nil
}

func (s *SelectInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		listCmd  tea.Cmd
		inputCmd tea.Cmd
	)

	switch msg := msg.(type) {

	case SetSelectOptionsMsg:
		if msg.InputID == s.id {
			s.SetOptions(msg.Options)
		}
	}

	selected := s.getSelectedItem()
	if selected != nil {
		s.inputModel.SetValue(selected.Label())
	}

	s.listModel, listCmd = s.listModel.Update(msg)
	s.inputModel, inputCmd = s.inputModel.Update(msg)

	return s, utils.BatchIfExists(listCmd, inputCmd)
}

func (s SelectInput) View() string {
	renderedLabel := ""
	if s.label != "" {
		renderedLabel = styles.InputLabelStyle.Render(s.label)
	}

	frameWidth, _ := styles.InputStyle.GetFrameSize()
	renderedInput := styles.InputStyle.
		Width(s.width - frameWidth).
		Render(s.inputModel.View())

	return lipgloss.JoinVertical(
		lipgloss.Left,
		renderedLabel,
		renderedInput,
	)
}

func (s SelectInput) ViewSidePanel() string {
	return s.listModel.View()
}

func (s SelectInput) Blur() tea.Cmd {
	return nil
}

func (s SelectInput) Focus() tea.Cmd {
	// Dont call focus on the input model because it will show the cursor
	return nil
}

func (s SelectInput) HasPanelContent() bool {
	return true
}

func (s SelectInput) GetID() string {
	return s.id
}

func (s SelectInput) GetValue() string {
	return s.getSelectedItem().Value()
}

func (s *SelectInput) SetValue(selected string) {
	for i, opt := range s.listModel.Items() {
		option, ok := opt.(SelectListOption)
		if !ok {
			return
		}

		if option.opt.Value() == selected {
			s.inputModel.SetValue(option.opt.Label())
			s.listModel.Select(i)
		}
	}
}

func (s *SelectInput) SetOptions(opts []SelectOption) {
	newItems := make([]tealist.Item, len(opts))

	for i, opt := range opts {
		newItems[i] = SelectListOption{opt: opt}
	}

	s.listModel.SetItems(newItems)
}

func (s *SelectInput) SetSize(width int, height int) {
	s.width = width

	inputFrameWidth, _ := styles.InputStyle.GetFrameSize()

	s.inputModel.Width = width - inputFrameWidth

	s.listModel.SetSize(width, height)
}

func (s *SelectInput) SetPanelSize(width int, height int) {
	s.listModel.SetSize(width, height)
}

func (m SelectInput) getSelectedItem() SelectOption {
	opt := m.listModel.SelectedItem()
	option, ok := opt.(SelectListOption)
	if !ok {
		return nil
	}

	return option.opt
}
