package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	DefaultPort                                 = 3000
	DefaultHost                                 = "localhost"
	DefaultReadTimeoutMiliseconds time.Duration = 100 * time.Millisecond
	DefaultControllerEngine                     = "inmem"
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

// WithReadTimeout sets server read timeout in  milliseconds.
func WithReadTimeout(timeout time.Duration) Option {
	return func(c *Config) error {
		c.readTimeout = timeout

		return nil
	}
}

// WithControllerEngine sets controller engine name that should be used.
func WithControllerEngine(engineName string) Option {
	return func(c *Config) error {
		c.controllerEngine = engineName
		// TODO:
		// let me think if I would like to also set actual object or pointer here
		// and validate if name is one of {'sqlite', 'imem'}
		// alternatively this loigc will stay in New function..

		return nil
	}
}

// WitSqlitePath sets SQLiteController data source if controller engine is sqlite.
func WitSqliteDataSource(dataSource string) Option {
	return func(c *Config) error {
		c.sqliteDataSource = dataSource

		return nil
	}
}

type Config struct {
	host             string
	port             int
	readTimeout      time.Duration
	controllerEngine string
	sqliteDataSource string
}

func DefaultConfig() Config {
	return Config{
		host:             DefaultHost,
		port:             DefaultPort,
		readTimeout:      DefaultReadTimeoutMiliseconds,
		controllerEngine: DefaultControllerEngine,
	}
}

func New(opts ...Option) (*http.Server, error) {
	c := DefaultConfig()

	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, fmt.Errorf("option failed %w", err)
		}
	}

	var controller Controller

	switch c.controllerEngine {
	case "inmem":
		controller = NewInMemoryController()
	case "sqlite":
		source := sqliteDatasePathEnvName
		if c.sqliteDataSource != "" {
			source = c.sqliteDataSource
		}

		controller = NewSQLiteController(source)
	default:
		return nil, errors.New("wrong controller engine provided")
	}

	if err := controller.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize controller: %w", err)
	}

	h := createHandler(controller)

	server := &http.Server{
		Addr:        fmt.Sprintf("%s:%d", c.host, c.port),
		Handler:     h,
		ReadTimeout: c.readTimeout,
	}

	return server, nil
}
