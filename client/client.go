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

// LoadTasksOrCreateUser check if user with given username exists
// if it does tasks belonging to user will be returned otherwise
// user is created and empty task list is returned.
func (client HTTPClient) LoadTasksOrCreateUser() (habit.TaskList, error) {
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
		return habit.TaskList{}, fmt.Errorf("failed to login error: %w", err)
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		return habit.TaskList{}, fmt.Errorf("failed to login incorrect status code %d", code)
	}

	tokenData := map[string]any{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		return habit.TaskList{}, fmt.Errorf("failed to decode token data to login error: %w", err)
	}

	token, _ := tokenData["access_token"].(string)

	// get habits
	ctx, cancel = context.WithTimeout(context.Background(), httpTimeout*time.Second)
	defer cancel()

	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, client.Address+"/user/habits", http.NoBody)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return habit.TaskList{}, fmt.Errorf("failed to get user habits error: %w", err)
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		return habit.TaskList{}, fmt.Errorf("failed to get user habits incorrect code: %d", code)
	}

	habits := habit.TaskList{}
	if err := json.NewDecoder(resp.Body).Decode(&habits); err != nil {
		return habit.TaskList{}, fmt.Errorf("failed to decode habits tasks %w", err)
	}

	return habits, nil
}
