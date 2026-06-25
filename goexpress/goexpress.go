package goexpress

import (
	"fmt"
	"net/http"
)

// UPGRADE 1: The Secret Sauce Pays Off!
// We simply swap out (w, r) for our new (*Context)
type RouteHandler func(c *Context)

type Engine struct {
	router map[string]RouteHandler
}

func New() *Engine {
	return &Engine{router: make(map[string]RouteHandler)}
}

func (e *Engine) GET(path string, handler RouteHandler) {
	e.router["GET-"+path] = handler
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if routeHandler, ok := e.router[key]; ok {

		// UPGRADE 2: Package the raw Go variables into our custom Context
		c := newContext(w, r)

		// Pass the Context to the developer's code
		routeHandler(c)

	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}
}

func (e *Engine) Run(addr string) error {
	fmt.Printf("GoExpress is running on %s...\n", addr)
	return http.ListenAndServe(addr, e)
}
