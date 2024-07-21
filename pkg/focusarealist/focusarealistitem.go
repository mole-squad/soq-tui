package focusarealist

import (
	soqapi "github.com/mole-squad/soq-api/api"
)

type FocusAreaListItem struct {
	focusArea soqapi.FocusAreaDTO
}

func (f FocusAreaListItem) Title() string {
	return f.focusArea.Name
}

func (f FocusAreaListItem) Description() string {
	return ""
}

func (f FocusAreaListItem) FilterValue() string {
	return f.focusArea.Name
}
