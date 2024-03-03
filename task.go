package habitui

import (
	"time"
)

// TaskList is a slice of tasks.
type TaskList []Task

// Task is an occurring event that has its own name identifier.
type Task struct {
	name              string
	description       string
	creationDate      time.Time
	completionHistory []time.Time
}

// NewTask constructs Task based on given name and description.
func NewTask(name, description string) Task {
	return Task{
		name,
		description,
		time.Now(),
		make([]time.Time, 0),
	}
}

// GetCompletionHistory returns list of task completion time dates.
func (task Task) GetCompletionHistory() []time.Time {
	return task.completionHistory
}

// MakeTaskCompleted adds current time (time.Now()) to the completionHistory if it wasn't completed yet.
// Each task can be completed once a day.
func (task *Task) MakeTaskCompleted() {
	if l := len(task.completionHistory); l > 0 {
		lastComplete := task.completionHistory[l-1]
		if areSameDates(time.Now(), lastComplete) {
			return
		}
	}

	task.completionHistory = append(task.completionHistory, time.Now())
}

// WasCompletedToday returns whether the Task t was completed at the day pointed
// by time.Now.
func (task Task) WasCompletedToday() bool {
	return task.WasCompletedAtDay(time.Now())
}

// WasCompletedAtDay returns whether the Task tk was completed at the
// day pointed time tme.
func (task Task) WasCompletedAtDay(tme time.Time) bool {
	for _, completion := range task.completionHistory {
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
