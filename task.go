package habitui

import (
	"time"
)

// Task is an occuring event that has its own name identifier.
type Task struct {
	name              string
	description       string
	creationDate      time.Time
	completionHistory []time.Time
}

// NewTask constructs Task based on given name and description.
// The Task
func NewTask(name, description string) (t Task) {
	t.name = name
	t.description = description
	t.creationDate = time.Now()
	t.completionHistory = make([]time.Time, 0)

	return t
}

func (t *Task) MakeTaskCompleted() {
	t.completionHistory = append(t.completionHistory, time.Now())
}

// WasCompletedToday returns whether the Task t was completed at the day pointed
// by time.Now.
func (t Task) WasCompletedToday() bool {
	return t.WasCompletedAtDay(time.Now())
}

// WasCompletedAtDay returns whether the Task tk was completed at the
// day pointed time tme.
func (tk Task) WasCompletedAtDay(tme time.Time) bool {
	type date struct {
		y int
		m time.Month
		d int
	}
	timeDate := date{}
	timeDate.y, timeDate.m, timeDate.d = tme.Date()
	for _, completion := range tk.completionHistory {
		completionDate := date{}
		completionDate.y, completionDate.m, completionDate.d = completion.Date()
		if timeDate == completionDate {
			return true
		}
	}

	return false
}
