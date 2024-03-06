package habitui

import (
	"slices"
	"time"
)

// TaskList is a slice of tasks.
type TaskList []Task

// Task is an occurring event that has its own name identifier.
// Each task can be completed once a day.
type Task struct {
	Name                 string
	Description          string
	CreationDate         time.Time
	YearlyTaskCompletion YearlyTaskCompletion
	GetTime              func() time.Time
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
		make(YearlyTaskCompletion),
		getTime,
	}
}

// MakeTaskCompleted adds current time (getTime()) to the completion history if it wasn't completed yet.
// Each task can be completed once a day.
func (task *Task) MakeTaskCompleted() {
	now := task.GetTime()

	completionsThisYear, exists := task.YearlyTaskCompletion[now.Year()]
	if !exists {
		completionsThisYear = make(MonthlyTaskCompletion)
		task.YearlyTaskCompletion[now.Year()] = completionsThisYear
		completionsThisYear[now.Month()] = []time.Time{now}

		return
	}

	completionsThisMonth, exists := completionsThisYear[now.Month()]
	if !exists {
		completionsThisYear[now.Month()] = []time.Time{now}

		return
	}

	if lastComplete := completionsThisMonth[len(completionsThisMonth)-1]; !AreSameDates(now, lastComplete) {
		completionsThisYear[now.Month()] = append(completionsThisMonth, now)
	}
}

// MonthCompletions returns task completion at given year and month.
func (task Task) MonthCompletions(year int, month time.Month) []time.Time {
	completionsYear, exists := task.YearlyTaskCompletion[year]
	if !exists {
		return nil
	}

	return completionsYear[month]
}

// WasCompletedAt returns whether the Task was completed at the given date.
func (task Task) WasCompletedAt(year int, month time.Month, day int) bool {
	mcmpl := task.MonthCompletions(year, month)
	if mcmpl == nil {
		return false
	}

	atDate := time.Date(year, month, day, 0, 0, 0, 0, &time.Location{})

	return slices.ContainsFunc(mcmpl, func(t time.Time) bool { return AreSameDates(t, atDate) })
}

// WasCompletedToday returns whether the Task was completed.
func (task Task) WasCompletedToday() bool {
	y, m, d := task.GetTime().Date()

	return task.WasCompletedAt(y, m, d)
}

// AreSameDates is a helper function that checks if t1 t2 time.Time
// have the same day date meaning year, month and day.
func AreSameDates(one, other time.Time) bool {
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
