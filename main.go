package main

import (
	"fmt"
	"net/http"

	"github.com/dev-manik-mia/goexpress" // Update to your actual module path
)

// Global Middleware
func GlobalLogger() goexpress.RouteHandler {
	return func(c *goexpress.Context) {
		fmt.Printf("[GLOBAL OMNIPRESENT LOG] Intercepted: %s %s\n", c.Method, c.Path)
		c.Next()
	}
}

// Middleware for /admin Group
func AdminGuard() goexpress.RouteHandler {
	return func(c *goexpress.Context) {
		token := c.Req.Header.Get("X-Admin-Token")
		if token != "super-secret-admin-pass" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"status": "Rejected",
				"reason": "Administrative clearance token missing or invalid.",
			})
			return // Short-circuit the execution chain!
		}
		fmt.Println("[GUARD] Clearance confirmed. Transitioning inward...")
		c.Next()
	}
}

// Middleware for /something-else Group
func ContextualTracker() goexpress.RouteHandler {
	return func(c *goexpress.Context) {
		fmt.Println("[TRACKER] Request routed directly into the 'Something-Else' cluster.")
		c.Next()
	}
}

func main() {
	g := goexpress.New()

	// Apply Global Middleware
	g.Use(GlobalLogger())

	// ==========================================
	// GROUP 1: /admin
	// ==========================================
	admin := g.Group("/admin")
	admin.Use(AdminGuard())
	{
		// Resolves to: PUT /admin/users/:id
		admin.PUT("/users/:id", func(c *goexpress.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, map[string]string{
				"action":  "PUT",
				"message": fmt.Sprintf("User structure for ID %s has been completely overwritten.", id),
			})
		})

		// Resolves to: DELETE /admin/users/:id
		admin.DELETE("/users/:id", func(c *goexpress.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, map[string]string{
				"action":  "DELETE",
				"message": fmt.Sprintf("User record with ID %s purged permanently from system storage.", id),
			})
		})
	}

	// ==========================================
	// GROUP 2: /something-else
	// ==========================================
	somethingElse := g.Group("/something-else")
	somethingElse.Use(ContextualTracker())
	{
		// Resolves to: PATCH /something-else/configurations
		somethingElse.PATCH("/configurations", func(c *goexpress.Context) {
			c.JSON(http.StatusOK, map[string]string{
				"action":  "PATCH",
				"message": "Specific configuration delta settings applied successfully.",
			})
		})
	}

	g.Run(":8080")
}
