package taskform

import (
	"strconv"

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

func (f *focusAreaOption) Value() string {
	return strconv.FormatUint(uint64(f.focusArea.ID), 10)
}
