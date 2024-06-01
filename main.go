package main

import (
	"Kserver/handlers"
	"Kserver/internal"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func main() {
	// I have to do it like this couse otherwise HelpMsg can't initialize before using it
	internal.Commands["help"] = internal.Command{
		Name:        "help",
		NameANSI:    "\033[33mhelp\033[0m",
		Description: "show this message",
		Callback: func() {
			fmt.Println(internal.HelpMsg)
		},
	}

	var logFile, err = os.Create("KserverLOG.txt")

	log.SetOutput(logFile)

	var srv = &http.Server{
		Handler:           nil,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ErrorLog:          log.Default().StandardLog(),
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
			fmt.Println("Starting with TLS")
			srv.Addr = port
			err := srv.ListenAndServeTLS(cfg.Cert, cfg.Key)
			if err != nil {
				panic(err)
			}
		}()
	} else {
		fmt.Println("TLS not provided starting locally")
		go func() {
			srv.Addr = port
			err := srv.ListenAndServe()
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
