package tasklist

import (
	soqapi "github.com/mole-squad/soq-api/api"
)

type TaskListItem struct {
	task soqapi.TaskDTO
}

func (t TaskListItem) Title() string {
	return t.task.Summary
}

func (t TaskListItem) Description() string {
	return ""
}

func (t TaskListItem) FilterValue() string {
	return t.task.Summary
}
