package main

import (
	"Kserver/handlers"
	"Kserver/internal"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	lua "github.com/yuin/gopher-lua"
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

	if len(os.Args) > 1 && os.Args[1] == "-lua" {
		l := lua.NewState(lua.Options{MinimizeStackMemory: true})
		defer l.Close()
		registerRouteType(l)
		l.SetGlobal("registerRoutes", l.NewFunction(luaRegisterRoutes))
		l.SetGlobal("kserver_lua_file", lua.LString("kserver.lua"))
		if err := l.DoString("assert(loadfile(kserver_lua_file))()"); err != nil {
			panic(err)
		}
	}

	cfg, err := handlers.GetHandlers()
	if err != nil {
		// panic(err)
		cfg.Port = 8000
	}

	// if os.Args[1] != "-lua" || len(os.Args) < 2 {
	for _, route := range cfg.Handlers {
		go handlers.RegisterRoutes(route)
	}
	// }

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

const luaRouteTypeName = "route"

func registerRouteType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaRouteTypeName)
	L.SetGlobal("route", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newRoute))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), routeMethods))
}

func newRoute(L *lua.LState) int {
	route := &handlers.Route{
		Route:       L.CheckString(1),
		Content:     L.CheckString(2),
		ContentType: L.CheckString(3),
	}
	ud := L.NewUserData()
	ud.Value = route
	L.SetMetatable(ud, L.GetTypeMetatable(luaRouteTypeName))
	L.Push(ud)
	return 1
}

func checkRoute(L *lua.LState) *handlers.Route {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*handlers.Route); ok {
		return v
	}
	L.ArgError(1, "route expected")
	return nil
}

var routeMethods = map[string]lua.LGFunction{
	"route":       routeGetSetRoute,
	"content":     routeGetSetContent,
	"contentType": routeGetSetContentType,
}

func routeGetSetRoute(L *lua.LState) int {
	r := checkRoute(L)
	if L.GetTop() == 2 {
		r.Route = L.CheckString(2)
		return 0
	}
	L.Push(lua.LString(r.Route))
	return 1
}

func routeGetSetContent(L *lua.LState) int {
	r := checkRoute(L)
	if L.GetTop() == 2 {
		r.Content = L.CheckString(2)
		return 0
	}
	L.Push(lua.LString(r.Content))
	return 1
}

func routeGetSetContentType(L *lua.LState) int {
	r := checkRoute(L)
	if L.GetTop() == 2 {
		r.ContentType = L.CheckString(2)
		return 0
	}
	L.Push(lua.LString(r.ContentType))
	return 1
}

func luaRegisterRoutes(L *lua.LState) int {
	route := checkRoute(L)
	fmt.Printf("luaRegisterRoutes called with route: %+v\n", route)
	handlers.RegisterRoutes(*route)
	return 0 // No failure condition, always returns success
}
