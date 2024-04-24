package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	habitui "github.com/bazko1/habitui/habit"
)

func main() {
	startDate := time.Date(2024, 3, 1, 12, 0, 0, 0, time.Local)
	now := func() time.Time {
		return startDate
	}
	tasks := habitui.TaskList{
		habitui.NewTaskWithCustomTime("work on habittui", "daily app grind", now),
		habitui.NewTaskWithCustomTime("go for a walk", "walking is relaxing and healthy activity", now),
	}

	for startDate.Before(time.Now().AddDate(0, 0, -2)) {
		startDate = startDate.AddDate(0, 0, 1)

		tasks[0].MakeCompleted()

		for i := 1; i < len(tasks); i++ {
			if rand.Int()%2 != 0 {
				tasks[i].MakeCompleted()
			}
		}
	}

	err := habitui.JSONSaveTasks(".habitui.json", tasks)
	if err != nil {
		fmt.Println("failed to save tasks state")
	}
}
