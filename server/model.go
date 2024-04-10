package server

import "github.com/bazko1/habitui/habit"

type UserModel struct {
	username  string
	email     string
	passwordd string
	habits    habit.TaskList
}
