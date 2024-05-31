package handlers

import (
	"io"
	"net/http"
	"os"
)

func RegisterRoutes(route Route) {
	// fmt.Println(route)
	http.HandleFunc(route.Route, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", route.ContentType)

		// fmt.Println(r.Method, r.URL)

		f, err := os.OpenFile(route.Content, os.O_RDONLY, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer f.Close()

		_, err = io.Copy(w, f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
