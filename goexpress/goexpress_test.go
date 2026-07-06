package goexpress

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareOrder(t *testing.T) {
	engine := New()

	// We will track the order of execution in this slice
	var executionOrder []string

	// Middleware A (Outer Layer)
	engine.Use(func(c *Context) {
		executionOrder = append(executionOrder, "A_Before")
		c.Next() // Suspend and go deeper
		executionOrder = append(executionOrder, "A_After")
	})

	// Middleware B (Inner Layer)
	engine.Use(func(c *Context) {
		executionOrder = append(executionOrder, "B_Before")
		c.Next() // Suspend and go to core handler
		executionOrder = append(executionOrder, "B_After")
	})

	// Core Route Handler
	engine.GET("/test", func(c *Context) {
		executionOrder = append(executionOrder, "Core_Handler")
		c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// Validate the Onion Model flow
	expectedOrder := []string{"A_Before", "B_Before", "Core_Handler", "B_After", "A_After"}

	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d steps, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, v := range executionOrder {
		if v != expectedOrder[i] {
			t.Errorf("At index %d: expected %s, got %s", i, expectedOrder[i], v)
		}
	}
}
