package main

import (
	"fmt"
	"os"

	habitui "github.com/bazko1/habitui/habit"
	"github.com/bazko1/habitui/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// \ TODO: this should alaso parsed from command line arguments
	// and searched in some standard places like $HOME/.config
	tasksFile := ".habitui.json"

	// TODO: we need to check if file exist first
	tasks, err := habitui.JSONLoadTasks(tasksFile)
	if err != nil {
		fmt.Println("failed to load tasks:", err) //nolint:forbidigo
		os.Exit(1)
	}

	prog := tea.NewProgram(tui.NewTuiAgent(tasks))

	defer func() {
		err := habitui.JSONSaveTasks(tasksFile, tasks)
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
