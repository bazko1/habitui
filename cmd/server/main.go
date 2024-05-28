//nolint:forbidigo //prints for client are not debug statements
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/bazko1/habitui/server"
)

func main() {
	host := flag.String("hostname", server.DefaultHost, "host name or ip to serve on")
	port := flag.Int("port", server.DefaultPort, "port to serve on")
	timeout := flag.Int64("timeout", server.DefaultReadTimeoutMiliseconds.Milliseconds(), "read timeout milliseconds")
	controllerEngine := flag.String("engine", "inmem", "engine to use for controller")
	flag.Parse()

	server, err := server.New(
		server.WithHost(*host),
		server.WithPort(*port),
		server.WithReadTimeout(time.Duration(*timeout)*time.Millisecond),
	)
	if err != nil {
		fmt.Printf("Failed to create new server: %v\n", err)

		return
	}

	fmt.Println("Server is listening at:", server.Addr)
	// TODO: Check if engine is one of "inmem", "sqlite"
	fmt.Println("Using controller engine", *controllerEngine)

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Failed to listen and serve: %v", err)

		return
	}
}
