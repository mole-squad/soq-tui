package forms

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type selectInputDelegate struct{}

func (d selectInputDelegate) Height() int {
	return 1
}

func (d selectInputDelegate) Spacing() int {
	return 0
}

func (d selectInputDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d selectInputDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(SelectListOption)
	if !ok {
		return
	}

	if index == m.Index() {
		fmt.Fprintf(w, "> %s", i.Label())
	} else {
		fmt.Fprintf(w, "  %s", i.Label())
	}
}
