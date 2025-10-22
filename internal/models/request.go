package models

import (
	"time"
)

// Request represents an HTTP request
type Request struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Method        string            `json:"method"`
	URL           string            `json:"url"`
	Headers       map[string]string `json:"headers"`
	QueryParams   map[string]string `json:"query_params"`
	Body          *RequestBody      `json:"body"`
	Auth          *AuthConfig       `json:"auth"`
	PreScript     string            `json:"pre_script"`
	PostScript    string            `json:"post_script"`
	Tests         []Test            `json:"tests"`
	CollectionID  string            `json:"collection_id"`
	FolderID      string            `json:"folder_id"`
	EnvironmentID string            `json:"environment_id"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// RequestBody represents the body of an HTTP request
type RequestBody struct {
	Type    string `json:"type"`    // json, xml, form, raw
	Content string `json:"content"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Type   string            `json:"type"`   // basic, bearer, api_key, oauth2
	Config map[string]string `json:"config"`
}

// Test represents a test assertion
type Test struct {
	Name     string `json:"name"`
	Script   string `json:"script"`
	Expected string `json:"expected"`
}
