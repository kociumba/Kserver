package handlers

import (
	"io"
	"net/http"
	"os"
)

func RegisterRoutes(route Route) {
	// fmt.Printf("FROM REGISTERROUTES Registering route: %+v\n", route) // Log the route being registered
	http.HandleFunc(route.Route, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", route.ContentType)

		// fmt.Printf("FROM REGISTERROUTES Handling request for route: %s\n", route.Route) // Log handling request

		f, err := os.OpenFile(route.Content, os.O_RDONLY, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			// fmt.Printf("FROM REGISTERROUTES Error opening file %s: %v\n", route.Content, err) // Log the error
			return
		}
		defer f.Close()

		_, err = io.Copy(w, f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			// fmt.Printf("FROM REGISTERROUTES Error copying file %s: %v\n", route.Content, err) // Log the error
		}
	})
}
