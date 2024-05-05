package habit_test

import (
	"os"
	"testing"
	"time"

	"github.com/bazko1/habitui/habit"
)

func TestTaskJSONState(t *testing.T) { //nolint:funlen
	t.Parallel()

	dit := dayIncreasingTime{time.Date(2023, time.October, 3, 15, 33, 0, 0, time.UTC)}
	tasks := habit.TaskList{
		habit.WithCustomTime("go for a walk", "walkin and dreamin...", dit.Now),
		habit.WithCustomTime("strength training", "gym or home calisthenics training", dit.Now),
		habit.WithCustomTime("english lesson", "mobile app lesson", dit.Now),
	}
	inARowCompl := 4

	for range inARowCompl {
		dit.AddDay()

		for i := range tasks {
			tasks[i].MakeCompleted()
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
		file.Close()
		os.Remove(file.Name())
	}()

	err = habit.JSONSaveTasks(file.Name(), tasks)
	if err != nil {
		t.Fatalf("Failed to json save: %v", err)
	}

	bytes, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	loadedTasks, err := habit.JSONLoadTasks(bytes)
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

func TestJSONLoadPartialData(t *testing.T) {
	t.Parallel()

	_, err := habit.JSONLoadTasks([]byte(`[{"Name":"go for a walk",
"Description":"walking is relaxing and healthy activity",
"CreationDate":"2024-03-01T12:00:00+01:00",
"YearlyTaskCompletion":{"2024":{"3":["2024-03-02T12:00:00+01:00"]}}}]`))
	if err != nil {
		t.Fatalf("Failed to load tasks from json: %v", err)
	}
}
