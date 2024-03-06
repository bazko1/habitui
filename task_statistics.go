package habitui

import (
	"time"
)

// MonthlyTaskCompletion is a record of task completion for each month.
type MonthlyTaskCompletion map[time.Month][]time.Time

// YearlyTaskCompletion keeps a record of MonthlyTaskCompletion for each year.
type YearlyTaskCompletion map[int]MonthlyTaskCompletion
