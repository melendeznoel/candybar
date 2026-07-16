package helper

import (
	"log"
	"net/http"
	"time"
)

// Logger will log a request info
func Logger(handler http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		handler.ServeHTTP(rw, r)

		log.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))
	})
}

// BuildRouter registers routes with the standard library mux.
// Method matching is enforced inside the wrapped handler.
func BuildRouter(routes Routes) *http.ServeMux {
	router := http.NewServeMux()

	for _, route := range routes {
		route := route

		var handler http.Handler

		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != route.Method {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			route.HandlerFunction(w, r)
		})

		handler = Logger(handler, route.Name)
		router.Handle(route.Pattern, handler)
	}

	return router
}
