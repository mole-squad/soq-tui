package settings

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Back       key.Binding
	FocusAreas key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		FocusAreas: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "focus areas"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.FocusAreas, k.Back}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.FocusAreas, k.Back},
	}
}
