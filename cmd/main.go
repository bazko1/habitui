package main

import (
	"fmt"
	"os"

	"github.com/bazko1/habitui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	habitui.NewTask("Work on habitui", `Daily work on personal project that is 
	also participating in coding challenge.`)

	p := tea.NewProgram(habitui.TuiModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
