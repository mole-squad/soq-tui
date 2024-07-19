package forms

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type FormkeyAction struct {
	Action func() tea.Cmd
	Key    key.Binding
	Label  string
}

type formKeyMap struct {
	Next   key.Binding
	Submit key.Binding
}

// TODO add toggle full help menu
func newFormKeyMap() formKeyMap {
	return formKeyMap{
		Next: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next field"),
		),
		Submit: key.NewBinding(
			// TODO use a different key combo
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit"),
		),
	}
}

func (k formKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Next,
		k.Submit,
	}
}

func (k formKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Submit},
	}
}
