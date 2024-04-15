package server

import (
	"github.com/bazko1/habitui/habit"
)

type UserModel struct {
	Username  string
	Email     string
	Passwordd string `json:"-"`
	habits    habit.TaskList
}
