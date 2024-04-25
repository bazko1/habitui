package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func renderCalendar() {
	now := time.Now()
	days := getDaysInMonth(time.Now())
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	firstDayWeekday := int(firstDay.Weekday())
	daysSlice := make([]string, 0, days)

	for i := firstDayWeekday; i > 0; i-- {
		daysSlice = append(daysSlice, fmt.Sprintf("%d", firstDay.AddDate(0, 0, -i).Day()))
	}
	for i := 1; i <= days; i++ {
		daysSlice = append(daysSlice, fmt.Sprintf("%d", i))
	}
	// numWeeks := days/7 + 1
	// weeksLabels :=
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
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(0, 1)
		})

	// files := labelStyle.Render(strings.Join([]string{" 1", "2", "3", "4", "5", "6", "7", "8"}, "\n\n "))

	dayNames := labelStyle.Render(strings.Join([]string{" Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}, "  "))
	fmt.Println(lipgloss.JoinVertical(lipgloss.Left, dayNames, t.Render()) + "\n")
}
