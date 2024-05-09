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
		Password: u.Password,
		Habits:   make(habit.TaskList, 0),
	}

	return controller.users[username], nil
}

func (controller *InMemoryController) UpdateUserHabits(user UserModel, habits habit.TaskList,
) error {
	if controller.IsValid(user) {
		return ErrNonExistentUserOrPassword
	}

	u := controller.users[user.Username]
	u.Habits = habits
	controller.users[user.Username] = u

	return nil
}

func (controller InMemoryController) GetUserHabits(user UserModel) (habit.TaskList, error) {
	if controller.IsValid(user) {
		return habit.TaskList{}, ErrNonExistentUserOrPassword
	}

	return controller.users[user.Username].Habits, nil
}

func (controller InMemoryController) IsValid(user UserModel) bool {
	u, exist := controller.users[user.Username]
	return !exist || u.Password != user.Password
}
