package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bazko1/habitui/habit"
)

const httpTimeout = 5

type HTTPClient struct {
	Address  string
	Username string
	Password string
}

func (client HTTPClient) login() (string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout*time.Second)
	defer cancel()

	userData := fmt.Sprintf(`{"Username":"%s","Password":"%s"}`, client.Username, client.Password)
	req, _ := http.NewRequestWithContext(ctx,
		http.MethodPost,
		client.Address+"/user/login",
		strings.NewReader(userData))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// TODO: Will need to handle user creation differently. Probably I need
		// special error from server that directly points that user does not exist.
		return "", 0, fmt.Errorf("failed to login: %w", err)
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		return "", resp.StatusCode, fmt.Errorf("failed to login incorrect status code %d", code)
	}

	tokenData := map[string]any{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		return "", 0, fmt.Errorf("failed to decode token data to login: %w", err)
	}

	token, _ := tokenData["access_token"].(string)

	return token, resp.StatusCode, nil
}

func (client HTTPClient) createUser() error {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout*time.Second)
	defer cancel()

	userData := fmt.Sprintf(`{"Username":"%s","Password":"%s"}`, client.Username, client.Password)
	req, _ := http.NewRequestWithContext(ctx,
		http.MethodPost,
		client.Address+"/user/create",
		strings.NewReader(userData))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("during post call user creation: %w", err)
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code == http.StatusNoContent {
		return fmt.Errorf("username %q is already taken wrong password provided", client.Username)
	}

	if code := resp.StatusCode; code != http.StatusCreated {
		return fmt.Errorf("incorrect status code for user create has '%d' but expected 201", code)
	}

	return nil
}

// LoadTasksOrCreateUser check if user with given username exists
// if it does tasks belonging to user will be returned otherwise
// user is created and empty task list is returned.
func (client HTTPClient) LoadTasksOrCreateUser() (habit.TaskList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout*time.Second)
	defer cancel()

	// try to login or create user on login failure
	// if also creation fails return error
	token, statusCode, err := client.login()
	if err != nil {
		if statusCode == http.StatusUnauthorized {
			if err := client.createUser(); err != nil {
				return habit.TaskList{}, fmt.Errorf("failed to create new user: %w", err)
			}
			token, _, err = client.login()
		}

		if err != nil {
			return habit.TaskList{}, fmt.Errorf("failed to login user to load tasks: %w", err)
		}
	}

	// get habits
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, client.Address+"/user/habits", http.NoBody)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return habit.TaskList{}, fmt.Errorf("failed to get user habits: %w", err)
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		return habit.TaskList{}, fmt.Errorf("failed to get user habits incorrect code: %d", code)
	}

	habits := habit.TaskList{}
	if err := json.NewDecoder(resp.Body).Decode(&habits); err != nil {
		return habit.TaskList{}, fmt.Errorf("failed to decode habits tasks: %w", err)
	}

	return habits, nil
}

// SaveUserTasks saves client habits to remote server.
func (client HTTPClient) SaveUserTasks(habits habit.TaskList) error {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout*time.Second)
	defer cancel()

	token, _, err := client.login()
	if err != nil {
		return fmt.Errorf("failed to login user to load tasks: %w", err)
	}

	b, _ := json.Marshal(habits)
	userJSON := string(b)
	newHabits := strings.NewReader(userJSON)

	req, _ := http.NewRequestWithContext(ctx,
		http.MethodPut,
		client.Address+"/user/habits",
		newHabits)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error during post call user creation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("put /user/habits should return %d while it returned %d", http.StatusOK, resp.StatusCode)
	}

	return nil
}
