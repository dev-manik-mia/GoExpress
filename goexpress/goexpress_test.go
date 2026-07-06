package goexpress

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// ... (Keep existing CRUD tests from Lab 3)

func TestDynamicRoute(t *testing.T) {
	engine := New()

	// Register a dynamic route
	engine.GET("/users/:id", func(c *Context) {
		id := c.Param("id")
		c.String(http.StatusOK, "User ID is "+id)
	})

	// Test Case 1: Standard parameter extraction
	req := httptest.NewRequest("GET", "/users/456", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expected := "User ID is 456"
	if w.Body.String() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, w.Body.String())
	}
}

func TestWildcardRoute(t *testing.T) {
	engine := New()

	// Register a catch-all wildcard route
	engine.GET("/assets/*filepath", func(c *Context) {
		c.String(http.StatusOK, "File: "+c.Param("filepath"))
	})

	// Test Case 2: Multi-segment parameter capture
	req := httptest.NewRequest("GET", "/assets/css/main.css", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	expected := "File: css/main.css"
	if w.Body.String() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, w.Body.String())
	}
}
