package common

type AppState int

const (
	AppStateLoading AppState = iota
	AppStateLogin
	AppStateFocusAreaList
	AppStateFocusAreaForm
	AppStateTaskList
	AppStateTaskForm
	AppStateSettings
)
