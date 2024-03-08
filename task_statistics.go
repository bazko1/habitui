package habitui

import (
	"time"
)

const WeekDuration int = 7

// MonthlyTaskCompletion is a record of task completion for each month.
type MonthlyTaskCompletion map[time.Month][]time.Time

// YearlyTaskCompletion keeps a record of MonthlyTaskCompletion for each year.
type YearlyTaskCompletion map[int]MonthlyTaskCompletion

func (task Task) CurrentYearCompletion() int {
	return task.YearCompletion(task.GetTime().Year())
}

// Returns number of completions over the given year.
func (task Task) YearCompletion(year int) int {
	completionsYear, exists := task.YearlyTaskCompletion[year]
	if !exists {
		return 0
	}

	completions := 0
	for _, monthCompletions := range completionsYear {
		completions += len(monthCompletions)
	}

	return completions
}

// Returns number of completions over the month represented by given year and month.
func (task Task) MonthCompletion(year int, month time.Month) int {
	mcp := task.MonthCompletionTime(year, month)
	if mcp == nil {
		return 0
	}

	return len(mcp)
}

func (task Task) CurrentMonthCompletion() int {
	y, m, _ := task.GetTime().Date()

	return task.MonthCompletion(y, m)
}

// Returns number of completions over the week represented by given date.
// Week is previous Monday up to given date.
func (task Task) WeekCompletion(year int, month time.Month, day int) int {
	mcp := task.MonthCompletionTime(year, month)
	if mcp == nil {
		return 0
	}
	// we count all completions that are less than end
	end := time.Date(year, month, day, 23, 59, 59, 59, task.GetTime().Location())

	weekDay := int(end.Weekday())
	// counting Sunday as seventh day of the week
	if weekDay == 0 {
		weekDay = 7
	}

	weekBegin := time.Date(year, month, day, 0, 0, 0, 1, task.GetTime().Location()).AddDate(0, 0, -weekDay)

	counter := 0

	for _, cdt := range mcp {
		if cdt.Before(weekBegin) {
			continue
		}

		if cdt.After(end) {
			break
		}

		counter++
	}

	return counter
}

// CurrentWeekCompletion returns number of task completions over the whole week up to
// current day. The week in this sense is treated as last Monday till current day.
func (task Task) CurrentWeekCompletion() int {
	y, m, d := task.GetTime().Date()

	return task.WeekCompletion(y, m, d)
}
