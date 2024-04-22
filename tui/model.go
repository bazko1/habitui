package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/bazko1/habitui/habit"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
)

const (
	sectionBoxWidth = 40
	numWinCols      = 2
)

type Model struct {
	tasks         habit.TaskList
	cursorRow     int
	cursorCol     int
	selectedRow   map[int]struct{}
	keys          keyMap
	editEnabled   bool
	addingNewTask bool
	editInput     textinput.Model
	help          help.Model
}

func NewTuiModel(tasks habit.TaskList) Model { //nolint:funlen
	editInput := textinput.New()
	model := Model{
		tasks:       tasks,
		cursorRow:   0,
		cursorCol:   0,
		selectedRow: make(map[int]struct{}),
		keys: keyMap{
			Up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "move up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "move down"),
			),
			Left: key.NewBinding(
				key.WithKeys("left", "h"),
				key.WithHelp("←/h", "move left"),
			),
			Right: key.NewBinding(
				key.WithKeys("right", "l"),
				key.WithHelp("→/l", "move right"),
			),
			Help: key.NewBinding(
				key.WithKeys("?"),
				key.WithHelp("?", "toggle help"),
			),
			Select: key.NewBinding(
				key.WithKeys("enter", " "),
				key.WithHelp("Enter/Space", "change task status"),
			),
			Add: key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "add task"),
			),
			Delete: key.NewBinding(
				key.WithKeys("d"),
				key.WithHelp("d", "delete task"),
			),
			Edit: key.NewBinding(
				key.WithKeys("e"),
				key.WithHelp("e", "edit task short name or description"),
			),
			Quit: key.NewBinding(
				key.WithKeys("q", "esc", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
		editEnabled:   false,
		addingNewTask: false,
		editInput:     editInput,
		help:          help.New(),
	}

	for tID, t := range tasks {
		if t.WasCompletedToday() {
			model.selectedRow[tID] = struct{}{}
		}
	}

	return model
}

func (model Model) Init() tea.Cmd {
	return textinput.Blink
}

func (model Model) Tasks() habit.TaskList {
	return model.tasks
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Select key.Binding
	Help   key.Binding
	Add    key.Binding
	Delete key.Binding
	Edit   key.Binding
	Quit   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Help, k.Quit, k.Select, k.Add},
		{k.Edit, k.Delete},
	}
}

var editMap = struct { //nolint:gochecknoglobals
	Quit    key.Binding
	Confirm key.Binding
}{
	key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "quit")),
	key.NewBinding(
		key.WithKeys("enter", "ctrl+j"),
		key.WithHelp("enter", "confirm")),
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint: ireturn, funlen, cyclop,gocognit
	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		if model.editEnabled {
			switch {
			case key.Matches(msg, editMap.Quit):
				model.editEnabled = false
				model.editInput.Blur()
				model.editInput.Reset()
			case key.Matches(msg, editMap.Confirm):
				model.editEnabled = false

				if model.cursorCol == 0 {
					model.tasks[model.cursorRow].Name = model.editInput.Value()
					if model.addingNewTask {
						model.editEnabled = true
						model.cursorCol = 1
						model.addingNewTask = false
					}
				} else {
					model.tasks[model.cursorRow].Description = model.editInput.Value()
					model.cursorCol = 0
					model.editInput.Blur()
				}

				model.editInput.Reset()
			}

			var cmd tea.Cmd
			model.editInput, cmd = model.editInput.Update(msg)

			return model, cmd
		}

		switch {
		case key.Matches(msg, model.keys.Quit):
			return model, tea.Quit

		case key.Matches(msg, model.keys.Up):
			if model.cursorRow > 0 && model.cursorCol == 0 {
				model.cursorRow--
			}

		case key.Matches(msg, model.keys.Down):
			if model.cursorRow < len(model.tasks)-1 && model.cursorCol == 0 {
				model.cursorRow++
			}

		case key.Matches(msg, model.keys.Right):
			if model.cursorCol < numWinCols && len(model.tasks) > 0 {
				model.cursorCol++
			}

		case key.Matches(msg, model.keys.Left):
			if model.cursorCol > 0 && len(model.tasks) > 0 {
				model.cursorCol--
			}

		case key.Matches(msg, model.keys.Select):
			_, ok := model.selectedRow[model.cursorRow]
			if ok {
				delete(model.selectedRow, model.cursorRow)
				model.tasks[model.cursorRow].MakeUnCompleted()
			} else {
				model.selectedRow[model.cursorRow] = struct{}{}
				model.tasks[model.cursorRow].MakeCompleted()
			}

		case key.Matches(msg, model.keys.Add):
			model.tasks = append(model.tasks, habit.NewTask("Set name", "Set description"))
			model.cursorRow = len(model.tasks) - 1
			model.editEnabled = true
			model.addingNewTask = true
			model.editInput.Focus()
			model.editInput.Cursor.Blink = true

		case key.Matches(msg, model.keys.Delete):
			model.tasks = append(model.tasks[:model.cursorRow], model.tasks[model.cursorRow+1:]...)

			if model.cursorRow > 0 {
				model.cursorRow--
			}
		case key.Matches(msg, model.keys.Edit):
			if len(model.tasks) == 0 {
				break
			}

			model.editEnabled = !model.editEnabled
			if model.editEnabled && !model.editInput.Focused() {
				model.editInput.Focus()
				model.editInput.Cursor.Blink = true
			}
		case key.Matches(msg, model.keys.Help):
			model.help.ShowAll = !model.help.ShowAll
		}
	}

	return model, nil
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

