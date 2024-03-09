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
	Name         string
	Description  string
	CreationDate time.Time
	GetTime      func() time.Time

	YearlyTaskCompletion YearlyTaskCompletion
	LastTimeCompleted    time.Time
	currentStrike        uint
	YearlyBestStrike     YearlyBestStrike
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
		getTime,
		make(YearlyTaskCompletion),
		time.Time{},
		0,
		make(YearlyBestStrike),
	}
}

// MakeTaskCompleted adds current time (getTime()) to the completion history if it wasn't completed yet.
// Each task can be completed once a day.
func (task *Task) MakeTaskCompleted() {
	now := task.GetTime()

	if AreSameDates(now, task.LastTimeCompleted) {
		return
	}

	if task.LastTimeCompleted.IsZero() || !task.IsStrikeContinued() {
		task.currentStrike = 1
	}

	if AreSameDates(now, task.LastTimeCompleted.AddDate(0, 0, 1)) {
		task.currentStrike++
	}

	task.LastTimeCompleted = now

	if checkHistoricRecodExistOrCreate(
		task.YearlyTaskCompletion,
		now.Year(),
		now.Month(), []time.Time{now}) {
		return
	}

	completionsThisYear := task.YearlyTaskCompletion[now.Year()]
	completionsThisMonth := completionsThisYear[now.Month()]

	if lastComplete := completionsThisMonth[len(completionsThisMonth)-1]; !AreSameDates(now, lastComplete) {
		completionsThisYear[now.Month()] = append(completionsThisMonth, now)
	}
}

// MonthCompletionTime returns task completion at given year and month.
func (task Task) MonthCompletionTime(year int, month time.Month) []time.Time {
	completionsYear, exists := task.YearlyTaskCompletion[year]
	if !exists {
		return nil
	}

	return completionsYear[month]
}

// WasCompletedAt returns whether the Task was completed at the given date.
func (task Task) WasCompletedAt(year int, month time.Month, day int) bool {
	mcmpl := task.MonthCompletionTime(year, month)
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

func (task Task) CurrentStrike() uint {
	if task.IsStrikeContinued() {
		return task.currentStrike
	}

	return 0
}

// IsStrikeContinued returns whether strike was broken
// meaning there was over 1 day break from finishing it.
func (task Task) IsStrikeContinued() bool {
	return AreSameDates(task.GetTime(), task.LastTimeCompleted) ||
		AreSameDates(task.GetTime().AddDate(0, 0, -1), task.LastTimeCompleted)
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
