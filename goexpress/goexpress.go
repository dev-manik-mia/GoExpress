package goexpress

import (
	"fmt"
	"net/http"
)

// Engine is the core framework instance
type Engine struct {
	*RouterGroup // Struct Embedding: Engine inherits all methods of RouterGroup
	router       *router
}

// New creates a new GoExpress Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	// The Engine initializes itself as the absolute Root Group ("")
	engine.RouterGroup = &RouterGroup{engine: engine}
	return engine
}

// NEW: Use adds global middlewares to the framework instance
func (e *Engine) Use(middlewares ...RouteHandler) {
	e.middlewares = append(e.middlewares, middlewares...)
}

func (e *Engine) GET(pattern string, handler RouteHandler) {
	e.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler RouteHandler) {
	e.addRoute("POST", pattern, handler)
}

func (e *Engine) PUT(pattern string, handler RouteHandler) {
	e.addRoute("PUT", pattern, handler)
}

func (e *Engine) DELETE(pattern string, handler RouteHandler) {
	e.addRoute("DELETE", pattern, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)

	node, params := e.router.getRoute(r.Method, r.URL.Path)
	if node != nil {
		c.Params = params
		key := r.Method + "-" + node.pattern

		// Fetch the pre-calculated execution chain directly from the router
		c.handlers = e.router.handlers[key]
	} else {
		// 404 Handler
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND")
		})
	}

	// Kick off the execution chain
	c.Next()
}

func (e *Engine) Run(addr string) error {
	fmt.Printf("GoExpress is running on %s...\n", addr)
	return http.ListenAndServe(addr, e)
}
