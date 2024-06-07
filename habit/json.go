package habit

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type taskJSON struct {
	Version      string
	Name         string
	Description  string
	CreationDate time.Time

	YearlyTaskCompletion   YearlyTaskCompletion
	LastTimeCompleted      time.Time
	CurrentStrike          int
	BestStrikeThisWeek     int
	StrikeThisMonth        Strike
	YearlyBestStrike       YearlyBestStrike
	BestStrikeLastFinished time.Time
}

func (t *TaskList) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), t)
}

func (t TaskList) Value() (driver.Value, error) {
	b, err := json.Marshal(t)

	return string(b), err
}

func (task Task) MarshalJSON() ([]byte, error) {
	if task.Version == "" {
		task.Version = TaskVersionLatest
	}

	bytes, err := json.Marshal(taskJSON{
		Version:                task.Version,
		Name:                   task.Name,
		Description:            task.Description,
		CreationDate:           task.CreationDate,
		YearlyTaskCompletion:   task.yearlyTaskCompletion,
		LastTimeCompleted:      task.lastTimeCompleted,
		CurrentStrike:          task.currentStrike,
		StrikeThisMonth:        task.strikeThisMonth,
		YearlyBestStrike:       task.yearlyBestStrike,
		BestStrikeLastFinished: task.bestStrikeLastFinished,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Task: %w", err)
	}
	// TODO: For now I stick with extra struct that capitalizes
	// all the fields. I might want to write un/marshaler that
	// converts struct to map[string]any and loops over field
	// names and capitalizes them so that I do not need to maintain
	// both taskJSON and Task.

	return bytes, nil
}

func (task *Task) UnmarshalJSON(data []byte) error {
	jsonTask := taskJSON{}
	if err := json.Unmarshal(data, &jsonTask); err != nil {
		return fmt.Errorf("failed to unmarshal Task: %w", err)
	}

	task.Version = jsonTask.Version
	task.Name = jsonTask.Name
	task.Description = jsonTask.Description
	task.CreationDate = jsonTask.CreationDate
	task.GetTime = time.Now
	task.yearlyTaskCompletion = jsonTask.YearlyTaskCompletion
	task.lastTimeCompleted = jsonTask.LastTimeCompleted
	task.currentStrike = jsonTask.CurrentStrike
	task.yearlyBestStrike = jsonTask.YearlyBestStrike
	task.bestStrikeLastFinished = jsonTask.BestStrikeLastFinished

	if task.Version == "" {
		task.Version = TaskVersionLatest
	}

	return nil
}

func JSONLoadTasks(bytes []byte) (TaskList, error) {
	taskList := TaskList{}
	if err := json.Unmarshal(bytes, &taskList); err != nil {
		return nil, fmt.Errorf("load json failed to umarshal bytes: %w", err)
	}

	return taskList, nil
}

func JSONSaveTasks(filename string, tasks TaskList) error {
	bytes, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("save json failed to marshall tasks(%v): %w", tasks, err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("save json failed to create file(%s): %w", filename, err)
	}

	defer func() {
		file.Close()
	}()

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("save json failed to write to file(%s): %w", filename, err)
	}

	return nil
}
