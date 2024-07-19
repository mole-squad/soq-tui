package oldforms

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func NewTextArea(label string, height int) textarea.Model {
	input := textarea.New()
	input.Placeholder = label
	input.ShowLineNumbers = false
	input.Prompt = ""

	input.MaxWidth = 0
	input.FocusedStyle.CursorLine = lipgloss.NewStyle()

	input.SetHeight(height)

	return input
}

func NewTextInput(label string, height int) textinput.Model {
	input := textinput.New()
	input.Placeholder = label
	input.Prompt = ""

	return input
}
