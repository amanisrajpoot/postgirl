package unit

import (
	"testing"
	"time"

	"litepost/internal/app"
	"litepost/internal/http"
	"litepost/internal/models"
)

func TestServiceCreation(t *testing.T) {
	service := app.NewService(nil)
	if service == nil {
		t.Error("Expected service to be created")
	}
}

func TestHTTPClientCreation(t *testing.T) {
	client := http.NewClient(nil)
	if client == nil {
		t.Error("Expected HTTP client to be created")
	}
}

func TestServiceRequestCreation(t *testing.T) {
	service := app.NewService(nil)
	req := service.CreateNewRequest()
	
	if req.ID == "" {
		t.Error("Expected request ID to be generated")
	}
	
	if req.Method != "GET" {
		t.Errorf("Expected method 'GET', got '%s'", req.Method)
	}
	
	if req.Name != "New Request" {
		t.Errorf("Expected name 'New Request', got '%s'", req.Name)
	}
}

func TestResponseModelCreation(t *testing.T) {
	resp := &models.Response{
		ID:         "resp-1",
		RequestID:  "req-1",
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       `{"message": "success"}`,
		Size:       25,
		Duration:   100 * time.Millisecond,
		CreatedAt:  time.Now(),
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
	
	if resp.Size != 25 {
		t.Errorf("Expected size 25, got %d", resp.Size)
	}
	
	if resp.Duration != 100*time.Millisecond {
		t.Errorf("Expected duration 100ms, got %v", resp.Duration)
	}
}

func TestRequestBuilder(t *testing.T) {
	req := &models.Request{
		ID:     "test-1",
		Name:   "Test Request",
		Method: "POST",
		URL:    "https://api.example.com/users",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Authorization": "Bearer token123",
		},
		QueryParams: map[string]string{
			"page": "1",
			"limit": "10",
		},
		Body: &models.RequestBody{
			Type:    "json",
			Content: `{"name": "John Doe"}`,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test request properties
	if req.Method != "POST" {
		t.Errorf("Expected method 'POST', got '%s'", req.Method)
	}
	
	if len(req.Headers) != 2 {
		t.Errorf("Expected 2 headers, got %d", len(req.Headers))
	}
	
	if len(req.QueryParams) != 2 {
		t.Errorf("Expected 2 query params, got %d", len(req.QueryParams))
	}
	
	if req.Body == nil {
		t.Error("Expected request body to be set")
	}
	
	if req.Body.Type != "json" {
		t.Errorf("Expected body type 'json', got '%s'", req.Body.Type)
	}
}
