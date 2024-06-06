//nolint:forbidigo //prints for client are not debug statements
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bazko1/habitui/server"
)

func main() {
	host := flag.String("hostname", server.DefaultHost, "host name or ip to serve on")
	port := flag.Int("port", server.DefaultPort, "port to serve on")
	timeout := flag.Int64("timeout", server.DefaultReadTimeoutMiliseconds.Milliseconds(), "read timeout milliseconds")
	controllerEngine := flag.String("engine", server.DefaultControllerEngine,
		"engine to use for controller supported: 'inmem', 'sqlite'")

	flag.Parse()

	retCode := 0
	defer func() { os.Exit(retCode) }()

	server, finalizefn, err := server.New(
		server.WithHost(*host),
		server.WithPort(*port),
		server.WithReadTimeout(time.Duration(*timeout)*time.Millisecond),
		server.WithControllerEngine(*controllerEngine),
	)
	if err != nil {
		fmt.Printf("Failed to create new server: %v\n", err)

		return
	}

	defer func() {
		if err := finalizefn(); err != nil {
			fmt.Printf("Failed to finalize server: %v\n", err)
		}
	}()

	fmt.Println("Server is listening at:", server.Addr)
	fmt.Println("Using controller engine:", *controllerEngine)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	serverServeError := make(chan error, 1)
	defer close(serverServeError)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			serverServeError <- err

			return
		}
	}()

	for {
		select {
		case err := <-serverServeError:
			fmt.Printf("Failed to listen and serve: %v\n", err)

			retCode = 1

			return
		case <-sigs:
			const shutdownTimeout = 10 * time.Second

			ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				fmt.Printf("Failed to shutdown server gracefully: %v\n", err)
				server.Close()
			}

			return
		default:
			break
		}
	}
}
