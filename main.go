package main

import (
	"github.com/dev-manik-mia/goexpress"
)

func main() {
	g := goexpress.New()

	// Using the new c.String() helper
	g.GET("/", func(c *goexpress.Context) {
		c.String(200, "Welcome to GoExpress Lab 2! The context works.\n")
	})

	// Using the new c.JSON() helper
	g.GET("/api/info", func(c *goexpress.Context) {
		// We can pass a Go map directly, and the framework handles the JSON translation
		c.JSON(200, map[string]interface{}{
			"framework": "GoExpress",
			"version":   "2.0",
			"author":    "Manik Mia",
		})
	})

	g.Run(":8082")
}
