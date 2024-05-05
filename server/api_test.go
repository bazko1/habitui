// // nolint:forbidigo
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bazko1/habitui/habit"
)

const testUserString = `{"Username":"foo","Email":"bar"}`

func startServer(t *testing.T) net.Listener {
	t.Helper()

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create new listener error: %v", err)
	}

	controller := NewInMemoryController()
	serverServeError := make(chan error)

	go func() {
		defer close(serverServeError)
		if err := http.Serve(ln, createHandler(&controller)); err != nil {
			serverServeError <- err
		}
	}()

	// one time check if errors occurred or continue
	select {
	case err := <-serverServeError:
		t.Fatalf("Failed to start server error %v:", err)
	default:
		break
	}

	return ln
}

func createUser(t *testing.T, addr string) UserModel {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx,
		http.MethodPost,
		addr+"/user/create",
		strings.NewReader(testUserString))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error during post call user creation: %v", err)
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != 201 {
		t.Fatalf("Incorrect status code has '%d' but expected 201", code)
	}

	user := UserModel{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatalf("Error decoding user body: %v", err)
	}

	return user
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	ln := startServer(t)
	defer ln.Close()
	address := "http://" + ln.Addr().String()
	fmt.Println("Server listening at: ", address)
	user := createUser(t, address)

	if user.Token == "" {
		t.Fatal("User has no token set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx,
		http.MethodPost,
		address+"/user/create",
		strings.NewReader(testUserString))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error during post call user recreation: %v", err)
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusNoContent {
		t.Fatalf("Incorrect status code has '%d' but expected %d", code, http.StatusNoContent)
	}

	bytes, _ := io.ReadAll(resp.Body)
	if data := string(bytes); data != "" {
		t.Fatalf("Expected empty body got %s", data)
	}
}

func TestUpdateUserTasks(t *testing.T) {
	t.Parallel()
	ln := startServer(t)

	defer ln.Close()
	address := "http://" + ln.Addr().String()
	fmt.Println("Server listening at: ", address)

	getHabits := func(body io.Reader) *http.Response {
		url := address + "/user/habits"

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, body)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Error failed to get user habits : %v", err)
		}

		return resp
	}

	user := createUser(t, address)
	b, _ := json.Marshal(user)
	userJSON := string(b)
	baseBody := strings.NewReader(userJSON)

	resp := getHabits(baseBody)
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	if string(bytes) != "[]" {
		t.Fatal("Error initial user habits should be empty")
	}

	startDate := time.Date(2024, 3, 30, 12, 0, 0, 0, time.Local)
	now := func() time.Time {
		return startDate
	}
	task := habit.WithCustomTime("work on habittui", "daily app grind", now)
	completions := 5

	for range completions {
		task.MakeCompleted()

		startDate = startDate.AddDate(0, 0, 1)
	}

	user.Habits = habit.TaskList{task}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	b, _ = json.Marshal(user)
	userJSON = string(b)
	bodyWithHabits := strings.NewReader(userJSON)

	req, _ := http.NewRequestWithContext(ctx,
		http.MethodPut,
		address+"/user/habits",
		bodyWithHabits)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error during post call user creation: %v", err)
	}
	defer resp.Body.Close()

	resp = getHabits(strings.NewReader(userJSON))
	defer resp.Body.Close()

	if code := resp.StatusCode; code != 200 {
		t.Fatalf("Incorrect status code has '%d' but expected 200", code)
	}

	updatedHabits := habit.TaskList{}
	if err := json.NewDecoder(resp.Body).Decode(&updatedHabits); err != nil {
		t.Fatalf("Error decoding user body: %v", err)
	}

	fmt.Println("updated:=", updatedHabits)
	fmt.Println(reflect.DeepEqual(updatedHabits, user.Habits))
	// fmt.Println(updatedHabits[0].AllCompletion())
}
