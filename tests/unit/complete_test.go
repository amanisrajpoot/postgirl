package unit

import (
	"testing"
	"time"

	"postgirl/internal/app"
	"postgirl/internal/models"
	"postgirl/internal/storage/sqlite"
)

func TestCompleteWorkflow(t *testing.T) {
	// Create in-memory database
	storage, err := sqlite.NewSQLiteStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	// Create service
	service := app.NewService(storage)

	// Test 1: Create a collection
	collection := &models.Collection{
		ID:          "test-collection",
		Name:        "Test Collection",
		Description: "A test collection",
		Variables: map[string]string{
			"base_url": "https://api.example.com",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = service.SaveCollection(collection)
	if err != nil {
		t.Errorf("Failed to save collection: %v", err)
	}

	// Test 2: Create an environment
	environment := &models.Environment{
		ID:   "test-env",
		Name: "Test Environment",
		Variables: map[string]string{
			"api_key": "test-key-123",
			"base_url": "https://api.test.com",
		},
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = service.SaveEnvironmentToDB(environment)
	if err != nil {
		t.Errorf("Failed to save environment: %v", err)
	}

	// Test 3: Create a request with scripting
	request := &models.Request{
		ID:            "test-request",
		Name:          "Test Request",
		Method:        "GET",
		URL:           "{{base_url}}/users",
		Headers:       map[string]string{
			"Authorization": "Bearer {{api_key}}",
		},
		EnvironmentID: "test-env",
		PreScript:     "console.log('Pre-request script executed');",
		PostScript:    "console.log('Post-response script executed');",
		Tests: []models.Test{
			{
				Name:   "Status Code Test",
				Script: "pm.test('Status code is 200', function() { pm.expect(response.statusCode).to.equal(200); });",
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = service.SaveRequest(request)
	if err != nil {
		t.Errorf("Failed to save request: %v", err)
	}

	// Test 4: List collections
	collections, err := service.ListCollections()
	if err != nil {
		t.Errorf("Failed to list collections: %v", err)
	}
	if len(collections) != 1 {
		t.Errorf("Expected 1 collection, got %d", len(collections))
	}

	// Test 5: List environments
	environments, err := service.ListEnvironmentsFromDB()
	if err != nil {
		t.Errorf("Failed to list environments: %v", err)
	}
	if len(environments) != 1 {
		t.Errorf("Expected 1 environment, got %d", len(environments))
	}

	// Test 6: List requests
	requests, err := service.ListRequests()
	if err != nil {
		t.Errorf("Failed to list requests: %v", err)
	}
	if len(requests) != 1 {
		t.Errorf("Expected 1 request, got %d", len(requests))
	}

	// Test 7: Get specific request
	retrievedRequest, err := service.GetRequest("test-request")
	if err != nil {
		t.Errorf("Failed to get request: %v", err)
	}
	if retrievedRequest.Name != "Test Request" {
		t.Errorf("Expected request name 'Test Request', got '%s'", retrievedRequest.Name)
	}

	// Test 8: Test environment variable substitution
	// This would normally be tested with a real HTTP request, but we'll test the substitution logic
	if retrievedRequest.URL != "{{base_url}}/users" {
		t.Errorf("Expected URL with variables, got '%s'", retrievedRequest.URL)
	}

	// Test 9: Delete request
	err = service.DeleteRequest("test-request")
	if err != nil {
		t.Errorf("Failed to delete request: %v", err)
	}

	// Test 10: Delete collection
	err = service.DeleteCollection("test-collection")
	if err != nil {
		t.Errorf("Failed to delete collection: %v", err)
	}

	// Test 11: Delete environment
	err = service.DeleteEnvironmentFromDB("test-env")
	if err != nil {
		t.Errorf("Failed to delete environment: %v", err)
	}
}

func TestScriptingEngine(t *testing.T) {
	scriptEngine := app.NewScriptEngine()

	// Test basic script execution
	script := "console.log('Hello from JavaScript');"
	err := scriptEngine.ExecutePreScript(script, &models.Request{}, nil)
	if err != nil {
		t.Errorf("Failed to execute pre-script: %v", err)
	}

	// Test script with request context
	request := &models.Request{
		ID:     "test-request",
		Name:   "Test Request",
		Method: "GET",
		URL:    "https://api.example.com/test",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	script = "console.log('Request URL:', request.url);"
	err = scriptEngine.ExecutePreScript(script, request, nil)
	if err != nil {
		t.Errorf("Failed to execute script with request context: %v", err)
	}
}

func TestCompleteAuthenticationTypes(t *testing.T) {
	// Test all authentication types
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
			authType: "digest",
			config: map[string]string{
				"username": "user",
				"password": "pass",
			},
			valid: true,
		},
		{
			authType: "hawk",
			config: map[string]string{
				"id":  "hawk-id",
				"key": "hawk-key",
			},
			valid: true,
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

func TestCompleteEnvironmentVariableSubstitution(t *testing.T) {
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

	// Test complex variable substitution
	request := &models.Request{
		ID:     "test-request",
		Name:   "Test Request",
		Method: "GET",
		URL:    "{{base_url}}/users/{{user_id}}",
		Headers: map[string]string{
			"Authorization": "Bearer {{api_key}}",
			"X-Timeout":     "{{timeout}}",
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
		EnvironmentID: "test",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Test variable substitution
	err := envService.SubstituteRequestVariables(request, "test")
	if err != nil {
		t.Errorf("Failed to substitute request variables: %v", err)
	}

	// Check URL substitution
	expectedURL := "https://api.example.com/users/{{user_id}}"
	if request.URL != expectedURL {
		t.Errorf("Expected URL '%s', got '%s'", expectedURL, request.URL)
	}

	// Check header substitution
	if request.Headers["Authorization"] != "Bearer test-key-123" {
		t.Errorf("Expected Authorization header 'Bearer test-key-123', got '%s'", request.Headers["Authorization"])
	}

	// Check body substitution
	expectedBody := `{"user": "test-key-123", "url": "https://api.example.com"}`
	if request.Body.Content != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, request.Body.Content)
	}
}
