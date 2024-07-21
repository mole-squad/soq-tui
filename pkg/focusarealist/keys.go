package focusarealist

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Back   key.Binding
	New    key.Binding
	Edit   key.Binding
	Delete key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new focus area"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit focus area"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete focus area"),
		),
	}
}
