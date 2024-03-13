package habitui

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type JSONTask struct {
	*Task
	YearlyTaskCompletion YearlyTaskCompletion
	LastTimeCompleted    time.Time
	CurrentStrike        uint
	YearlyBestStrike     YearlyBestStrike
}

func (task Task) ToJSONTask() JSONTask {
	return JSONTask{
		&task,
		task.yearlyTaskCompletion,
		task.lastTimeCompleted,
		task.currentStrike,
		task.yearlyBestStrike,
	}
}

func JSONLoadTasks(filename string) (TaskList, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("load json failed to read file(%s): %w", filename, err)
	}

	jsonTasks := []JSONTask{}

	if err := json.Unmarshal(bytes, &jsonTasks); err != nil {
		return nil, fmt.Errorf("load json failed to umarshal bytes: %w", err)
	}

	taskList := make(TaskList, 0, len(jsonTasks))
	for _, jsonTask := range jsonTasks {
		taskList = append(taskList, Task{
			jsonTask.Name,
			jsonTask.Description,
			jsonTask.CreationDate,
			time.Now,
			jsonTask.YearlyTaskCompletion,
			jsonTask.LastTimeCompleted,
			jsonTask.CurrentStrike,
			jsonTask.YearlyBestStrike,
		})
	}

	return taskList, nil
}

func JSONSaveTasks(filename string, tasks TaskList) error {
	exportableTasks := make([]JSONTask, 0, len(tasks))
	for _, t := range tasks {
		exportableTasks = append(exportableTasks, t.ToJSONTask())
	}

	bytes, err := json.Marshal(exportableTasks)
	if err != nil {
		return fmt.Errorf("save json failed to marshall tasks(%v): %w", tasks, err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("save json failed to create file(%s): %w", filename, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("save json failed to write to file(%s): %w", filename, err)
	}

	return nil
}
