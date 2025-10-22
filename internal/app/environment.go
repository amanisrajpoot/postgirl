package app

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"postgirl/internal/models"
)

// EnvironmentService handles environment variable operations
type EnvironmentService struct {
	environments map[string]*models.Environment
}

// NewEnvironmentService creates a new environment service
func NewEnvironmentService() *EnvironmentService {
	return &EnvironmentService{
		environments: make(map[string]*models.Environment),
	}
}

// SetEnvironment sets an environment
func (es *EnvironmentService) SetEnvironment(env *models.Environment) {
	es.environments[env.ID] = env
}

// GetEnvironment gets an environment by ID
func (es *EnvironmentService) GetEnvironment(id string) (*models.Environment, error) {
	env, exists := es.environments[id]
	if !exists {
		return nil, fmt.Errorf("environment not found: %s", id)
	}
	return env, nil
}

// ListEnvironments returns all environments
func (es *EnvironmentService) ListEnvironments() []*models.Environment {
	var envs []*models.Environment
	for _, env := range es.environments {
		envs = append(envs, env)
	}
	return envs
}

// SubstituteVariables substitutes environment variables in a string
func (es *EnvironmentService) SubstituteVariables(text string, envID string) (string, error) {
	if envID == "" {
		return text, nil
	}

	env, err := es.GetEnvironment(envID)
	if err != nil {
		return text, err
	}

	// Pattern to match {{variable}} syntax
	pattern := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	
	result := pattern.ReplaceAllStringFunc(text, func(match string) string {
		// Extract variable name from {{variable}}
		variableName := strings.Trim(match, "{}")
		
		// Check if variable exists in environment
		if value, exists := env.Variables[variableName]; exists {
			return value
		}
		
		// Return original match if variable not found
		return match
	})

	return result, nil
}

// SubstituteRequestVariables substitutes variables in a request
func (es *EnvironmentService) SubstituteRequestVariables(req *models.Request, envID string) error {
	var err error

	// Substitute URL
	req.URL, err = es.SubstituteVariables(req.URL, envID)
	if err != nil {
		return fmt.Errorf("failed to substitute URL variables: %w", err)
	}

	// Substitute headers
	for key, value := range req.Headers {
		req.Headers[key], err = es.SubstituteVariables(value, envID)
		if err != nil {
			return fmt.Errorf("failed to substitute header variables: %w", err)
		}
	}

	// Substitute query parameters
	for key, value := range req.QueryParams {
		req.QueryParams[key], err = es.SubstituteVariables(value, envID)
		if err != nil {
			return fmt.Errorf("failed to substitute query parameter variables: %w", err)
		}
	}

	// Substitute body content
	if req.Body != nil {
		req.Body.Content, err = es.SubstituteVariables(req.Body.Content, envID)
		if err != nil {
			return fmt.Errorf("failed to substitute body variables: %w", err)
		}
	}

	// Substitute auth config
	if req.Auth != nil {
		for key, value := range req.Auth.Config {
			req.Auth.Config[key], err = es.SubstituteVariables(value, envID)
			if err != nil {
				return fmt.Errorf("failed to substitute auth variables: %w", err)
			}
		}
	}

	return nil
}

// CreateDefaultEnvironment creates a default environment
func (es *EnvironmentService) CreateDefaultEnvironment() *models.Environment {
	return &models.Environment{
		ID:        "default",
		Name:      "Default Environment",
		Variables: map[string]string{
			"base_url": "https://api.example.com",
			"api_key":  "your-api-key-here",
			"timeout":  "30",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateDevelopmentEnvironment creates a development environment
func (es *EnvironmentService) CreateDevelopmentEnvironment() *models.Environment {
	return &models.Environment{
		ID:        "development",
		Name:      "Development Environment",
		Variables: map[string]string{
			"base_url": "http://localhost:3000",
			"api_key":  "dev-api-key",
			"timeout":  "10",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateProductionEnvironment creates a production environment
func (es *EnvironmentService) CreateProductionEnvironment() *models.Environment {
	return &models.Environment{
		ID:        "production",
		Name:      "Production Environment",
		Variables: map[string]string{
			"base_url": "https://api.production.com",
			"api_key":  "prod-api-key",
			"timeout":  "30",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
