package main

import (
	"Kserver/handlers"
	"fmt"
	"net/http"
	"os"
)

func main() {
	cfg, err := handlers.GetHandlers()
	if err != nil {
		panic(err)
	}

	for _, route := range cfg.Handlers {
		go handlers.RegisterRoutes(route)
	}

	go http.ListenAndServe(":8080", nil)

	runREPL()
}

func runREPL() {
	var input string

	for {
		fmt.Print("\033[5m> \033[0m")
		fmt.Scan(&input)

		switch input {
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("unimplemented")
		}
	}
}
