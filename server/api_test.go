package server

import (
	"fmt"
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
	if code := resp.StatusCode; code != 200 {
		t.Fatalf("Incorrect status code has '%d' but expected 200", code)
	}
}
