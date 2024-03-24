package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/bazko1/habitui/habit"
	"github.com/bazko1/habitui/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// \ TODO: this should alaso parsed from command line arguments
	// and searched in some standard places like $HOME/.config
	tasksFile := ".habitui.json"

	tasks, err := habit.JSONLoadTasks(tasksFile)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Println("failed to load tasks:", err) //nolint:forbidigo
		os.Exit(1)
	}

	if tasks == nil {
		tasks = habit.TaskList{}
	}

	prog := tea.NewProgram(tui.NewTuiAgent(tasks))

	defer func() {
		err := habit.JSONSaveTasks(tasksFile, tasks)
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
