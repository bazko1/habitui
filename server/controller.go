package server

import (
	"errors"

	"github.com/bazko1/habitui/habit"
)

var (
	ErrUsernameExists  = errors.New("user with given username already exists")
	ErrInccorectInput  = errors.New("incorrect input provided")
	ErrEmailRegistered = errors.New("some user already registered with given email")
	ErrNonExistentUser = errors.New("user does not exists")
)

type Controller interface {
	CreateNewUser(user UserModel) (UserModel, error)
	UpdateUserHabits(user UserModel, habits habit.TaskList) error
	GetUserHabits(user UserModel) (habit.TaskList, error)
}
