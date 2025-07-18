package habit

import (
	"slices"
	"time"
)

const TaskVersionLatest = "v1"

// TaskList is a slice of tasks.
type TaskList []Task

// Task is an occurring event that has its own name identifier.
// Each task can be completed once a day.
type Task struct {
	Version      string
	Name         string
	Description  string
	CreationDate time.Time
	GetTime      func() time.Time `json:"-"`

	yearlyTaskCompletion YearlyTaskCompletion
	lastTimeCompleted    time.Time
	currentStrike        int
	// TODO: CurrentMonthBestStrike actually shows best strike ever that was finished in
	// current month. I think I would like it to show best strikes in current month
	// in a sense that it is value from 0-31.
	strikeThisMonth        Strike
	yearlyBestStrike       YearlyBestStrike
	bestStrikeLastFinished time.Time
}

// NewTask creates new task based on name and description with default get time function set to time.Now.
func NewTask(name, description string) Task {
	return WithCustomTime(name, description, time.Now)
}

// WithCustomTime constructs Task based on given name, description and time returning function.
func WithCustomTime(name, description string, getTime func() time.Time) Task {
	if getTime == nil {
		getTime = time.Now
	}

	return Task{
		Version:                TaskVersionLatest,
		Name:                   name,
		Description:            description,
		CreationDate:           getTime(),
		GetTime:                getTime,
		yearlyTaskCompletion:   make(YearlyTaskCompletion),
		lastTimeCompleted:      time.Time{},
		currentStrike:          0,
		strikeThisMonth:        Strike{},
		yearlyBestStrike:       make(YearlyBestStrike),
		bestStrikeLastFinished: time.Time{},
	}
}

// LastTimeCompleted returns last date when task completion was done.
func (task *Task) LastTimeCompleted() time.Time {
	return task.lastTimeCompleted
}

// MakeCompleted updates all task tracking states with information about finishing task now.
// Date is added to completion history if it wasn't completed this day yet.
// This method also updates statistics information such as day streak number.
func (task *Task) MakeCompleted() {
	now := task.GetTime()

	if task.WasCompletedToday() {
		return
	}

	if task.lastTimeCompleted.Year() != now.Year() {
		initializeDateMaps(
			task.yearlyTaskCompletion,
			now.Year(),
			now.Month(), func() []time.Time {
				longestHalfMonthUpper := 16
				init := make([]time.Time, 0, longestHalfMonthUpper)
				init = append(init, now)

				return init
			})

		initializeDateMaps(task.yearlyBestStrike, now.Year(),
			now.Month(), func() int { return 1 })
	}

	if task.lastTimeCompleted.IsZero() {
		task.lastTimeCompleted = now
		task.currentStrike = 1
		task.bestStrikeLastFinished = now

		return
	}

	if !task.IsStrikeContinued() {
		task.currentStrike = 1
	}

	if AreSameDates(now, task.lastTimeCompleted.AddDate(0, 0, 1)) {
		task.currentStrike++
	}

	monthBestStrike := task.yearlyBestStrike[now.Year()][now.Month()]

	if task.currentStrike > monthBestStrike {
		task.yearlyBestStrike[now.Year()][now.Month()] = task.currentStrike
		task.bestStrikeLastFinished = now
	}

	task.lastTimeCompleted = now

	completionsThisYear := task.yearlyTaskCompletion[now.Year()]
	completionsThisMonth := completionsThisYear[now.Month()]

	if len(completionsThisMonth) == 0 {
		completionsThisYear[now.Month()] = append(completionsThisMonth, now)

		return
	}

	if lastComplete := completionsThisMonth[len(completionsThisMonth)-1]; !AreSameDates(now, lastComplete) {
		completionsThisYear[now.Month()] = append(completionsThisMonth, now)
	}
}

// MakeUnCompleted makes reverts task completion for current day.
func (task *Task) MakeUnCompleted() {
	if !task.WasCompletedToday() {
		return
	}

	complDate := task.lastTimeCompleted
	task.lastTimeCompleted = time.Time{}

	if monthlyCompletions, exists := task.yearlyTaskCompletion[complDate.Year()]; exists {
		monthly := monthlyCompletions[complDate.Month()]
		if completionNum := len(monthly); completionNum > 0 && monthly[completionNum-1].Equal(complDate) {
			monthlyCompletions[complDate.Month()] = monthly[:completionNum-1]

			if completionNum > 1 {
				task.lastTimeCompleted = monthly[len(monthly)-2]
			}
		}
	}

	task.currentStrike--

	if task.bestStrikeLastFinished.Equal(complDate) {
		task.yearlyBestStrike[complDate.Year()][complDate.Month()] = task.currentStrike
	}
}

// MonthCompletionTime returns task completion at given year and month.
func (task *Task) MonthCompletionTime(year int, month time.Month) []time.Time {
	completionsYear, exists := task.yearlyTaskCompletion[year]
	if !exists {
		return nil
	}

	return completionsYear[month]
}

// WasCompletedAt returns whether the Task was completed at the given date.
func (task *Task) WasCompletedAt(year int, month time.Month, day int) bool {
	mcmpl := task.MonthCompletionTime(year, month)
	if mcmpl == nil {
		return false
	}

	atDate := time.Date(year, month, day, 0, 0, 0, 0, &time.Location{})

	return slices.ContainsFunc(mcmpl, func(t time.Time) bool { return AreSameDates(t, atDate) })
}

// WasCompletedToday returns whether the Task was completed at current GetTime day.
func (task *Task) WasCompletedToday() bool {
	return AreSameDates(task.GetTime(), task.lastTimeCompleted)
}

// CurrentStrike returns how many days in a row were task finished.
func (task *Task) CurrentStrike() int {
	if task.IsStrikeContinued() {
		return task.currentStrike
	}

	return 0
}

// IsStrikeContinued returns whether strike was broken
// meaning there was over 1 day break from finishing it.
func (task *Task) IsStrikeContinued() bool {
	return AreSameDates(task.GetTime(), task.lastTimeCompleted) ||
		AreSameDates(task.GetTime().AddDate(0, 0, -1), task.lastTimeCompleted)
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
