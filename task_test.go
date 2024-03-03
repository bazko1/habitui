package habitui_test

import (
	"testing"

	"github.com/bazko1/habitui"
)

func TestTaskCompletion(t *testing.T) {
	t.Parallel()

	task := habitui.NewTask("test", "test description")

	task.MakeTaskCompleted()
	task.MakeTaskCompleted()

	if len(task.GetCompletionHistory()) == 0 {
		t.Fatal("Task completion wasn't archived")
	}

	if len(task.GetCompletionHistory()) > 1 {
		t.Fatal("Task shouldn't be completed twice a day")
	}

	if !task.WasCompletedToday() {
		t.Fatal("Task should return that it was completed today")
	}
}
