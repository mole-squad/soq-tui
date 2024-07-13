package taskform

import (
	soqapi "github.com/mole-squad/soq-api/api"
)

type focusAreaOption struct {
	focusArea soqapi.FocusAreaDTO
}

func NewFocusAreaOption(fa soqapi.FocusAreaDTO) *focusAreaOption {
	return &focusAreaOption{
		focusArea: fa,
	}
}

func (f *focusAreaOption) Label() string {
	return f.focusArea.Name
}

func (f *focusAreaOption) Value() interface{} {
	return f.focusArea.ID
}
