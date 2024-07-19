package utils

import tea "github.com/charmbracelet/bubbletea"

func GatherUpdates[M tea.Model](model M, msg tea.Msg, cmds []tea.Cmd) (M, []tea.Cmd) {
	updatedModel, cmd := ApplyUpdate(model, msg)
	return updatedModel, AppendIfNotNil(cmds, cmd)
}

func ApplyUpdate[M tea.Model](model M, msg tea.Msg) (M, tea.Cmd) {
	updatedModel, cmd := model.Update(msg)

	return updatedModel.(M), cmd
}

func AppendIfNotNil(cmds []tea.Cmd, cmd tea.Cmd) []tea.Cmd {
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return cmds
}

func BatchIfExists(cmds []tea.Cmd) tea.Cmd {
	if len(cmds) > 0 {
		return tea.Batch(cmds...)
	}

	return nil
}
