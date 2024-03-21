package tui

import (
	"fmt"
	"strings"

	"github.com/bazko1/habitui/habit"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
)

const (
	sectionBoxWidth = 40
	numWinCols      = 2
)

type Agent struct {
	tasks       habit.TaskList
	cursorRow   int
	cursorCol   int
	selectedRow map[int]struct{}
}

func NewTuiAgent(tasks habit.TaskList) Agent {
	agent := Agent{
		tasks:       tasks,
		cursorRow:   0,
		cursorCol:   0,
		selectedRow: make(map[int]struct{}),
	}

	for tID, t := range tasks {
		if t.WasCompletedToday() {
			agent.selectedRow[tID] = struct{}{}
		}
	}

	return agent
}

func (agent Agent) Init() tea.Cmd {
	return nil
}

func (agent Agent) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint: ireturn, cyclop
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
				agent.tasks[agent.cursorRow].MakeTaskUnCompleted()
			} else {
				agent.selectedRow[agent.cursorRow] = struct{}{}
				agent.tasks[agent.cursorRow].MakeTaskCompleted()
			}
		}
	}

	return agent, nil
}

func formatSelectedText(text string) string {
	style := lipgloss.NewStyle().
		Bold(true).
		Faint(true).
		Reverse(true)

	return style.Render(text)
}

func createUpperTextPanelBox(text string, height int) string {
	style := lipgloss.NewStyle().
		PaddingTop(0).
		PaddingLeft(0).
		Border(lipgloss.NormalBorder()).
		Height(height).
		Width(sectionBoxWidth)

	return style.Render(text)
}

func createDescriptionBox(desc string, height int, selected bool) string {
	height++
	style := lipgloss.NewStyle().
		PaddingTop(0).
		PaddingLeft(0).
		Border(lipgloss.NormalBorder()).
		Height(height).
		Width(sectionBoxWidth)

	// width based wrapping seems to
	// work incorrect if we want to
	// format text previously
	desc = wordwrap.String(desc, sectionBoxWidth)
	desc = wrap.String(desc, sectionBoxWidth) // force-wrap long strings

	if selected {
		desc = formatSelectedText(desc)
	}

	return style.Render("Description:\n" + desc)
}

func createLowerPanelTextBox(text string, height int) string {
	style := lipgloss.NewStyle().
		PaddingTop(0).
		PaddingLeft(0).
		Border(lipgloss.NormalBorder()).
		Height(height).
		Width(sectionBoxWidth)

	return style.Render(text)
}

func (agent Agent) View() string {
	description := ""
	habits := ""
	selectedID := 0
	descriptionSelected := false

	if agent.cursorCol == 1 {
		descriptionSelected = true
	}

	for taskID, task := range agent.tasks {
		taskName := task.Name

		if agent.cursorRow == taskID {
			selectedID = taskID
			description = task.Description

			if agent.cursorCol == 0 {
				taskName = formatSelectedText(taskName)
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

	view += lipgloss.JoinHorizontal(1, createUpperTextPanelBox(habits, height),
		createDescriptionBox(description, height, descriptionSelected))

	selectedTask := agent.tasks[selectedID]
	numOfStats := 4
	lowerPanel := lipgloss.JoinHorizontal(
		1,
		createLowerPanelTextBox(fmt.Sprintf("Strike:\n\tCurrent: %d\n\tBest monthly: %d\n\tBest yearly: %d",
			selectedTask.CurrentStrike(),
			selectedTask.CurrentMonthBestStrike(),
			selectedTask.CurrentYearBestStrike()), numOfStats),

		createLowerPanelTextBox(
			fmt.Sprintf("Completion:\n\tThis week: %d\n\tThis month: %d\n\tThis year: %d",
				selectedTask.CurrentWeekCompletion(),
				selectedTask.CurrentMonthCompletion(),
				selectedTask.CurrentYearCompletion()), numOfStats),
	)

	view = lipgloss.JoinVertical(1, view, lowerPanel)

	// The footer
	view += "\n\nPress q to quit.\n"

	// Send the UI for rendering
	return view
}
