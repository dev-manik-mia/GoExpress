package goexpress

import (
	"encoding/json"
	"net/http"
)

// RouteHandler defines the function signature for both middleware and standard routes
type RouteHandler func(c *Context)

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	Path   string
	Method string
	Params map[string]string

	// NEW: The execution chain and state pointer
	handlers []RouteHandler
	index    int8 // Tracks which handler in the chain we are currently executing
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1, // Start at -1, so the first Next() call increments it to 0
	}
}

// NEW: Next executes the pending handlers in the chain inside the calling handler.
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// ... (Keep existing Param, String, JSON, BindJSON methods below)

// Param retrieves a dynamic path parameter by its name
func (c *Context) Param(key string) string {
	return c.Params[key]
}

// --- HELPER METHODS (SAME AS BEFORE) ---

// String sends a plain text response with a status code
func (c *Context) String(code int, text string) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(text))
}

// JSON sends a formatted JSON response with a status code
func (c *Context) JSON(code int, obj interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)

	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// BindJSON reads the incoming HTTP request body and decodes it into a Go struct.
func (c *Context) BindJSON(obj interface{}) error {
	decoder := json.NewDecoder(c.Req.Body)
	defer c.Req.Body.Close()
	return decoder.Decode(obj)
}
