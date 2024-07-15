package utils

import tea "github.com/charmbracelet/bubbletea"

func GatherUpdates[M tea.Model](model M, msg tea.Msg, cmds []tea.Cmd) (M, []tea.Cmd) {
	updatedModel, cmd := ApplyUpdate(model, msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return updatedModel, cmds
}

func ApplyUpdate[M tea.Model](model M, msg tea.Msg) (M, tea.Cmd) {
	updatedModel, cmd := model.Update(msg)

	return updatedModel.(M), cmd
}
