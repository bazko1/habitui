package client

import "github.com/bazko1/habitui/habit"

type HTTPClient struct {
	Address  string
	Username string
	Password string
}

// LoadTasksOrCreateUser check if user with given username exists
// if it does tasks belonging to user will be returned otherwise
// user is created and empty task list is returned.
func (client HTTPClient) LoadTasksOrCreateUser() (habit.TaskList, error) {
	return habit.TaskList{}, nil
}
