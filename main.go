package main

import (
	"net/http"
	"strconv"

	"github.com/dev-manik-mia/goexpress" // Update path accordingly
)

// 1. The Data Model
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// 2. The In-Memory Datastore & Auto-Incrementing ID
var userDB = make(map[string]User)
var idCounter = 1

func main() {
	g := goexpress.New()

	// CREATE: Add a new user (Targeting the collection: /users)
	g.POST("/users", func(c *goexpress.Context) {
		var newUser User

		if err := c.BindJSON(&newUser); err != nil {
			c.String(http.StatusBadRequest, "Invalid JSON data")
			return
		}

		// Generate a new ID as a string
		id := strconv.Itoa(idCounter)
		idCounter++

		userDB[id] = newUser
		c.JSON(http.StatusCreated, map[string]string{
			"id":      id,
			"message": "User created successfully",
		})
	})

	// READ ALL: Get the entire collection of users
	g.GET("/users", func(c *goexpress.Context) {
		// The framework's JSON helper automatically encodes the entire map
		if len(userDB) == 0 {
			// Return an empty object or message if no users exist yet
			c.JSON(http.StatusOK, map[string]string{"message": "No users found"})
			return
		}
		c.JSON(http.StatusOK, userDB)
	})

	// READ: Get a specific user by dynamic ID (Targeting an item: /users/:id)
	g.GET("/users/:id", func(c *goexpress.Context) {
		id := c.Param("id") // Extract ID from the URL

		if user, exists := userDB[id]; exists {
			c.JSON(http.StatusOK, user)
		} else {
			c.String(http.StatusNotFound, "User not found")
		}
	})

	// UPDATE: Modify a specific user by dynamic ID
	g.PUT("/users/:id", func(c *goexpress.Context) {
		id := c.Param("id") // Extract ID from the URL

		// Ensure the user exists before updating
		if _, exists := userDB[id]; !exists {
			c.String(http.StatusNotFound, "User not found")
			return
		}

		var updatedUser User
		if err := c.BindJSON(&updatedUser); err != nil {
			c.String(http.StatusBadRequest, "Invalid JSON data")
			return
		}

		// Overwrite the existing data at this specific ID
		userDB[id] = updatedUser
		c.JSON(http.StatusOK, map[string]string{"message": "User " + id + " updated"})
	})

	// DELETE: Remove a specific user by dynamic ID
	g.DELETE("/users/:id", func(c *goexpress.Context) {
		id := c.Param("id") // Extract ID from the URL

		if _, exists := userDB[id]; exists {
			delete(userDB, id)
			c.String(http.StatusOK, "User "+id+" deleted")
		} else {
			c.String(http.StatusNotFound, "User not found")
		}
	})

	// (Optional) Keep the wildcard route from earlier to show multiple features coexisting
	g.GET("/static/*filepath", func(c *goexpress.Context) {
		file := c.Param("filepath")
		c.String(http.StatusOK, "Simulating serving file: "+file)
	})

	g.Run(":8080")
}
