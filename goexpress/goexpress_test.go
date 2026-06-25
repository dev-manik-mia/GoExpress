package goexpress

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test 1: Testing the upgraded Context String helper
func TestEngineRouting(t *testing.T) {
	engine := New()

	// NEW: Notice we are using our Lab 2 Context signature!
	engine.GET("/hello", func(c *Context) {
		c.String(http.StatusOK, "hello world")
	})

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
	if w.Body.String() != "hello world" {
		t.Errorf("Expected body 'hello world', got '%s'", w.Body.String())
	}
}

// Test 2: The 404 Fallback (Remains largely the same)
func TestEngineNotFound(t *testing.T) {
	engine := New()

	req := httptest.NewRequest("GET", "/missing-page", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", w.Code)
	}
}

// Test 3: NEW - Testing the Context JSON Helper
func TestContextJSON(t *testing.T) {
	engine := New()

	// 1. Setup a route that sends a JSON response
	engine.GET("/api/data", func(c *Context) {
		c.JSON(http.StatusCreated, map[string]string{"status": "ok"})
	})

	// 2. Create our fakes and execute
	req := httptest.NewRequest("GET", "/api/data", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// 3. Assertion A: Check the Status Code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", w.Code)
	}

	// 4. Assertion B: Check the HTTP Headers
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got '%s'", contentType)
	}

	// 5. Assertion C: Parse and verify the exact JSON payload
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected JSON status 'ok', got '%s'", response["status"])
	}
}
