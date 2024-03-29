package habit_test

import (
	"os"
	"testing"
	"time"

	habitui "github.com/bazko1/habitui/habit"
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

	if strike := task.CurrentMonthBestStrike(); strike != 3 {
		t.Fatalf("Task CurrentMonthBestStrike should be 3 times while it returned %d", strike)
	}
}

func TestTaskUnCompletion(t *testing.T) {
	t.Parallel()

	dit := dayIncreasingTime{time.Date(2023, time.October, 3, 15, 33, 0, 0, time.UTC)}
	task := habitui.NewTaskWithCustomTime("hit the gym", "test description", dit.Now)
	numCompletions := 6

	for range numCompletions - 1 {
		dit.AddDay()
		task.MakeTaskCompleted()
	}

	notUnCompleted := dit.CurrentTime
	dit.AddDay()
	task.MakeTaskCompleted()

	if task.CurrentMonthCompletion() != numCompletions {
		t.Fatalf("Task should be completed %d times this month while it returned %d", numCompletions, task.CurrentStrike())
	}

	if strike := task.CurrentMonthBestStrike(); strike != numCompletions {
		t.Fatalf("Task CurrentMonthBestStrike should be %d times while it returned %d", numCompletions, strike)
	}

	task.MakeTaskUnCompleted()

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

func TestTaskJSONState(t *testing.T) {
	t.Parallel()

	dit := dayIncreasingTime{time.Date(2023, time.October, 3, 15, 33, 0, 0, time.UTC)}
	tasks := habitui.TaskList{
		habitui.NewTaskWithCustomTime("go for a walk", "walkin and dreamin...", dit.Now),
		habitui.NewTaskWithCustomTime("strength training", "gym or home calistenics training", dit.Now),
		habitui.NewTaskWithCustomTime("english lesson", "mobile app lesson", dit.Now),
	}
	inARowCompl := 4

	for range inARowCompl {
		dit.AddDay()

		for i := range tasks {
			tasks[i].MakeTaskCompleted()
		}
	}

	for _, task := range tasks {
		if compl := task.CurrentMonthCompletion(); compl != inARowCompl {
			t.Fatalf("Task '%s' should be completed %d times this month while it returned %d",
				task.Name,
				inARowCompl,
				compl)
		}

		if strike := task.CurrentMonthBestStrike(); strike != inARowCompl {
			t.Fatalf("Task '%s' CurrentMonthBestStrike should be %d times while it returned %d", task.Name, inARowCompl, strike)
		}
	}

	file, err := os.CreateTemp("", "tmpfile-json-test")
	if err != nil {
		t.Fatalf("Failed to create tempfile: %v", err)
	}

	defer func() {
		os.Remove(file.Name())
	}()

	err = habitui.JSONSaveTasks(file.Name(), tasks)
	if err != nil {
		t.Fatalf("Failed to json save: %v", err)
	}

	loadedTasks, err := habitui.JSONLoadTasks(file.Name())
	if err != nil {
		t.Fatalf("Failed to load tasks from json: %v", err)
	}

	for _, task := range loadedTasks {
		task.GetTime = dit.Now
		if compl := task.CurrentMonthCompletion(); compl != inARowCompl {
			t.Fatalf("Task '%s' should be completed %d times this month while it returned %d",
				task.Name,
				inARowCompl,
				compl)
		}
	}
}
