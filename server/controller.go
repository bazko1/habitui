package server

import (
	"errors"

	"github.com/bazko1/habitui/habit"
)

var (
	ErrUsernameAlreadyExists     = errors.New("user with given username already exists")
	ErrUsernameDoesNotExist      = errors.New("user with given username does not exist")
	ErrInccorectInput            = errors.New("incorrect input provided")
	ErrEmailRegistered           = errors.New("some user already registered with given email")
	ErrNonExistentUserOrPassword = errors.New("user with given name does not exist or incorrect password")
)

type Controller interface {
	Initialize() error
	GetUserByName(name string) (UserModel, error)
	CreateNewUser(user UserModel) (UserModel, error)
	UpdateUserHabits(user UserModel, habits habit.TaskList) error
	GetUserHabits(user UserModel) (habit.TaskList, error)
	IsValid(user UserModel) (bool, error)
	Finalize() error
}
