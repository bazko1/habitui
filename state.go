package habitui

type TasksState interface {
	Save(tasks TaskList) bool
	Load() TaskList
}
