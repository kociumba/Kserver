package main

import (
	"flag"
	"fmt"
	"kserver/handlers"
	"kserver/internal"
	luainternal "kserver/luaintegration"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	lua "github.com/yuin/gopher-lua"
)

var (
	port   = flag.Int("port", 0, "port to listen on")
	useLua = flag.Bool("lua", false, "run with lua configuration")
)

func main() {

	flag.Parse()

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

	cfg, err := handlers.GetHandlers(*useLua)
	if err != nil {
		panic(err)
	}

	if *useLua {
		l := lua.NewState(lua.Options{MinimizeStackMemory: true})
		defer l.Close()
		luainternal.RegisterRouteType(l)
		l.SetGlobal("registerRoutes", l.NewFunction(luainternal.LuaRegisterRoutes))

		// Create a temporary file to hold the Lua script
		tempFile := luainternal.CreateTempLUA()
		defer os.Remove(tempFile.Name())

		// Set the temporary file as the value of the kserver_lua_file global
		l.SetGlobal("kserver_lua_file", lua.LString(tempFile.Name()))
		if err := l.DoString("assert(loadfile(kserver_lua_file))()"); err != nil {
			panic(err)
		}
	}

	// if os.Args[1] != "-lua" || len(os.Args) < 2 {
	for _, route := range cfg.Handlers {
		fmt.Printf("Registering route: \033[34m%+v\033[0m from YAML\n", route.Route)
		go handlers.RegisterRoutes(route)
	}
	// }

	if *port != 0 {
		cfg.Port = *port
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

		// fmt.Println("\033[1m" + os.Interrupt.String() + "\033[0m")

		if !executed {
			if os.Interrupt == nil || os.Getpid() == 1 {
				fmt.Println("\033[31mCtrl+C pressed, force closing\033[0m")
				continue
			}
			fmt.Println("\033[31mUnknown command, run \033[1mhelp\033[0m")
		}
		executed = false
	}
}
