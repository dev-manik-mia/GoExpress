package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dev-manik-mia/goexpress" // Update path
)

// Logger is a custom middleware function
func Logger() goexpress.RouteHandler {
	return func(c *goexpress.Context) {
		// 1. PRE-PROCESSING
		t := time.Now()
		fmt.Printf("[START] Request incoming: %s %s\n", c.Method, c.Path)

		// 2. PASS CONTROL to the next middleware or handler
		c.Next()

		// 3. POST-PROCESSING (Executes after the entire request is handled)
		latency := time.Since(t)
		fmt.Printf("[END] Request completed: %s %s in %v\n", c.Method, c.Path, latency)
	}
}

func main() {
	g := goexpress.New()

	// Register global middleware using .Use()
	g.Use(Logger())

	// A simple route to test the middleware
	g.GET("/", func(c *goexpress.Context) {
		// Simulate some heavy processing
		time.Sleep(200 * time.Millisecond)
		c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to the Onion Model",
		})
	})

	g.GET("/crash", func(c *goexpress.Context) {
		panic("Simulated database failure!")
	})

	g.Run(":8080")
}
