package tui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bazko1/habitui/habit"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// renderCalendar
// TODO: this will render calendar that will be added as separate panel to tui
// calendar will be different for each task as completed days will be colored in green.
func renderCalendar(task habit.Task) string {
	now := task.GetTime()
	days := getDaysInMonth(time.Now())
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	firstDayWeekday := int(firstDay.Weekday())
	daysSlice := make([]string, 0, days)

	for i := firstDayWeekday; i > 0; i-- {
		daysSlice = append(daysSlice, fmt.Sprintf("%d", firstDay.AddDate(0, 0, -i).Day()))
	}

	for i := 1; i <= days; i++ {
		daysSlice = append(daysSlice, strconv.Itoa(i))
	}

	re := lipgloss.NewRenderer(os.Stdout)
	labelStyle := re.NewStyle().Foreground(lipgloss.Color("241"))

	weeks := [][]string{}

	var i int
	for ; i < days/7; i++ {
		weeks = append(weeks, daysSlice[i*7:i*7+7])
	}

	weeks = append(weeks, daysSlice[i*7:])

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderRow(true).
		BorderColumn(true).
		Rows(weeks...).
		StyleFunc(func(_, _ int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(0, 1)
		})

	dayNames := labelStyle.Render(strings.Join([]string{" Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}, "  "))
	calendar := lipgloss.JoinVertical(lipgloss.Left, dayNames, t.Render()) + "\n"

	return calendar
}
