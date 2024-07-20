package common

type AppState int

const (
	AppStateLoading AppState = iota
	AppStateLogin
	AppStateTaskList
	AppStateTaskForm
	AppStateSettings
)
