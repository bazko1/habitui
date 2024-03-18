package main

import (
	"fmt"
	"os"

	"github.com/bazko1/habitui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	tasks := habitui.TaskList{
		habitui.NewTask("go for a walk", "walkin and dreamin..."),
		habitui.NewTask("strength training", "gym or home calistenics training"),
		habitui.NewTask("english lesson", "mobile app lesson"),
		habitui.NewTask("Work on habitui", "Daily work on personal project that is "+
			"also participating in coding challenge."),
	}

	p := tea.NewProgram(habitui.NewTuiAgent(tasks))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Running tui error: %v", err)
		os.Exit(1)
	}
}
