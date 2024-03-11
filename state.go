package habitui

// a FIXME: Probably will not be needed as such interface
// as io.Reader io.Writer will be better interfaces to reuse.
type TasksState interface {
	Save(tasks TaskList) bool
	Load() TaskList
}
