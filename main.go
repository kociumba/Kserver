package main

import (
	"Kserver/handlers"
	"Kserver/internal"
	"fmt"
	"net/http"
)

func main() {
	// I have to do it like this couse otherwise HelpMsg can't initialize
	internal.Commands["help"] = internal.Command{
		Name:        "help",
		NameANSI:    "\033[33mhelp\033[0m",
		Description: "show this message",
		Callback: func() {
			fmt.Println(internal.HelpMsg)
		},
	}

	cfg, err := handlers.GetHandlers()
	if err != nil {
		panic(err)
	}

	for _, route := range cfg.Handlers {
		go handlers.RegisterRoutes(route)
	}

	port := ":" + fmt.Sprint(cfg.Port)
	fmt.Println("Server starting on port " + port)

	if cfg.Cert != "" && cfg.Key != "" {
		// do it with an anonymous func to still catch the error if any
		go func() {
			err := http.ListenAndServeTLS(port, cfg.Cert, cfg.Key, nil)
			if err != nil {
				panic(err)
			}
		}()
	} else {
		go func() {
			err := http.ListenAndServe(port, nil)
			if err != nil {
				panic(err)
			}
		}()
	}

	runREPL()
}

func runREPL() {
	var input string
	var executed bool

	// fmt.Println(internal.Commands)

	for {
		fmt.Print("\033[5m\033[1m> \033[0m")
		fmt.Scanln(&input)

		for _, command := range internal.Commands {
			if command.Name == input {
				if command.Callback != nil {
					command.Callback()
				}
				executed = true
				break
			}
		}

		if !executed {
			fmt.Println("\033[31mUnknown command, run \033[1mhelp\033[0m")
		}
		executed = false
	}
}
