package handlers

import (
	"io"
	"net/http"
	"os"

	lua "github.com/yuin/gopher-lua"
)

func RegisterRoutes(route Route) {
	// fmt.Println(route)
	http.HandleFunc(route.Route, func(w http.ResponseWriter, r *http.Request) {
		if route.ContentType == "application/json" && route.LuaFunc != nil {
			state := lua.NewState()
			defer state.Close()
			state.SetGlobal("request", lua.LString(r.Method))
			state.SetGlobal("url", lua.LString(r.URL.String()))
			if err := state.CallByParam(lua.P{
				Fn:      route.LuaFunc,
				NRet:    1,
				Protect: true,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ret := state.Get(-1)
			state.Pop(1)
			w.Header().Set("Content-Type", route.ContentType)
			w.Write([]byte(ret.String()))
		} else {
			w.Header().Set("Content-Type", route.ContentType)

			f, err := os.OpenFile(route.Content, os.O_RDONLY, os.ModePerm)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()

			_, err = io.Copy(w, f)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})
}
