package goexpress

import (
	"fmt"
	"net/http"
)

// RouteHandler is our custom type for GoExpress route logic
type RouteHandler func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	// Updated the map to use the new RouteHandler type
	router map[string]RouteHandler
}

func New() *Engine {
	return &Engine{router: make(map[string]RouteHandler)}
}

// GET now accepts the GoExpress RouteHandler type
func (e *Engine) GET(path string, handler RouteHandler) {
	e.router["GET-"+path] = handler
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if routeHandler, ok := e.router[key]; ok {
		routeHandler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}
}

func (e *Engine) Run(addr string) error {
	fmt.Printf("GoExpress is running on %s...\n", addr)
	return http.ListenAndServe(addr, e)
}
