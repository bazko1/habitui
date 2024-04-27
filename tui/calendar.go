package tui

import (
	"fmt"
	"os"
	"slices"
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
func RenderCalendar(task habit.Task) string {
	now := task.GetTime()
	days := getDaysInMonth(time.Now())
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	firstDayWeekday := int(firstDay.Weekday())
	daysSlice := make([]string, 0, days)

	completedDays := [][2]int{}
	for i := firstDayWeekday; i > 0; i-- {
		prvMonthDay := firstDay.AddDate(0, 0, -i)
		daysSlice = append(daysSlice, fmt.Sprintf("%d", prvMonthDay.Day()))
		if task.WasCompletedAt(prvMonthDay.Year(), prvMonthDay.Month(), prvMonthDay.Day()) {
			completedDays = append(completedDays, [2]int{0, firstDayWeekday - i})
		}
	}

	for i := 1; i <= days; i++ {
		daysSlice = append(daysSlice, strconv.Itoa(i))
		row := (i - 1) / 7
		col := (i - firstDayWeekday) % 7
		if task.WasCompletedAt(now.Year(), now.Month(), i) {
			completedDays = append(completedDays, [2]int{row, col})
		}
	}

	weeks := [][]string{}

	var i int
	for ; i < days/7; i++ {
		weeks = append(weeks, daysSlice[i*7:i*7+7])
	}

	weeks = append(weeks, daysSlice[i*7:])

	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	labelStyle := re.NewStyle().Foreground(lipgloss.Color("241"))
	selectedStyle := baseStyle.Copy().Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderRow(true).
		BorderColumn(true).
		Rows(weeks...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if slices.ContainsFunc(completedDays, func(a [2]int) bool { return a[0] == row && a[1] == col }) {
				return selectedStyle
			}

			return baseStyle
		})

	dayNames := labelStyle.Render(strings.Join([]string{" Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}, "  "))
	calendar := lipgloss.JoinVertical(lipgloss.Left, dayNames, t.Render()) + "\n"

	return calendar
}
