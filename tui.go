package habitui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ListModel struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

func TuiModel() ListModel {
	return ListModel{
		// Our to-do list is a grocery list
		choices: []string{"work on habitui", "go for a walk", "app english lesson"},

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

func (m ListModel) View() string {
	// The header
	view := "Complete habit:\n\n"

	// Iterate over our choices
	for chID, choice := range m.choices {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == chID {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[chID]; ok {
			checked = "x" // selected!
		}

		// Render the row
		view += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	view += "\nPress q to quit.\n"

	// Send the UI for rendering
	return view
}
