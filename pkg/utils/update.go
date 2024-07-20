package utils

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Updater interface {
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
}

func GatherUpdates[M Updater](model M, msg tea.Msg, cmds []tea.Cmd) (M, []tea.Cmd) {
	updatedModel, cmd := ApplyUpdate(model, msg)
	return updatedModel, AppendIfNotNil(cmds, cmd)
}

func ApplyUpdate[M Updater](model M, msg tea.Msg) (M, tea.Cmd) {
	updatedModel, cmd := model.Update(msg)

	return updatedModel.(M), cmd
}

func AppendIfNotNil(cmds []tea.Cmd, cmd tea.Cmd) []tea.Cmd {
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return cmds
}

func BatchIfExists(cmds ...tea.Cmd) tea.Cmd {
	if len(cmds) > 0 {
		return tea.Batch(cmds...)
	}

	return nil
}
