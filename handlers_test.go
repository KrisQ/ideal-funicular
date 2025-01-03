package main

import (
	"bytes"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

// Test handlerCreateChirp
func TestHandlerCreateChirp(t *testing.T) {
	// Create an apiConfig with a mock DB
	mockDB := &database.MockDB{}
	cfg := apiConfig{
		db: mockDB,
	}

	// Build a sample request
	payload := `{"body": "Hello #chirpy", "user_id": "` + uuid.New().String() + `"}`
	req, err := http.NewRequest("POST", "/api/chirps", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	// Usually you'd set headers if needed, e.g. content-type.
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler
	http.HandlerFunc(cfg.handlerCreateChirp).ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, status)
	}

	// Check the response body
	var respBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	if err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	if respBody["body"] != "Hello #chirpy" {
		t.Errorf("expected 'Hello #chirpy' in response, got %v", respBody["body"])
	}
}

func TestHandlerGetAllChirps(t *testing.T) {
	mockDB := &database.MockDB{}
	cfg := apiConfig{
		db: mockDB,
	}

	req, err := http.NewRequest("GET", "/api/chirps", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rr := httptest.NewRecorder()

	http.HandlerFunc(cfg.handlerGetAllChirps).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %v", rr.Code)
	}

	var chirps []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &chirps)
	if err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	if len(chirps) == 0 {
		t.Error("expected at least one chirp from mock data, got zero")
	}
}
