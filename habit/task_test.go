package habit_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bazko1/habitui/habit"
)

func habitCountError(hName string, period string, expected int, count int) error {
	return fmt.Errorf("Task '%s' should be completed %d times over the %s while it returned %d", //nolint:goerr113
		hName,
		expected,
		period,
		count)
}

func habitStrikeError(hName string, period string, expected int, count int) error {
	return fmt.Errorf("Task '%s' should have %s strike equal to %d while it returned %d", //nolint:goerr113
		hName,
		period,
		expected,
		count)
}

type dayIncreasingTime struct {
	CurrentTime time.Time
}

func (dit *dayIncreasingTime) Now() time.Time {
	return dit.CurrentTime
}

func (dit *dayIncreasingTime) AddDay() {
	dit.CurrentTime = dit.CurrentTime.AddDate(0, 0, 1)
}

func validateCompletion(habit habit.Task, expectedWeek int, expectedMonth int, expectedYear int) error {
	wct, mct, yct := habit.AllCompletion()

	if wct != expectedWeek {
		return habitCountError(habit.Name,
			"week",
			expectedWeek,
			wct)
	}

	if mct != expectedMonth {
		return habitCountError(habit.Name,
			"month",
			expectedMonth,
			mct)
	}

	if yct != expectedYear {
		return habitCountError(habit.Name,
			"year",
			expectedYear,
			yct)
	}

	return nil
}

func validateStrike(habit habit.Task, expectedCurrent, expectedMonth, expectedYear int) error {
	wct, mct, yct := habit.AllStrike()

	if wct != expectedCurrent {
		return habitStrikeError(habit.Name,
			"current",
			expectedCurrent,
			wct)
	}

	if mct != expectedMonth {
		return habitStrikeError(habit.Name,
			"best monthly",
			expectedMonth,
			mct)
	}

	if yct != expectedYear {
		return habitStrikeError(habit.Name,
			"best yearly",
			expectedYear,
			yct)
	}

	return nil
}

func TestTaskCompletionSingleDay(t *testing.T) {
	t.Parallel()

	task := habit.NewTask("test", "test description")

	task.MakeCompleted()
	task.MakeCompleted()

	y, m, _ := time.Now().Date()
	monthCompletions := task.MonthCompletionTime(y, m)

	if len(monthCompletions) == 0 {
		t.Fatal("Task completion wasn't archived")
	}

	if len(monthCompletions) > 1 {
		t.Fatal("Task shouldn't be completed twice a day")
	}

	if !task.WasCompletedToday() {
		t.Fatal("Task should return that it was completed today")
	}
}

func TestTaskWithChangingDay(t *testing.T) {
	t.Parallel()

	// new year new me resolution
	dit := dayIncreasingTime{time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)}
	task := habit.WithCustomTime("hit the gym", "test description", dit.Now)

	dit.AddDay()
	task.MakeCompleted()

	dit.AddDay()
	task.MakeCompleted()

	dit.AddDay()

	task.MakeCompleted()

	if compl := task.CurrentMonthCompletion(); compl != 3 {
		t.Fatalf("Task completion is %d when it should be 3", compl)
	}

	if task.CurrentStrike() != 3 {
		t.Fatalf("Task strike is %d when it should be 3", task.CurrentStrike())
	}

	dit.AddDay()
	dit.AddDay()
	// 2 days passed and task wasn't completed so strike should be zeroed.
	if task.CurrentStrike() != 0 {
		t.Fatalf("Task strike is %d when it should be 0", task.CurrentStrike())
	}

	dit.AddDay()
	dit.AddDay()
	// this moves date to sunday
	dit.AddDay()

	if weekCpl := task.CurrentWeekCompletion(); weekCpl != 2 {
		t.Fatalf("Task should be completed %d times this week while it returned %d", 2, weekCpl)
	}

	// this moves date to monday so the weekly counter should be 0 now
	dit.AddDay()

	if task.CurrentWeekCompletion() != 0 {
		t.Fatalf("Task should be completed 0 times this week while it returned %d", task.CurrentStrike())
	}

	task.MakeCompleted()

	if task.CurrentMonthCompletion() != 4 {
		t.Fatalf("Task should be completed 4 times this month while it returned %d", task.CurrentStrike())
	}

	if strike := task.CurrentMonthBestStrike(); strike != 3 {
		t.Fatalf("Task CurrentMonthBestStrike should be 3 times while it returned %d", strike)
	}
}

