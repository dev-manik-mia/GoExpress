package main

import (
	"net/http"

	"github.com/dev-manik-mia/goexpress"
)

// 1. The Data Model
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// 2. The In-Memory Datastore
// Maps a string ID (like "1") to a User struct
var userDB = make(map[string]User)

func main() {
	g := goexpress.New()

	// CREATE: Add a new user
	g.POST("/user", func(c *goexpress.Context) {
		var newUser User

		// Parse the incoming JSON into our struct
		if err := c.BindJSON(&newUser); err != nil {
			c.String(http.StatusBadRequest, "Invalid JSON data")
			return
		}

		// Save to our fake database (Hardcoding ID "1" for simplicity)
		userDB["1"] = newUser
		c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
	})

	// READ: Get the user
	g.GET("/user", func(c *goexpress.Context) {
		if user, exists := userDB["1"]; exists {
			c.JSON(http.StatusOK, user)
		} else {
			c.String(http.StatusNotFound, "User not found")
		}
	})

	// UPDATE: Modify the user
	g.PUT("/user", func(c *goexpress.Context) {
		var updatedUser User
		if err := c.BindJSON(&updatedUser); err != nil {
			c.String(http.StatusBadRequest, "Invalid JSON data")
			return
		}

		// Overwrite the existing data
		userDB["1"] = updatedUser
		c.JSON(http.StatusOK, map[string]string{"message": "User updated"})
	})

	// DELETE: Remove the user
	g.DELETE("/user", func(c *goexpress.Context) {
		delete(userDB, "1")
		c.String(http.StatusOK, "User deleted")
	})

	g.Run(":8082")
}
