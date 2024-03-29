package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/bazko1/habitui/habit"
	"github.com/bazko1/habitui/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	tasksFile := flag.String("data", ".habitui.json", "a name of for loading/saving tasks data")
	flag.Parse()

	tasks, err := habit.JSONLoadTasks(*tasksFile)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Println("failed to load tasks:", err) //nolint:forbidigo
		os.Exit(1)
	}

	if tasks == nil {
		tasks = habit.TaskList{}
	}

	prog := tea.NewProgram(tui.NewTuiModel(tasks))

	defer func() {
		err := habit.JSONSaveTasks(*tasksFile, tasks)
		if err != nil {
			fmt.Println("failed to save tasks: %w", err) //nolint:forbidigo
			os.Exit(1)
		}

		os.Exit(0)
	}()

	if _, err := prog.Run(); err != nil {
		fmt.Printf("Running tui error: %v", err) //nolint:forbidigo
	}
}
