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

	var tasks habit.TaskList

	file, err := os.ReadFile(*tasksFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("failed to open tasks file '%s': %v\n", *tasksFile, err) //nolint:forbidigo
		os.Exit(1)
	}

	if !errors.Is(err, os.ErrNotExist) {
		tasks, err = habit.JSONLoadTasks(file)
		if err != nil {
			fmt.Println("failed to load tasks:", err) //nolint:forbidigo
			os.Exit(1)
		}
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
