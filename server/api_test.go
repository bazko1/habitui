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

const testUserString = `{"Username":"foo","Email":"bar","Password":"test"}`

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

func createUser(t *testing.T, addr string) {
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
		t.Fatalf("Incorrect status code for user create has '%d' but expected 201", code)
	}
}

func loginUser(t *testing.T, addr string) map[string]any {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx,
		http.MethodPost,
		addr+"/user/login",
		strings.NewReader(testUserString))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error during post call user login: %v", err)
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != 200 {
		t.Fatalf("Incorrect status code for login has '%d' but expected 201", code)
	}

	tokenData := map[string]any{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		t.Fatalf("Error decoding login token body: %v", err)
	}

	return tokenData
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	ln := startServer(t)
	defer ln.Close()
	address := "http://" + ln.Addr().String()
	fmt.Println("Server listening at: ", address)
	createUser(t, address)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check if recreating user returns status no content
	// with empty body
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

	getHabits := func(token string) *http.Response {
		url := address + "/user/habits"

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
		req.Header.Add("Authorization", token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Error failed to get user habits : %v", err)
		}

		return resp
	}

	createUser(t, address)
	tokenData := loginUser(t, address)
	token, _ := tokenData["access_token"].(string)

	resp := getHabits(fmt.Sprintf("Bearer %s", token))
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	if string(bytes) != "[]" {
		t.Fatalf("Error initial user habits should be empty while it is %s", string(bytes))
	}

	startDate := time.Date(2024, 3, 30, 12, 0, 0, 0, time.Local)
	now := func() time.Time {
		return startDate
	}
	task := habit.WithCustomTime("work on habitui", "daily app grind", now)
	completions := 5

	for range completions {
		task.MakeCompleted()

		startDate = startDate.AddDate(0, 0, 1)
	}

	user := UserModel{}
	_ = json.Unmarshal([]byte(testUserString), &user)
	user.Habits = habit.TaskList{task}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	b, _ := json.Marshal(user.Habits)
	userJSON := string(b)
	newHabits := strings.NewReader(userJSON)

	token, _ = loginUser(t, address)["access_token"].(string)
	req, _ := http.NewRequestWithContext(ctx,
		http.MethodPut,
		address+"/user/habits",
		newHabits)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error during post call user creation: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Put /user/habits should return %d while it returned %d", http.StatusOK, resp.StatusCode)
	}

	token, _ = loginUser(t, address)["access_token"].(string)

	resp = getHabits("Bearer " + token)
	defer resp.Body.Close()

	if code := resp.StatusCode; code != 200 {
		t.Fatalf("Incorrect status code has '%d' but expected 200", code)
	}

	updatedHabits := habit.TaskList{}
	if err := json.NewDecoder(resp.Body).Decode(&updatedHabits); err != nil {
		t.Fatalf("Error decoding user body: %v", err)
	}

	// TODO: Lets for now compare jsons marshal
	// I might want to write custom comparisor
	// that looks at fields in detail.
	bytesOld, _ := json.MarshalIndent(user.Habits, "", "  ")
	bytesNew, _ := json.MarshalIndent(updatedHabits, "", "  ")
	if oldH, newH := string(bytesOld), string(bytesNew); !reflect.DeepEqual(newH, oldH) {
		t.Fatalf(`Updated tasks returned by server differ from template:
			old:= %v
		  new:= %v`, oldH, newH)
	}
}
