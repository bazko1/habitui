package server

import (
	"github.com/bazko1/habitui/habit"
)

type InMemoryController struct {
	users map[string]UserModel
}

func NewInMemoryController() InMemoryController {
	return InMemoryController{users: make(map[string]UserModel)}
}

func (controller *InMemoryController) CreateNewuser(username,
	email,
	password string,
) (bool, error) {
	if _, exists := controller.users[username]; exists {
		return false, ErrUsernameExists
	}

	controller.users[username] = UserModel{
		Username:  username,
		Email:     email,
		Passwordd: password,
		habits:    make(habit.TaskList, 0),
	}

	return true, nil
}

func (controller *InMemoryController) UpdateUserHabits(username string,
	habits habit.TaskList,
) {
	if user, exist := controller.users[username]; exist {
		user.habits = habits
	}
}

func (controller InMemoryController) GetUserHabits(username string) habit.TaskList {
	return controller.users[username].habits
}
