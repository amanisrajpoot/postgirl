package unit

import (
	"testing"
	"time"

	"postgirl/internal/app"
	"postgirl/internal/models"
)

func TestEnvironmentService(t *testing.T) {
	envService := app.NewEnvironmentService()
	
	// Test creating default environment
	defaultEnv := envService.CreateDefaultEnvironment()
	if defaultEnv.ID != "default" {
		t.Errorf("Expected default environment ID 'default', got '%s'", defaultEnv.ID)
	}
	
	if defaultEnv.Name != "Default Environment" {
		t.Errorf("Expected default environment name 'Default Environment', got '%s'", defaultEnv.Name)
	}
	
	// Test setting and getting environment
	envService.SetEnvironment(defaultEnv)
	retrievedEnv, err := envService.GetEnvironment("default")
	if err != nil {
		t.Errorf("Failed to get environment: %v", err)
	}
	
	if retrievedEnv.ID != defaultEnv.ID {
		t.Errorf("Expected environment ID 'default', got '%s'", retrievedEnv.ID)
	}
}

func TestEnvironmentVariableSubstitution(t *testing.T) {
	envService := app.NewEnvironmentService()
	
	// Create test environment
	env := &models.Environment{
		ID:   "test",
		Name: "Test Environment",
		Variables: map[string]string{
			"base_url": "https://api.example.com",
			"api_key":  "test-key-123",
			"timeout":  "30",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	envService.SetEnvironment(env)
	
	// Test variable substitution
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "{{base_url}}/users",
			expected: "https://api.example.com/users",
		},
		{
			input:    "https://{{base_url}}/api/v1",
			expected: "https://https://api.example.com/api/v1",
		},
		{
			input:    "{{api_key}}",
			expected: "test-key-123",
		},
		{
			input:    "{{unknown_var}}",
			expected: "{{unknown_var}}",
		},
		{
			input:    "no variables here",
			expected: "no variables here",
		},
	}
	
	for _, test := range tests {
		result, err := envService.SubstituteVariables(test.input, "test")
		if err != nil {
			t.Errorf("Failed to substitute variables in '%s': %v", test.input, err)
			continue
		}
		
		if result != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, result)
		}
	}
}

func TestRequestVariableSubstitution(t *testing.T) {
	envService := app.NewEnvironmentService()
	
	// Create test environment
	env := &models.Environment{
		ID:   "test",
		Name: "Test Environment",
		Variables: map[string]string{
			"base_url": "https://api.example.com",
			"api_key":  "test-key-123",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	envService.SetEnvironment(env)
	
	// Create test request
	req := &models.Request{
		ID:     "test-request",
		Name:   "Test Request",
		Method: "GET",
		URL:    "{{base_url}}/users",
		Headers: map[string]string{
			"Authorization": "Bearer {{api_key}}",
			"Content-Type":  "application/json",
		},
		QueryParams: map[string]string{
			"page": "1",
			"limit": "{{timeout}}",
		},
		Body: &models.RequestBody{
			Type:    "json",
			Content: `{"user": "{{api_key}}", "url": "{{base_url}}"}`,
		},
		Auth: &models.AuthConfig{
			Type: "api_key",
			Config: map[string]string{
				"value":  "{{api_key}}",
				"header": "X-API-Key",
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Test variable substitution
	err := envService.SubstituteRequestVariables(req, "test")
	if err != nil {
		t.Errorf("Failed to substitute request variables: %v", err)
	}
	
	// Check URL substitution
	if req.URL != "https://api.example.com/users" {
		t.Errorf("Expected URL 'https://api.example.com/users', got '%s'", req.URL)
	}
	
	// Check header substitution
	if req.Headers["Authorization"] != "Bearer test-key-123" {
		t.Errorf("Expected Authorization header 'Bearer test-key-123', got '%s'", req.Headers["Authorization"])
	}
	
	// Check body substitution
	expectedBody := `{"user": "test-key-123", "url": "https://api.example.com"}`
	if req.Body.Content != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, req.Body.Content)
	}
	
	// Check auth substitution
	if req.Auth.Config["value"] != "test-key-123" {
		t.Errorf("Expected auth value 'test-key-123', got '%s'", req.Auth.Config["value"])
	}
}

func TestAuthenticationTypes(t *testing.T) {
	// Test different authentication types
	authTests := []struct {
		authType string
		config   map[string]string
		valid    bool
	}{
		{
			authType: "basic",
			config: map[string]string{
				"username": "user",
				"password": "pass",
			},
			valid: true,
		},
		{
			authType: "bearer",
			config: map[string]string{
				"token": "bearer-token-123",
			},
			valid: true,
		},
		{
			authType: "api_key",
			config: map[string]string{
				"value":  "api-key-123",
				"header": "X-API-Key",
			},
			valid: true,
		},
		{
			authType: "oauth2",
			config: map[string]string{
				"access_token": "oauth-token-123",
				"token_type":   "Bearer",
			},
			valid: true,
		},
		{
			authType: "invalid",
			config:   map[string]string{},
			valid:    false,
		},
	}
	
	for _, test := range authTests {
		req := &models.Request{
			ID:     "test-request",
			Name:   "Test Request",
			Method: "GET",
			URL:    "https://api.example.com/test",
			Auth: &models.AuthConfig{
				Type:   test.authType,
				Config: test.config,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		// Test that the request can be created with the auth config
		if test.valid {
			if req.Auth.Type != test.authType {
				t.Errorf("Expected auth type '%s', got '%s'", test.authType, req.Auth.Type)
			}
		}
	}
}
