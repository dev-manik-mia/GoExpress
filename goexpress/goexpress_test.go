package goexpress

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test 1: Does the Engine successfully find and execute a valid route?
func TestEngineRouting(t *testing.T) {
	// 1. Setup a fresh Engine
	engine := New()

	// 2. Register a test route
	engine.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	})

	// 3. Create a fake "offline" request (Simulating a browser)
	req := httptest.NewRequest("GET", "/hello", nil)

	// 4. Create a fake ResponseWriter to capture the output
	w := httptest.NewRecorder()

	// 5. Execute our Engine's core router directly
	engine.ServeHTTP(w, req)

	// 6. Assertions: Did the Engine do the right thing?
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	if w.Body.String() != "hello world" {
		t.Errorf("Expected body 'hello world', got '%s'", w.Body.String())
	}
}

// Test 2: Does the Engine correctly return a 404 when a route is missing?
func TestEngineNotFound(t *testing.T) {
	// 1. Setup an empty Engine (We are NOT registering any routes)
	engine := New()

	// 2. Create a fake request for a page that does not exist
	req := httptest.NewRequest("GET", "/missing-page", nil)
	w := httptest.NewRecorder()

	// 3. Execute
	engine.ServeHTTP(w, req)

	// 4. Assertions
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", w.Code)
	}
}
