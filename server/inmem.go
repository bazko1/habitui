package server

import (
	"fmt"

	"github.com/bazko1/habitui/habit"
)

type InMemoryController struct {
	users map[string]UserModel
}

func NewInMemoryController() *InMemoryController {
	return &InMemoryController{users: make(map[string]UserModel)}
}

func (InMemoryController) Initialize() error {
	return nil
}

func (controller *InMemoryController) CreateNewUser(u UserModel) (UserModel, error) {
	username := u.Username
	email := u.Email

	if username == "" {
		return UserModel{}, fmt.Errorf("%w: username can not be empty", ErrInccorectInput)
	}

	if u.Password == "" {
		return UserModel{}, fmt.Errorf("%w: password can not be empty", ErrInccorectInput)
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
	u := controller.users[user.Username]
	u.Habits = habits

	controller.users[user.Username] = u

	return nil
}

func (controller InMemoryController) GetUserHabits(user UserModel) (habit.TaskList, error) {
	if u, exist := controller.users[user.Username]; exist {
		return u.Habits, nil
	}

	return habit.TaskList{}, ErrNonExistentUserOrPassword
}

func (controller InMemoryController) IsValid(user UserModel) bool {
	u, exist := controller.users[user.Username]

	return exist && u.Password == user.Password
}

func (controller InMemoryController) GetUserByName(name string) (UserModel, bool) {
	u, exist := controller.users[name]
	if !exist {
		return UserModel{}, false
	}

	return u, true
}
