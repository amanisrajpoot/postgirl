package unit

import (
	"testing"
	"postgirl/internal/models"
	"postgirl/internal/http"
)

func TestRequestModel(t *testing.T) {
	req := &models.Request{
		ID:     "test-1",
		Name:   "Test Request",
		Method: "GET",
		URL:    "https://httpbin.org/get",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	if req.ID != "test-1" {
		t.Errorf("Expected ID 'test-1', got '%s'", req.ID)
	}

	if req.Method != "GET" {
		t.Errorf("Expected method 'GET', got '%s'", req.Method)
	}
}

func TestHTTPClient(t *testing.T) {
	client := http.NewClient(nil)
	if client == nil {
		t.Error("Expected HTTP client to be created")
	}
}

func TestResponseModel(t *testing.T) {
	resp := &models.Response{
		ID:         "resp-1",
		RequestID:  "req-1",
		StatusCode: 200,
		Body:       "Hello World",
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}
