package forms

import (
	teatextinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mole-squad/soq-tui/pkg/styles"
)

type TextInput struct {
	id    string
	label string

	teaInput teatextinput.Model

	width int
}

type TextInputOption func(*TextInput)

func NewTextInput(id string, label string, opts ...TextInputOption) FormField {
	teaInput := teatextinput.New()
	teaInput.Prompt = ""

	input := &TextInput{
		id:       id,
		label:    label,
		teaInput: teaInput,
	}

	for _, opt := range opts {
		opt(input)
	}

	return input
}

func (t *TextInput) Init() tea.Cmd {
	return nil
}

func (t *TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	t.teaInput, cmd = t.teaInput.Update(msg)

	return t, cmd
}

func (t *TextInput) View() string {
	renderedLabel := ""
	if t.label != "" {
		renderedLabel = styles.InputLabelStyle.Render(t.label)
	}

	frameWidth, _ := styles.InputStyle.GetFrameSize()
	renderedInput := styles.InputStyle.
		Width(t.width - frameWidth).
		Render(t.teaInput.View())

	return lipgloss.JoinVertical(
		lipgloss.Left,
		renderedLabel,
		renderedInput,
	)
}

func (t *TextInput) Blur() tea.Cmd {
	t.teaInput.Blur()

	return nil
}

func (t *TextInput) Focus() tea.Cmd {
	return t.teaInput.Focus()
}

func (t *TextInput) HasPanelContent() bool {
	return false
}

func (t *TextInput) GetID() string {
	return t.id
}

func (t *TextInput) GetValue() string {
	return t.teaInput.Value()
}

func (t *TextInput) SetValue(value string) {
	t.teaInput.SetValue(value)
}

func (t *TextInput) SetWidth(width int) {
	t.width = width

	inputFrameWidth, _ := styles.InputStyle.GetFrameSize()

	t.teaInput.Width = width - inputFrameWidth
}

func (t *TextInput) SetHeight(int) {}

func WithHiddenTextInput() TextInputOption {
	return func(t *TextInput) {
		t.teaInput.EchoMode = teatextinput.EchoPassword
		t.teaInput.EchoCharacter = '•'
	}
}
