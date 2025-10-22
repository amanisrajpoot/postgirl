package models

import (
	"time"
)

// Environment represents an environment with variables
type Environment struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Variables map[string]string `json:"variables"`
	IsActive  bool              `json:"is_active"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// EnvironmentVariable represents a single environment variable
type EnvironmentVariable struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// EnvironmentTemplate represents a template for creating environments
type EnvironmentTemplate struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Variables   []EnvironmentVariable  `json:"variables"`
}
