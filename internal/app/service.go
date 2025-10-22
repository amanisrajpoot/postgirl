package app

import (
	"fmt"
	"time"

	"postgirl/internal/http"
	"postgirl/internal/models"
	"postgirl/internal/storage/sqlite"
)

// Service represents the main application service
type Service struct {
	httpClient        *http.Client
	storage           *sqlite.SQLiteStorage
	environmentService *EnvironmentService
	scriptEngine      *ScriptEngine
}

// NewService creates a new service instance
func NewService(storage *sqlite.SQLiteStorage) *Service {
	// Create environment service
	envService := NewEnvironmentService()
	
	// Add default environments
	envService.SetEnvironment(envService.CreateDefaultEnvironment())
	envService.SetEnvironment(envService.CreateDevelopmentEnvironment())
	envService.SetEnvironment(envService.CreateProductionEnvironment())
	
	// Create script engine
	scriptEngine := NewScriptEngine()
	
	return &Service{
		httpClient:        http.NewClient(nil),
		storage:           storage,
		environmentService: envService,
		scriptEngine:      scriptEngine,
	}
}

// ExecuteRequest executes an HTTP request and returns the response
func (s *Service) ExecuteRequest(req *models.Request) (*models.Response, error) {
	// Create a copy of the request to avoid modifying the original
	requestCopy := *req
	
	// Get environment if specified
	var environment *models.Environment
	if req.EnvironmentID != "" {
		var err error
		environment, err = s.environmentService.GetEnvironment(req.EnvironmentID)
		if err != nil {
			// Try to get from database
			environment, err = s.storage.GetEnvironment(req.EnvironmentID)
			if err != nil {
				return nil, fmt.Errorf("environment not found: %s", req.EnvironmentID)
			}
		}
	}
	
	// Apply environment variable substitution if environment is specified
	if environment != nil {
		if err := s.environmentService.SubstituteRequestVariables(&requestCopy, req.EnvironmentID); err != nil {
			return nil, fmt.Errorf("failed to substitute environment variables: %w", err)
		}
	}
	
	// Execute pre-request script
	if req.PreScript != "" {
		if err := s.scriptEngine.ExecutePreScript(req.PreScript, &requestCopy, environment); err != nil {
			return nil, fmt.Errorf("pre-script execution failed: %w", err)
		}
	}
	
	// Execute the HTTP request
	resp, err := s.httpClient.Execute(&requestCopy)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	
	// Execute post-response script
	if req.PostScript != "" {
		if err := s.scriptEngine.ExecutePostScript(req.PostScript, &requestCopy, resp, environment); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: post-script execution failed: %v\n", err)
		}
	}
	
	// Execute test scripts
	if len(req.Tests) > 0 {
		for _, test := range req.Tests {
			testResult, err := s.scriptEngine.ExecuteTestScript(test.Script, &requestCopy, resp, environment)
			if err != nil {
				fmt.Printf("Warning: test execution failed: %v\n", err)
			} else {
				fmt.Printf("Test '%s': %s\n", test.Name, testResult.Message)
			}
		}
	}

	// Save the response to storage
	if err := s.storage.SaveResponse(resp); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to save response: %v\n", err)
	}

	return resp, nil
}

// SaveRequest saves a request to storage
func (s *Service) SaveRequest(req *models.Request) error {
	req.UpdatedAt = time.Now()
	return s.storage.SaveRequest(req)
}

// GetRequest retrieves a request by ID
func (s *Service) GetRequest(id string) (*models.Request, error) {
	return s.storage.GetRequest(id)
}

// GetResponses retrieves responses for a request
func (s *Service) GetResponses(requestID string) ([]*models.Response, error) {
	return s.storage.GetResponses(requestID)
}

// CreateNewRequest creates a new request with default values
func (s *Service) CreateNewRequest() *models.Request {
	return &models.Request{
		ID:          generateID(),
		Name:        "New Request",
		Method:      "GET",
		URL:         "",
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
		Body:        nil,
		Auth:        nil,
		PreScript:   "",
		PostScript:  "",
		Tests:       []models.Test{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// GetEnvironment gets an environment by ID
func (s *Service) GetEnvironment(id string) (*models.Environment, error) {
	return s.environmentService.GetEnvironment(id)
}

// ListEnvironments returns all environments
func (s *Service) ListEnvironments() []*models.Environment {
	return s.environmentService.ListEnvironments()
}

// SetEnvironment sets an environment
func (s *Service) SetEnvironment(env *models.Environment) {
	s.environmentService.SetEnvironment(env)
}

// ListRequests returns all requests
func (s *Service) ListRequests() ([]*models.Request, error) {
	return s.storage.ListRequests()
}

// DeleteRequest deletes a request by ID
func (s *Service) DeleteRequest(id string) error {
	return s.storage.DeleteRequest(id)
}

// SaveCollection saves a collection
func (s *Service) SaveCollection(col *models.Collection) error {
	col.UpdatedAt = time.Now()
	return s.storage.SaveCollection(col)
}

// GetCollection retrieves a collection by ID
func (s *Service) GetCollection(id string) (*models.Collection, error) {
	return s.storage.GetCollection(id)
}

// ListCollections returns all collections
func (s *Service) ListCollections() ([]*models.Collection, error) {
	return s.storage.ListCollections()
}

// DeleteCollection deletes a collection by ID
func (s *Service) DeleteCollection(id string) error {
	return s.storage.DeleteCollection(id)
}

// SaveEnvironmentToDB saves an environment to the database
func (s *Service) SaveEnvironmentToDB(env *models.Environment) error {
	env.UpdatedAt = time.Now()
	return s.storage.SaveEnvironment(env)
}

// GetEnvironmentFromDB retrieves an environment from the database
func (s *Service) GetEnvironmentFromDB(id string) (*models.Environment, error) {
	return s.storage.GetEnvironment(id)
}

// ListEnvironmentsFromDB returns all environments from the database
func (s *Service) ListEnvironmentsFromDB() ([]*models.Environment, error) {
	return s.storage.ListEnvironments()
}

// DeleteEnvironmentFromDB deletes an environment from the database
func (s *Service) DeleteEnvironmentFromDB(id string) error {
	return s.storage.DeleteEnvironment(id)
}

// generateID generates a unique ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
