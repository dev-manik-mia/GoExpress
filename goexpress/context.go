package goexpress

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	Path   string
	Method string
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
	}
}

// --- HELPER METHODS ---

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

	// Use Go's built-in JSON encoder to translate the object and write it
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// BindJSON reads the incoming HTTP request body and decodes it into a Go struct.
func (c *Context) BindJSON(obj interface{}) error {
	decoder := json.NewDecoder(c.Req.Body)
	// Important: We must close the body after reading it to prevent memory leaks
	defer c.Req.Body.Close()
	return decoder.Decode(obj)
}
