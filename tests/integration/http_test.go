package integration

import (
	"testing"
	"time"

	"litepost/internal/app"
	"litepost/internal/http"
	"litepost/internal/models"
	"litepost/internal/storage/sqlite"
)

func TestHTTPRequest(t *testing.T) {
	// Create a temporary database
	storage, err := sqlite.NewSQLiteStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	// Create service
	service := app.NewService(storage)

	// Create a test request
	req := &models.Request{
		ID:     "test-1",
		Name:   "Test Request",
		Method: "GET",
		URL:    "https://httpbin.org/get",
		Headers: map[string]string{
			"User-Agent": "Litepost/1.0",
		},
		QueryParams: map[string]string{
			"test": "value",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Execute the request
	resp, err := service.ExecuteRequest(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	// Verify response
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	if resp.Body == "" {
		t.Error("Expected response body to not be empty")
	}

	// Check if response contains expected data
	if !contains(resp.Body, "httpbin.org") {
		t.Error("Expected response to contain 'httpbin.org'")
	}

	t.Logf("Request executed successfully: %d - %d bytes", resp.StatusCode, resp.Size)
}

func TestHTTPClient(t *testing.T) {
	client := http.NewClient(nil)
	if client == nil {
		t.Error("Expected HTTP client to be created")
	}

	// Test with a simple request
	req := &models.Request{
		ID:     "test-2",
		Name:   "Test Request",
		Method: "GET",
		URL:    "https://httpbin.org/status/200",
		Headers: map[string]string{
			"User-Agent": "Litepost/1.0",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
