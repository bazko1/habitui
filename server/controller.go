package server

import (
	"errors"

	"github.com/bazko1/habitui/habit"
)

var (
	ErrUsernameExists  = errors.New("User with given username already exists")
	ErrInccorectInput  = errors.New("Incorrect input provided")
	ErrEmailRegistered = errors.New("Some user already registered with given email")
	ErrNonExistentUser = errors.New("User does not exists")
)

type Controller interface {
	CreateNewUser(user UserModel) (bool, error)
	UpdateUserHabits(user UserModel, habits habit.TaskList) error
	GetUserHabits(user UserModel) (habit.TaskList, error)
}
