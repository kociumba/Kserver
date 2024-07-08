package luainternal

import (
	"fmt"

	"github.com/kociumba/kserver/handlers"

	lua "github.com/yuin/gopher-lua"
)

const luaRouteTypeName = "route"

func RegisterRouteType(L *lua.LState) {
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

	if route.ContentType == "application/json" {
		if L.GetTop() == 4 {
			route.LuaFunc = L.CheckFunction(4)
		}
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
		if r.ContentType == "application/json" && L.GetTop() == 3 {
			r.LuaFunc = L.CheckFunction(3)
		}
		return 0
	}
	L.Push(lua.LString(r.ContentType))
	return 1
}

func LuaRegisterRoutes(L *lua.LState) int {
	route := checkRoute(L)
	if route.ContentType == "application/json" && route.LuaFunc != nil {
		fmt.Printf("Registering JSON route: \033[34m%+v\033[0m from LUA\n", route.Route)
	} else {
		fmt.Printf("Registering route: \033[34m%+v\033[0m from LUA\n", route.Route)
	}
	handlers.RegisterRoutes(*route)
	L.Push(lua.LBool(true)) // Return true to indicate success
	return 1                // Number of return values
}