func TestTaskUnCompletion(t *testing.T) {
	t.Parallel()

	dit := dayIncreasingTime{time.Date(2023, time.October, 3, 15, 33, 0, 0, time.UTC)}
	task := habit.WithCustomTime("hit the gym", "test description", dit.Now)
	numCompletions := 6

	for range numCompletions - 1 {
		dit.AddDay()
		task.MakeCompleted()
	}

	notUnCompleted := dit.CurrentTime
	dit.AddDay()
	task.MakeCompleted()

	if task.CurrentMonthCompletion() != numCompletions {
		t.Fatalf("Task should be completed %d times this month while it returned %d", numCompletions, task.CurrentStrike())
	}

	if strike := task.CurrentMonthBestStrike(); strike != numCompletions {
		t.Fatalf("Task CurrentMonthBestStrike should be %d times while it returned %d", numCompletions, strike)
	}

	task.MakeUnCompleted()

	numCompletions--
	if task.CurrentMonthCompletion() != numCompletions {
		t.Fatalf("Task should be completed %d times this month while it returned %d", numCompletions, task.CurrentStrike())
	}

	if strike := task.CurrentMonthBestStrike(); strike != numCompletions {
		t.Fatalf("Task CurrentMonthBestStrike should be %d times while it returned %d", numCompletions, strike)
	}

	if last := task.LastTimeCompleted(); last != notUnCompleted {
		t.Fatalf("Uncompleting did not update last completion properly it is %v while should be %v", last, notUnCompleted)
	}
}

func TestCompletionChangingMonth(t *testing.T) {
	t.Parallel()
	// 2024-03-30
	startDate := time.Date(2024, 3, 30, 12, 0, 0, 0, time.Local)
	now := func() time.Time {
		return startDate
	}
	task := habit.WithCustomTime("work on habittui", "daily app grind", now)

	compl := 9

	// first check completions in March
	for c := range 2 {
		task.MakeCompleted()

		if err := validateCompletion(task, c+1, c+1, c+1); err != nil {
			t.Fatal(err.Error())
		}

		if err := validateStrike(task, c+1, c+1, c+1); err != nil {
			t.Fatal(err.Error())
		}

		startDate = startDate.AddDate(0, 0, 1)
	}

	// check completions in April
	for c := range compl - 2 {
		task.MakeCompleted()

		if err := validateCompletion(task, c+1, c+1, c+3); err != nil {
			t.Fatal(err.Error())
		}

		// TODO: Probably the monthly will change that it will return c+1 as
		// number of completions only in April will count.
		if err := validateStrike(task, c+3, c+3, c+3); err != nil {
			t.Fatal(err.Error())
		}

		startDate = startDate.AddDate(0, 0, 1)
	}

	// monday
	if err := validateCompletion(task, 0, compl-2, compl); err != nil {
		t.Fatal(err.Error())
	}
}

func TestCompletionChangingYear(t *testing.T) {
	t.Parallel()
	// 2023-12-26
	startDate := time.Date(2024, 12, 26, 8, 0, 0, 0, time.Local)
	now := func() time.Time {
		return startDate
	}
	task := habit.WithCustomTime("work on habittui", "daily app grind", now)
	compl := 10

	for range compl {
		task.MakeCompleted()

		startDate = startDate.AddDate(0, 0, 1)
	}

	// 2025-01-05 and not completed yet
	if err := validateStrike(task, compl, compl, compl); err != nil {
		t.Fatal(err.Error())
	}

	if err := validateCompletion(task, 6, 4, 4); err != nil {
		t.Fatal(err.Error())
	}
}
