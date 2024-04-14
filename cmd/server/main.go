//nolint:forbidigo //prints for client are not debug statements
package main

import (
	"fmt"

	"github.com/bazko1/habitui/server"
)

func main() {
	server, err := server.New()
	if err != nil {
		fmt.Printf("Failed to create new server: %v\n", err)

		return
	}

	fmt.Println("Server is listening at:", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Failed to listen and serve: %v", err)

		return
	}
}
