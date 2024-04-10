package server

import "github.com/bazko1/habitui/habit"

type Controller interface {
	CreateNewUser(username, email, password string)
	UpdateUserHabits(username, password string, habits habit.TaskList)
	GetUserHabits(username, password string) habit.TaskList
}
