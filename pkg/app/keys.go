package app

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}
}
