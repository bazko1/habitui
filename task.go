package habitui

import (
	"time"
)

// TaskList is a slice of tasks.
type TaskList []Task

// Task is an occurring event that has its own name identifier.
type Task struct {
	Name              string
	Description       string
	CreationDate      time.Time
	CompletionHistory []time.Time
	GetTime           func() time.Time
}

func NewTask(name, description string) Task {
	return NewTaskWithCustomTime(name, description, time.Now)
}

// NewTaskWithCustomTime constructs Task based on given name, description and time returning function.
func NewTaskWithCustomTime(name, description string, getTime func() time.Time) Task {
	if getTime == nil {
		getTime = time.Now
	}

	return Task{
		name,
		description,
		getTime(),
		make([]time.Time, 0),
		getTime,
	}
}

// MakeTaskCompleted adds current time (getTime()) to the CompletionHistory if it wasn't completed yet.
// Each task can be completed once a day.
func (task *Task) MakeTaskCompleted() {
	if l := len(task.CompletionHistory); l > 0 {
		lastComplete := task.CompletionHistory[l-1]
		if areSameDates(task.GetTime(), lastComplete) {
			return
		}
	}

	task.CompletionHistory = append(task.CompletionHistory, task.GetTime())
}

// WasCompletedToday returns whether the Task t was completed at the day pointed
// by time.Now.
func (task Task) WasCompletedToday() bool {
	return task.WasCompletedAtDay(task.GetTime())
}

// WasCompletedAtDay returns whether the Task tk was completed at the
// day pointed time tme.
func (task Task) WasCompletedAtDay(tme time.Time) bool {
	for _, completion := range task.CompletionHistory {
		if areSameDates(completion, tme) {
			return true
		}
	}

	return false
}

// areSameDates is a helper function that checks if t1 t2 time.Time
// have the same day date meaning year, month and day.
func areSameDates(one, other time.Time) bool {
	type date struct {
		y int
		m time.Month
		d int
	}

	y, m, d := one.Date()
	oneDate := date{y, m, d}
	y, m, d = other.Date()
	otherDate := date{y, m, d}

	return oneDate == otherDate
}
