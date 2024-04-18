package server

import (
	"fmt"

	"github.com/bazko1/habitui/habit"
)

type InMemoryController struct {
	users map[string]UserModel
}

func NewInMemoryController() InMemoryController {
	return InMemoryController{users: make(map[string]UserModel)}
}

func (controller *InMemoryController) CreateNewUser(u UserModel) (bool, error) {
	username := u.Username
	email := u.Email
	if username == "" {
		return false, fmt.Errorf("%w: username can not be empty", ErrInccorectInput)
	}

	if _, exists := controller.users[username]; exists {
		return false, ErrUsernameExists
	}

	controller.users[username] = UserModel{
		Username: username,
		Email:    email,
		// TODO: Generate token
		Token:  "token",
		habits: make(habit.TaskList, 0),
	}

	return true, nil
}

func (controller *InMemoryController) UpdateUserHabits(user UserModel, habits habit.TaskList,
) error {
	if u, exist := controller.users[user.Username]; !exist || u.Token != user.Token {
		return fmt.Errorf("%w: user with given name does not exist or incorrect token", ErrNonExistentUser)
	}

	u := controller.users[user.Username]
	u.habits = habits
	controller.users[user.Username] = u

	return nil
}

func (controller InMemoryController) GetUserHabits(user UserModel) (habit.TaskList, error) {
	if u, exist := controller.users[user.Username]; !exist || u.Token != user.Token {
		return habit.TaskList{}, fmt.Errorf("%w: user with given name does not exist or incorrect token", ErrNonExistentUser)
	}
	return controller.users[user.Username].habits, nil
}
