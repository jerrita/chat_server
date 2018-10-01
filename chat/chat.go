package main

import (
	"client"
	"fmt"
	"os"
	"server"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage: chat <server/client>")
		os.Exit(1)
	}
	switch args[1] {
	case "server":
		server.Run()
	case "client":
		client.Run()
	default:
		fmt.Println("Usage: chat <server/client>")
	}
}
