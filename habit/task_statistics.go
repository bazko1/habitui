package habit

import (
	"time"
)

const (
	numberOfMonths int = 12
)

// MonthlyTaskCompletion is a record of task completion for each month.
type MonthlyTaskCompletion map[time.Month][]time.Time

// YearlyTaskCompletion keeps a record of MonthlyTaskCompletion for each year.
type YearlyTaskCompletion map[int]MonthlyTaskCompletion

// MonthlyBestStrike is type for storing task best longest strike
// that happened over a month.
type MonthlyBestStrike map[time.Month]int

// YearlyBestStrike stores a record of MonthlyBestStrike
// for each year.
type YearlyBestStrike map[int]MonthlyBestStrike

func (task *Task) CurrentYearCompletion() int {
	return task.YearCompletion(task.GetTime().Year())
}

// AllCompletion returns completion for current week, month and year.
func (task *Task) AllCompletion() (int, int, int) {
	return task.CurrentWeekCompletion(),
		task.CurrentMonthCompletion(),
		task.CurrentYearCompletion()
}

// AllStrike returns current, best weekly and best yearly strike information.
func (task *Task) AllStrike() (int, int, int) {
	return task.CurrentStrike(),
		task.CurrentMonthBestStrike(),
		task.CurrentYearBestStrike()
}

// Returns number of completions over the given year.
func (task *Task) YearCompletion(year int) int {
	completionsYear, exists := task.yearlyTaskCompletion[year]
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
func (task *Task) MonthCompletion(year int, month time.Month) int {
	mcp := task.MonthCompletionTime(year, month)
	if mcp == nil {
		return 0
	}

	return len(mcp)
}

func (task *Task) CurrentMonthCompletion() int {
	y, m, _ := task.GetTime().Date()

	return task.MonthCompletion(y, m)
}

// Returns number of completions over the week represented by given date.
// Week is previous Monday up to given date.
func (task *Task) WeekCompletion(year int, month time.Month, day int) int {
	mcp := task.MonthCompletionTime(year, month)
	if mcp == nil {
		return 0
	}
	// we count all completions that are less than end
	// a TODO: To be checked if time here and in weekBegin is proper with edge cases.
	end := time.Date(year, month, day, 23, 59, 59, 59, task.GetTime().Location())

	weekDay := int(end.Weekday())
	// counting Sunday as seventh day of the week
	if weekDay == 0 {
		weekDay = 7
	}

	weekBegin := time.Date(year, month, day, 0, 0, 0, 0, task.GetTime().Location()).AddDate(0, 0, 1-weekDay)

	// handle overlapping months
	if weekBegin.Month() != end.Month() {
		mcp = append(task.MonthCompletionTime(weekBegin.Year(), weekBegin.Month()), mcp...)
	}

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
func (task *Task) CurrentWeekCompletion() int {
	y, m, d := task.GetTime().Date()

	return task.WeekCompletion(y, m, d)
}

func (task *Task) YearBestStrike(year int) int {
	monthlyStrikes, exist := task.yearlyBestStrike[year]
	if !exist {
		return 0
	}

	var max int

	for _, strike := range monthlyStrikes {
		if strike > max {
			max = strike
		}
	}

	return max
}

// TODO: CurrentMonthBestStrike actually shows best strike ever that happens to also
// have completions in given month. This probably should be refactored so that month best
// strike means longest completion over period of ~30 days (given month) and returned values
// should be in range 0-31.
func (task *Task) MonthBestStrike(year int, month time.Month) int {
	monthlyStrikes, exist := task.yearlyBestStrike[year]
	if !exist {
		return 0
	}

	return monthlyStrikes[month]
}

func (task *Task) CurrentMonthBestStrike() int {
	y, m, _ := task.GetTime().Date()

	return task.MonthBestStrike(y, m)
}

func (task *Task) CurrentYearBestStrike() int {
	return task.YearBestStrike(task.GetTime().Year())
}

// initializeDateMaps checks if YearlyBestStrike or YearlyTaskHistory
// are properly initialized and if not does that with given value init function.
func initializeDateMaps[Y ~map[int]M, M ~map[time.Month]V,
	V int | []time.Time](
	yearlyHistory Y, year int,
	month time.Month, initFunction func() V,
) {
	yearlyHistory[year] = make(M, numberOfMonths)
	yearlyHistory[year][month] = initFunction()
}
