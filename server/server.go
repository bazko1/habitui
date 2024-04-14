package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	DefaultPort                      = 3000
	DefaultHost                      = "localhost"
	DefaultReadTimeout time.Duration = 100 * time.Millisecond
)

type Option func(*Config) error

func WithHost(host string) Option {
	return func(c *Config) error {
		c.host = host

		return nil
	}
}

func WithPort(port int) Option {
	return func(c *Config) error {
		c.port = port

		return nil
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(c *Config) error {
		c.readTimeout = timeout

		return nil
	}
}

type Config struct {
	host        string
	port        int
	readTimeout time.Duration
}

func DefaultConfig() Config {
	return Config{
		host:        DefaultHost,
		port:        DefaultPort,
		readTimeout: DefaultReadTimeout,
	}
}

func New(opts ...Option) (*http.Server, error) {
	c := DefaultConfig()
	h := createHandler()

	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, fmt.Errorf("option failed %w", err)
		}
	}

	server := &http.Server{
		Addr:        fmt.Sprintf("%s:%d", c.host, c.port),
		Handler:     h,
		ReadTimeout: c.readTimeout,
	}

	return server, nil
}

func createHandler() http.Handler {
	logRequestMiddleware := func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Printf("method: %s, path: %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		}
	}

	r := http.NewServeMux()

	r.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.HandleFunc("POST /user/create", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("creating user"))
	})

	r.HandleFunc("GET /user/habits", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("getting user habits"))
	})

	r.HandleFunc("PUT /user/habits", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("updating user habits"))
	})

	r.HandleFunc("PUT /user/token", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("revoking user token"))
	})

	return logRequestMiddleware(r)
}
