package habitui_test

import (
	"testing"
	"time"

	"github.com/bazko1/habitui"
)

type dayIncreasingTime struct {
	CurrentTime time.Time
}

func (dit *dayIncreasingTime) Now() time.Time {
	return dit.CurrentTime
}

func (dit *dayIncreasingTime) AddDay() {
	dit.CurrentTime = dit.CurrentTime.AddDate(0, 0, 1)
}

func TestTaskCompletionSingleDay(t *testing.T) {
	t.Parallel()

	task := habitui.NewTask("test", "test description")

	task.MakeTaskCompleted()
	task.MakeTaskCompleted()

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
	task := habitui.NewTaskWithCustomTime("hit the gym", "test description", dit.Now)

	dit.AddDay()
	task.MakeTaskCompleted()

	dit.AddDay()
	task.MakeTaskCompleted()

	dit.AddDay()
	task.MakeTaskCompleted()

	year, month, _ := dit.CurrentTime.Date()
	monthCompletions := task.MonthCompletionTime(year, month)

	if len(monthCompletions) != 3 {
		t.Fatalf("Task completion is %d when it should be 3", len(monthCompletions))
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

	if task.CurrentWeekCompletion() != 3 {
		t.Fatalf("Task should be completed 3 times this week while it returned %d", task.CurrentStrike())
	}

	// this moves date to monday so the weekly counter should be 0 now
	dit.AddDay()

	if task.CurrentWeekCompletion() != 0 {
		t.Fatalf("Task should be completed 0 times this week while it returned %d", task.CurrentStrike())
	}

	task.MakeTaskCompleted()

	if task.CurrentMonthCompletion() != 4 {
		t.Fatalf("Task should be completed 4 times this month while it returned %d", task.CurrentStrike())
	}
}
