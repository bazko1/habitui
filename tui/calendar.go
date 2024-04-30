package tui

import (
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bazko1/habitui/habit"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	weekDays       = 7
	calendarRows   = 5
	calendarCols   = weekDays
	calendarFields = calendarRows * calendarCols
)

// RenderCalendar creates calendar string based on task that
// shows days in a month when task was completed.
func RenderCalendar(task habit.Task) string { //nolint:funlen // lets keep it as long blob for now
	now := task.GetTime()
	monthDays := getDaysInMonth(time.Now())
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	firstDayWeekday := int(firstDay.Weekday())

	// count Monday as first day instead of Sunday
	if firstDayWeekday == 0 {
		firstDayWeekday = 7
	}
	// Mon is 1, Tue 2, ..., Sun 7 if we have Monday
	// no days are filled before 1st day of month.
	firstDayWeekday--

	daysSlice := make([]string, 0, calendarFields)

	completedDays := [][2]int{}

	for i := firstDayWeekday; i > 0; i-- {
		prvMonthDay := firstDay.AddDate(0, 0, -i)
		daysSlice = append(daysSlice, strconv.Itoa(prvMonthDay.Day()))

		if task.WasCompletedAt(prvMonthDay.Year(), prvMonthDay.Month(), prvMonthDay.Day()) {
			completedDays = append(completedDays, [2]int{0, firstDayWeekday - i})
		}
	}

	for i := 1; i <= monthDays; i++ {
		daysSlice = append(daysSlice, strconv.Itoa(i))
		calendarDay := firstDayWeekday + i
		row := (calendarDay - 1) / weekDays
		col := (calendarDay - 1) % weekDays

		if task.WasCompletedAt(now.Year(), now.Month(), i) {
			completedDays = append(completedDays, [2]int{row, col})
		}
	}

	if leftDays := calendarFields - len(daysSlice); leftDays > 0 {
		for i := range leftDays {
			daysSlice = append(daysSlice, strconv.Itoa(i+1))
		}
	}

	weeks := [][]string{}

	var i int
	for ; i < monthDays/weekDays; i++ {
		weeks = append(weeks, daysSlice[i*weekDays:i*weekDays+weekDays])
	}

	weeks = append(weeks, daysSlice[i*weekDays:])

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
			if slices.ContainsFunc(completedDays, func(a [2]int) bool { return a[0] == row-1 && a[1] == col }) {
				return selectedStyle
			}

			return baseStyle
		})

	dayNames := labelStyle.Render(strings.Join([]string{" Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}, "  "))
	calendar := lipgloss.JoinVertical(lipgloss.Left, dayNames, t.Render()) + "\n"

	return calendar
}
