package habitui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListModel struct {
	choices      []string // items on the to-do list
	descriptions []string
	cursor       int              // which to-do list item our cursor is pointing at
	selected     map[int]struct{} // which to-do items are selected
}

func TuiModel() ListModel {
	return ListModel{
		// Our to-do list is a grocery list
		choices: []string{"work on habitui", "go for a walk", "app english lesson"},
		descriptions: []string{
			"Longer description for\nthe task 1",
			"Longer description for\nthe task 2",
			"Longer description for\nthe task 3",
		},

		// A map which indicates which choices are selected. We're using
		// the map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
		cursor:   0,
	}
}

func (m ListModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint: ireturn
	switch msg := msg.(type) { //nolint: gocritic
	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
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
		Width(40)

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
		Width(40)

	return style.Render(text)
}

func (m ListModel) View() string {
	// The header
	habits := ""

	description := ""
	selectedID := 0
	// Iterate over our choices
	for chID, choice := range m.choices {
		// Is the cursor pointing at this choice?
		if m.cursor == chID {
			selectedID = chID
			description = fmt.Sprintf("%s", m.descriptions[chID])
			choice = formatSelectedText(choice)
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[chID]; ok {
			checked = "x" // selected!
		}

		// Render the row
		habits += fmt.Sprintf("[%s] %s\n", checked, choice)
	}
	view := ""
	habits = "Habits:\n" + habits[:len(habits)-1]
	height := strings.Count(habits, "\n")
	view += lipgloss.JoinHorizontal(1, createUpperTextPanelBox(habits, height), createUpperTextPanelBox(description, height+1))
	// view = lipgloss.JoinHorizontal(0, view, lipgloss.Place(10, 10, 0, 0, description, lipgloss.WithWhitespaceForeground(lipgloss.Color("0xffff"))))

	lowerPanel := lipgloss.JoinHorizontal(1,
		createLowerPanelTextBox(fmt.Sprintf("Strike (task %d):\n\tCurrent: 0\n\tBest monthly: 0\n\tLongest: 0", selectedID), 4),
		createLowerPanelTextBox(fmt.Sprintf("Completion (task %d):\n\tThis week: 0\n\tThis month: 0\n\tThis year: 0", selectedID), 4))

	view = lipgloss.JoinVertical(1, view, lowerPanel)

	// The footer
	view += "\n\nPress q to quit.\n"

	// Send the UI for rendering
	return view
}
