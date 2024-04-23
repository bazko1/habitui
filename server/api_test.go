package server

import (
	"fmt"
	"net"
	"net/http"
	"testing"
)

func startServer(t *testing.T) (net.Listener, chan error) {
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

	return ln, serverServeError
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	ln, serverErrors := startServer(t)
	select {
	case err := <-serverErrors:
		t.Fatalf("Failed to start server error %v:", err)
	default:
		fmt.Println("Server listening at: ", ln.Addr())
		resp, err := http.Get(fmt.Sprintf("http://%s/user/habits", ln.Addr().String()))
		fmt.Println(resp, err)
	}
}
