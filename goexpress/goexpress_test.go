package goexpress

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouteGroupingAndMiddleware(t *testing.T) {
	engine := New()

	// Tracker array for execution flow
	var executionOrder []string

	// 1. Apply Global Middleware to the Engine (Root Group)
	engine.Use(func(c *Context) {
		executionOrder = append(executionOrder, "Global_Middleware")
		c.Next()
	})

	// 2. Create a generic API Group (inherits Global)
	api := engine.Group("/api")

	// 3. Create a versioned subgroup under API (inherits Global)
	v1 := api.Group("/v1")
	// Apply middleware specifically to v1 ONLY
	v1.Use(func(c *Context) {
		executionOrder = append(executionOrder, "V1_Middleware")
		c.Next()
	})

	// 4. Attach the core handler
	v1.GET("/users", func(c *Context) {
		executionOrder = append(executionOrder, "V1_Users_Handler")
		c.String(http.StatusOK, "OK")
	})

	// Simulate an incoming request
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// Expected order: Parent -> Child -> Core
	expectedOrder := []string{"Global_Middleware", "V1_Middleware", "V1_Users_Handler"}

	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d steps, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, v := range executionOrder {
		if v != expectedOrder[i] {
			t.Errorf("At index %d: expected %s, got %s", i, expectedOrder[i], v)
		}
	}
}
