package tui

import (
	"fmt"
	"strings"

	"github.com/bazko1/habitui/habit"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	tasks       habit.TaskList
	cursorRow   int
	cursorCol   int
	selectedRow map[int]struct{}
	keys        keyMap
	help        help.Model
	addSelected bool
}

func NewTuiModel(tasks habit.TaskList) Model {
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
			Quit: key.NewBinding(
				key.WithKeys("q", "esc", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
		help:        help.New(),
		addSelected: false,
	}

	for tID, t := range tasks {
		if t.WasCompletedToday() {
			model.selectedRow[tID] = struct{}{}
		}
	}

	return model
}

func (model Model) Init() tea.Cmd {
	return nil
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Select key.Binding
	Help   key.Binding
	Add    key.Binding
	Quit   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},   // first column
		{k.Help, k.Quit, k.Select, k.Add}, // second column
	}
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint: ireturn, cyclop
	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
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
			if model.cursorCol < numWinCols {
				model.cursorCol++
			}

		case key.Matches(msg, model.keys.Left):
			if model.cursorCol > 0 {
				model.cursorCol--
			}

		case key.Matches(msg, model.keys.Select):
			_, ok := model.selectedRow[model.cursorRow]
			if ok {
				delete(model.selectedRow, model.cursorRow)
				model.tasks[model.cursorRow].MakeTaskUnCompleted()
			} else {
				model.selectedRow[model.cursorRow] = struct{}{}
				model.tasks[model.cursorRow].MakeTaskCompleted()
			}

		case key.Matches(msg, model.keys.Add):
			model.addSelected = true

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
	if model.addSelected {
		return addTaskView()
	}
	description := ""
	habits := strings.Builder{}
	selectedID := 0
	descriptionSelected := false
	height := 1

	habits.WriteString("Habits:\n")

	if model.cursorCol == 1 {
		descriptionSelected = true
	}

	for taskID, task := range model.tasks {
		height++
		taskName := task.Name

		if model.cursorRow == taskID {
			selectedID = taskID
			description = task.Description

			if model.cursorCol == 0 {
				taskName = formatSelectedText(taskName)
			}
		}

		completed := " "
		if _, ok := model.selectedRow[taskID]; ok {
			completed = "x"
		}

		habits.WriteString(fmt.Sprintf("[%s] %s\n", completed, taskName))
	}

	view := ""

	if len(model.tasks) == 0 {
		habits.WriteString("No habits.")

		description = "Add new task and start forming habit."
	}

	view += lipgloss.JoinHorizontal(1, createUpperTextPanelBox(strings.TrimSuffix(habits.String(), "\n"), height),
		createDescriptionBox(description, height, descriptionSelected))

	if len(model.tasks) != 0 {
		selectedTask := model.tasks[selectedID]
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
	}

	helpView := model.help.View(model.keys)
	view += "\n" + helpView

	return view
}

func addTaskView() string {
	buttonStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#888B7E")).
		Padding(0, 3).
		MarginTop(1)

	activeButtonStyle := buttonStyle.Copy().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#F25D94")).
		MarginRight(2).
		Underline(true)

	okButton := activeButtonStyle.Render("Yes")
	cancelButton := buttonStyle.Render("Maybe")

	question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Are you sure you want to eat marmalade?")
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	width := 2
	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 0).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)
	dialog := lipgloss.Place(width, 9,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars("猫咪"),

		lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}),
	)

	return dialog
}
