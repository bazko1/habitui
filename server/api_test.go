package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
)

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

func TestCreateUser(t *testing.T) {
	t.Parallel()

	ln := startServer(t)
	defer ln.Close()
	address := fmt.Sprintf("http://%s", ln.Addr().String())
	fmt.Println("Server listening at: ", address)
	resp, err := http.Post(address+"/user/create", "application/x-www-form-urlencoded", strings.NewReader(`{"Username":"foo","Email":"bar"}`))
	if err != nil {
		t.Fatalf("Error during post call user creation: %v", err)
	}
	if code := resp.StatusCode; code != 201 {
		t.Fatalf("Incorrect status code has '%d' but expected 201", code)
	}

	user := UserModel{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatalf("Error decoding user body: %v", err)
	}
	if user.Token == "" {
		t.Fatal("User has no token set")
	}

	resp, err = http.Post(address+"/user/create", "application/x-www-form-urlencoded", strings.NewReader(`{"Username":"foo","Email":"bar"}`))
	if err != nil {
		t.Fatalf("Error during post call user recreation: %v", err)
	}

	if code := resp.StatusCode; code != http.StatusNoContent {
		t.Fatalf("Incorrect status code has '%d' but expected %d", code, http.StatusNoContent)
	}

	bytes, _ := io.ReadAll(resp.Body)
	if data := string(bytes); data != "" {
		t.Fatalf("Expected empty body got %s", data)
	}
}