func (model Model) View() string { //nolint:funlen
	description := ""
	habits := strings.Builder{}
	selectedID := 0
	descriptionSelected := false
	height := 1

	habits.WriteString("Habits:\n")

	for taskID, task := range model.tasks {
		height++
		taskName := task.Name

		if model.cursorRow == taskID {
			selectedID = taskID
			description = task.Description

			if model.cursorCol == 0 {
				if model.editEnabled {
					model.editInput.Placeholder = taskName
					taskName = model.editInput.View()
				} else {
					taskName = formatSelectedText(taskName)
				}
			}
		}

		completed := " "
		if _, ok := model.selectedRow[taskID]; ok {
			completed = "x"
		}

		habits.WriteString(fmt.Sprintf("[%s] %s\n", completed, taskName))
	}

	if model.cursorCol == 1 {
		if model.editEnabled {
			model.editInput.Placeholder = description
			description = model.editInput.View()
		} else {
			descriptionSelected = true
		}
	}

	view := ""

	if len(model.tasks) == 0 {
		habits.WriteString("No habits.")

		description = "Add new task and start forming habit."
	}

	view += lipgloss.JoinHorizontal(1, createUpperTextPanelBox(strings.TrimSuffix(habits.String(), "\n"), height),
		createDescriptionBox(description, height, descriptionSelected))

	if len(model.tasks) != 0 {
		dateToday := model.tasks[0].GetTime()
		selectedTask := model.tasks[selectedID]
		numOfStats := 4
		lowerPanel := lipgloss.JoinHorizontal(
			1,
			createLowerPanelTextBox(fmt.Sprintf("Strike:\n\tCurrent: %d\n\tBest monthly: %d\n\tBest yearly: %d",
				selectedTask.CurrentStrike(),
				selectedTask.CurrentMonthBestStrike(),
				selectedTask.CurrentYearBestStrike()), numOfStats),

			createLowerPanelTextBox(
				fmt.Sprintf("Completion:\n\tThis week: %d / 7\n\tThis month: %d / %d \n\tThis year: %d / %d",
					selectedTask.CurrentWeekCompletion(),
					selectedTask.CurrentMonthCompletion(), getDaysInMonth(dateToday),
					selectedTask.CurrentYearCompletion(), getDaysInYear(dateToday)),
				numOfStats),
		)

		view = lipgloss.JoinVertical(1, view, lowerPanel)
	}

	helpView := model.help.View(model.keys)
	view += "\n" + helpView

	return view
}

func getDaysInMonth(t time.Time) int {
	return 32 - time.Date(t.Year(), t.Month(), 32, 0, 0, 0, 0, time.UTC).Day()
}

func getDaysInYear(t time.Time) int {
	return time.Date(t.Year(), 12, 31, 0, 0, 0, 0, time.UTC).YearDay()
}
