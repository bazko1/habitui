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

	if len(task.CompletionHistory) == 0 {
		t.Fatal("Task completion wasn't archived")
	}

	if len(task.CompletionHistory) > 1 {
		t.Fatal("Task shouldn't be completed twice a day")
	}

	if !task.WasCompletedToday() {
		t.Fatal("Task should return that it was completed today")
	}
}

func TestTaskWithChangingDay(t *testing.T) {
	t.Parallel()

	dit := dayIncreasingTime{time.Now()}
	task := habitui.NewTaskWithCustomTime("test", "test description", dit.Now)

	dit.AddDay()
	task.MakeTaskCompleted()

	dit.AddDay()
	task.MakeTaskCompleted()

	dit.AddDay()
	task.MakeTaskCompleted()

	if len(task.CompletionHistory) != 3 {
		t.Fatalf("Task completion is %d when it shoudld be 3", len(task.CompletionHistory))
	}
}
