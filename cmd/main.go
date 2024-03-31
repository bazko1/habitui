package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
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

	// TODO: This probably should be optional
	// and if not log file created the debug logs should not be printed
	// or should be forwarded to dev null or somewhere
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err) //// nolint:forbidigo
		os.Exit(1)
	}
	defer f.Close()

	logger := log.Default()
	logger.Println("starting tui program")

	model := tui.NewTuiModel(tasks)
	prog := tea.NewProgram(model)

	out, err := prog.Run()
	if err != nil {
		logger.Printf("Running tui error: %v", err)
	}

	model, _ = out.(tui.Model)

	defer func() {
		err := habit.JSONSaveTasks(*tasksFile, model.Tasks())
		if err != nil {
			logger.Printf("failed to save tasks: %v", err)
			os.Exit(1)
		}

		os.Exit(0)
	}()
}
