package goexpress

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test for the new POST method and BindJSON functionality
func TestPostAndBindJSON(t *testing.T) {
	engine := New()

	// 1. Define a dummy struct for testing
	type Payload struct {
		Message string `json:"message"`
	}

	// 2. Register a POST route that reads JSON and echoes it back
	engine.POST("/echo", func(c *Context) {
		var p Payload
		if err := c.BindJSON(&p); err != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}

		// If successful, send a 201 Created status
		c.String(http.StatusCreated, "Received: "+p.Message)
	})

	// 3. Forge the incoming JSON Body
	// bytes.NewBuffer turns a raw string into an io.Reader stream that httptest requires
	jsonBody := []byte(`{"message": "hello framework"}`)
	bodyReader := bytes.NewBuffer(jsonBody)

	// 4. Create the fake POST request
	req := httptest.NewRequest("POST", "/echo", bodyReader)
	w := httptest.NewRecorder()

	// 5. Execute
	engine.ServeHTTP(w, req)

	// 6. Assertions
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	expectedBody := "Received: hello framework"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, w.Body.String())
	}
}

// Test for the PUT method routing and JSON parsing
func TestPutRouting(t *testing.T) {
	engine := New()

	type UpdatePayload struct {
		Role string `json:"role"`
	}

	// 1. Register a PUT route for updates
	engine.PUT("/users/1/role", func(c *Context) {
		var p UpdatePayload
		if err := c.BindJSON(&p); err != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}

		// If successful, confirm the update
		c.String(http.StatusOK, "Role updated to: "+p.Role)
	})

	// 2. Forge the incoming JSON Body
	jsonBody := []byte(`{"role": "admin"}`)
	bodyReader := bytes.NewBuffer(jsonBody)

	// 3. Create the fake PUT request
	req := httptest.NewRequest("PUT", "/users/1/role", bodyReader)
	w := httptest.NewRecorder()

	// 4. Execute
	engine.ServeHTTP(w, req)

	// 5. Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expectedBody := "Role updated to: admin"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, w.Body.String())
	}
}

// Test for the DELETE method routing
func TestDeleteRouting(t *testing.T) {
	engine := New()

	engine.DELETE("/remove", func(c *Context) {
		c.String(http.StatusOK, "deleted")
	})

	req := httptest.NewRequest("DELETE", "/remove", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
