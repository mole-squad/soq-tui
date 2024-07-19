package forms

import (
	"log/slog"

	teatextinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mole-squad/soq-tui/pkg/styles"
)

type TextInput struct {
	id    string
	label string

	teaInput teatextinput.Model
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

	renderedInput := styles.InputStyle.Render(t.teaInput.View())

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

func (t *TextInput) GetID() string {
	return t.id
}

func (t *TextInput) GetValue() string {
	return t.teaInput.Value()
}

func (t *TextInput) SetValue(string) {
}

func (t *TextInput) SetWidth(width int) {
	inputFrameWidth, _ := styles.InputStyle.GetFrameSize()

	slog.Info("Setting width of text input", "width", width, "inputFrameWidth", inputFrameWidth)

	t.teaInput.Width = width - inputFrameWidth
}

func (t *TextInput) SetHeight(int) {
}

func WithHiddenTextInput() TextInputOption {
	return func(t *TextInput) {
		t.teaInput.EchoMode = teatextinput.EchoPassword
		t.teaInput.EchoCharacter = '•'
	}
}
