package main

import (
	"net/http"

	"github.com/dev-manik-mia/goexpress"
)

func main() {
	// 1. Create a new, isolated instance of our framework
	g := goexpress.New()

	// 2. Register our route using the custom helper
	g.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to Go Framework Lab 1!\n"))
	})

	// 3. Launch the server, passing our Engine in to replace 'nil'
	g.Run(":8080")
}
