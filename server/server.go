package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, fmt.Errorf("option failed %w", err)
		}
	}

	server := &http.Server{
		Addr:        fmt.Sprintf("%s:%d", c.host, c.port),
		Handler:     r,
		ReadTimeout: c.readTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		return nil, fmt.Errorf("failed to listen and serve: %w", err)
	}

	return server, nil
}
