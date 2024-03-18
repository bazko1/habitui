package habitui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	sectionBoxWidth = 40
	numWinCols      = 2
)

type TuiAgent struct {
	tasks       TaskList
	cursorRow   int
	cursorCol   int
	selectedRow map[int]struct{}
}

func NewTuiAgent(tasks TaskList) TuiAgent {
	return TuiAgent{
		tasks:       tasks,
		cursorRow:   0,
		cursorCol:   0,
		selectedRow: make(map[int]struct{}),
	}
}

func (agent TuiAgent) Init() tea.Cmd {
	return nil
}

func (agent TuiAgent) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint: ireturn, cyclop
	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return agent, tea.Quit

		case "up", "k":
			if agent.cursorRow > 0 && agent.cursorCol == 0 {
				agent.cursorRow--
			}

		case "down", "j":
			if agent.cursorRow < len(agent.tasks)-1 && agent.cursorCol == 0 {
				agent.cursorRow++
			}

		case "right", "l":
			if agent.cursorCol < numWinCols {
				agent.cursorCol++
			}

		case "left", "h":
			if agent.cursorCol > 0 {
				agent.cursorCol--
			}

		case "enter", " ":
			_, ok := agent.selectedRow[agent.cursorRow]
			if ok {
				delete(agent.selectedRow, agent.cursorRow)
			} else {
				agent.selectedRow[agent.cursorRow] = struct{}{}
			}
		}
	}

	return agent, nil
}

func formatSelectedText(text string) string {
	style := lipgloss.NewStyle().
		Bold(true).
		Italic(true).
		Faint(true).
		Reverse(true)

	return style.Render(text)
}

func createUpperTextPanelBox(text string, height int) string {
	style := lipgloss.NewStyle().
		// Foreground(lipgloss.Color("#FAFAFA")).
		// Background(lipgloss.Color("#7D56F4")).
		PaddingTop(0).
		PaddingLeft(0).
		Border(lipgloss.NormalBorder()).
		Height(height).
		Width(sectionBoxWidth)

	return style.Render(text)
}

func createLowerPanelTextBox(text string, height int) string {
	style := lipgloss.NewStyle().
		// Foreground(lipgloss.Color("#FAFAFA")).
		// Background(lipgloss.Color("#7D56F4")).
		PaddingTop(0).
		PaddingLeft(0).
		Border(lipgloss.NormalBorder()).
		Height(height).
		Width(sectionBoxWidth)

	return style.Render(text)
}

func (agent TuiAgent) View() string {
	description := ""
	habits := ""
	selectedID := 0

	for taskID, task := range agent.tasks {
		taskName := task.Name

		if agent.cursorRow == taskID {
			selectedID = taskID
			description = task.Description

			switch agent.cursorCol {
			case 0:
				taskName = formatSelectedText(taskName)
			case 1:
				description = formatSelectedText(description)
			}
		}

		completed := " "
		if _, ok := agent.selectedRow[taskID]; ok {
			completed = "x"
		}

		habits += fmt.Sprintf("[%s] %s\n", completed, taskName)
	}

	view := ""
	habits = "Habits:\n" + habits[:len(habits)-1]
	height := strings.Count(habits, "\n")

	// /TODO: Need to format description so that if its longer than some characers
	// newlines need to be added
	view += lipgloss.JoinHorizontal(1, createUpperTextPanelBox(habits, height),
		createUpperTextPanelBox("Description:\n"+description, height+1))

	numOfStats := 4
	lowerPanel := lipgloss.JoinHorizontal(
		1,
		createLowerPanelTextBox(fmt.Sprintf("Strike (task %d):\n\tCurrent: 0\n\tBest monthly: 0\n\tLongest: 0",
			selectedID), numOfStats),
		createLowerPanelTextBox(
			fmt.Sprintf("Completion (task %d):\n\tThis week: 0\n\tThis month: 0\n\tThis year: 0", selectedID),
			numOfStats,
		),
	)

	view = lipgloss.JoinVertical(1, view, lowerPanel)

	// The footer
	view += "\n\nPress q to quit.\n"

	// Send the UI for rendering
	return view
}
