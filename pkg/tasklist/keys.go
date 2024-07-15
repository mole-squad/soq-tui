package tasklist

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	New      key.Binding
	Edit     key.Binding
	Delete   key.Binding
	Resolve  key.Binding
	Settings key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new task"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit task"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete task"),
		),
		Resolve: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "resolve task"),
		),
		Settings: key.NewBinding(
			key.WithKeys(","),
			key.WithHelp(",", "settings"),
		),
	}
}
