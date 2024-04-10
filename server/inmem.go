package server

import "github.com/bazko1/habitui/habit"

type InMemoryController struct {
	users map[string]UserModel
}

func NewInMemoryController() InMemoryController {
	return InMemoryController{users: make(map[string]UserModel)}
}

func (controller *InMemoryController) CreateNewuser(username,
	email,
	password string,
) {
	controller.users[username] = UserModel{
		username:  username,
		email:     email,
		passwordd: password,
		habits:    make(habit.TaskList, 0),
	}
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
