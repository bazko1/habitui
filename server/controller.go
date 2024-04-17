package server

import (
	"errors"

	"github.com/bazko1/habitui/habit"
)

var (
	ErrUsernameExists  = errors.New("User with given username already exists")
	ErrInccorectInput  = errors.New("Incorrect input provided")
	ErrEmailRegistered = errors.New("Some user already registered with given email")
)

type Controller interface {
	CreateNewUser(username, email, password string) (bool, error)
	UpdateUserHabits(username, password string, habits habit.TaskList)
	GetUserHabits(username, password string) habit.TaskList
}
