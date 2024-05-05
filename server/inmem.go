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

func (controller *InMemoryController) CreateNewUser(u UserModel) (UserModel, error) {
	username := u.Username
	email := u.Email

	if username == "" {
		return UserModel{}, fmt.Errorf("%w: username can not be empty", ErrInccorectInput)
	}

	if _, exists := controller.users[username]; exists {
		return UserModel{}, ErrUsernameExists
	}

	controller.users[username] = UserModel{
		Username: username,
		Email:    email,
		// TODO: Generate token
		Token:  "token",
		Habits: make(habit.TaskList, 0),
	}

	return controller.users[username], nil
}

func (controller *InMemoryController) UpdateUserHabits(user UserModel, habits habit.TaskList,
) error {
	if u, exist := controller.users[user.Username]; !exist || u.Token != user.Token {
		return fmt.Errorf("%w: user with given name does not exist or incorrect token", ErrNonExistentUser)
	}

	u := controller.users[user.Username]
	u.Habits = habits
	controller.users[user.Username] = u

	return nil
}

func (controller InMemoryController) GetUserHabits(user UserModel) (habit.TaskList, error) {
	if u, exist := controller.users[user.Username]; !exist || u.Token != user.Token {
		return habit.TaskList{}, fmt.Errorf("%w: user with given name does not exist or incorrect token", ErrNonExistentUser)
	}

	return controller.users[user.Username].Habits, nil
}
