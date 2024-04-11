package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	defaultPort = 3000
	defaultHost = "localhost"
)

type Option func(*Config) error

type Config struct {
	host string
	port int
}

func DefaultConfig() Config {
	return Config{
		host: defaultHost,
		port: defaultPort,
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

	http.ListenAndServe(fmt.Sprintf("%s:%d", c.host, c.port), r)
	server := &http.Server{Addr: fmt.Sprintf("%s:%d", c.host, c.port), Handler: r}
	server.ListenAndServe()

	return server, nil
}
