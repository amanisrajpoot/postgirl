package http

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"litepost/internal/models"
)

// Client represents the HTTP client
type Client struct {
	restyClient *resty.Client
	config      *Config
}

// Config represents HTTP client configuration
type Config struct {
	Timeout     time.Duration
	RetryCount  int
	RetryDelay  time.Duration
	UserAgent   string
	FollowRedirects bool
}

// NewClient creates a new HTTP client
func NewClient(config *Config) *Client {
	if config == nil {
		config = &Config{
			Timeout:         30 * time.Second,
			RetryCount:      3,
			RetryDelay:      1 * time.Second,
			UserAgent:       "Litepost/1.0",
			FollowRedirects: true,
		}
	}

	client := resty.New()
	client.SetTimeout(config.Timeout)
	client.SetRetryCount(config.RetryCount)
	client.SetRetryWaitTime(config.RetryDelay)
	client.SetHeader("User-Agent", config.UserAgent)
	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(10))

	return &Client{
		restyClient: client,
		config:      config,
	}
}

// Execute executes an HTTP request
func (c *Client) Execute(req *models.Request) (*models.Response, error) {
	start := time.Now()

	// Create resty request
	r := c.restyClient.R()
	
	// Set headers
	if req.Headers != nil {
		r.SetHeaders(req.Headers)
	}

	// Set query parameters
	if req.QueryParams != nil {
		r.SetQueryParams(req.QueryParams)
	}

	// Set request body
	if req.Body != nil {
		switch req.Body.Type {
		case "json":
			r.SetHeader("Content-Type", "application/json")
			r.SetBody(req.Body.Content)
		case "xml":
			r.SetHeader("Content-Type", "application/xml")
			r.SetBody(req.Body.Content)
		case "form":
			r.SetHeader("Content-Type", "application/x-www-form-urlencoded")
			r.SetBody(req.Body.Content)
		case "raw":
			r.SetBody(req.Body.Content)
		}
	}

	// Apply authentication
	if req.Auth != nil {
		if err := c.applyAuth(r, req.Auth); err != nil {
			return nil, fmt.Errorf("failed to apply authentication: %w", err)
		}
	}

	// Set method and URL
	r.Method = req.Method
	r.URL = req.URL

	// Execute request
	resp, err := r.Send()
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	duration := time.Since(start)

	// Convert headers to map
	headers := make(map[string]string)
	for key, values := range resp.Header() {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	return &models.Response{
		ID:         generateID(),
		RequestID:  req.ID,
		StatusCode: resp.StatusCode(),
		Headers:    headers,
		Body:       string(resp.Body()),
		Size:       int64(len(resp.Body())),
		Duration:   duration,
		CreatedAt:  time.Now(),
	}, nil
}

// applyAuth applies authentication to the request
func (c *Client) applyAuth(r *resty.Request, auth *models.AuthConfig) error {
	switch auth.Type {
	case "basic":
		username, ok := auth.Config["username"]
		if !ok {
			return fmt.Errorf("username required for basic auth")
		}
		password, ok := auth.Config["password"]
		if !ok {
			return fmt.Errorf("password required for basic auth")
		}
		r.SetBasicAuth(username, password)

	case "bearer":
		token, ok := auth.Config["token"]
		if !ok {
			return fmt.Errorf("token required for bearer auth")
		}
		r.SetAuthToken(token)

	case "api_key":
		value, ok := auth.Config["value"]
		if !ok {
			return fmt.Errorf("value required for API key auth")
		}
		header, ok := auth.Config["header"]
		if !ok {
			header = "X-API-Key"
		}
		r.SetHeader(header, value)

	case "oauth2":
		// OAuth2 implementation
		accessToken, ok := auth.Config["access_token"]
		if !ok {
			return fmt.Errorf("access_token required for OAuth2 auth")
		}
		r.SetAuthToken(accessToken)
		
		// Add additional OAuth2 headers if needed
		if tokenType, ok := auth.Config["token_type"]; ok && tokenType != "" {
			r.SetHeader("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
		}

	case "digest":
		// Digest authentication
		username, ok := auth.Config["username"]
		if !ok {
			return fmt.Errorf("username required for digest auth")
		}
		password, ok := auth.Config["password"]
		if !ok {
			return fmt.Errorf("password required for digest auth")
		}
		// Note: resty doesn't have built-in digest auth, would need custom implementation
		r.SetBasicAuth(username, password) // Fallback to basic auth

	case "hawk":
		// Hawk authentication
		id, ok := auth.Config["id"]
		if !ok {
			return fmt.Errorf("id required for hawk auth")
		}
		key, ok := auth.Config["key"]
		if !ok {
			return fmt.Errorf("key required for hawk auth")
		}
		// Note: resty doesn't have built-in hawk auth, would need custom implementation
		r.SetHeader("Authorization", fmt.Sprintf("Hawk id=\"%s\", key=\"%s\"", id, key))

	default:
		return fmt.Errorf("unsupported authentication type: %s", auth.Type)
	}

	return nil
}

// generateID generates a unique ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
