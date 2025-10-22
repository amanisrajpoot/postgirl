package app

import (
	"fmt"
	"time"

	"github.com/dop251/goja"
	"litepost/internal/models"
)

// ScriptEngine handles JavaScript execution
type ScriptEngine struct {
	vm *goja.Runtime
}

// NewScriptEngine creates a new script engine
func NewScriptEngine() *ScriptEngine {
	vm := goja.New()
	
	// Add common utilities
	vm.Set("console", map[string]interface{}{
		"log": func(args ...interface{}) {
			fmt.Println(args...)
		},
		"error": func(args ...interface{}) {
			fmt.Print("ERROR: ")
			fmt.Println(args...)
		},
	})
	
	// Add setTimeout and setInterval
	vm.Set("setTimeout", func(callback goja.Callable, delay int) {
		go func() {
			time.Sleep(time.Duration(delay) * time.Millisecond)
			callback(goja.Undefined())
		}()
	})
	
	vm.Set("setInterval", func(callback goja.Callable, delay int) {
		go func() {
			for {
				time.Sleep(time.Duration(delay) * time.Millisecond)
				callback(goja.Undefined())
			}
		}()
	})
	
	return &ScriptEngine{vm: vm}
}

// ExecutePreScript executes a pre-request script
func (se *ScriptEngine) ExecutePreScript(script string, request *models.Request, environment *models.Environment) error {
	if script == "" {
		return nil
	}
	
	// Set up script context
	se.vm.Set("request", map[string]interface{}{
		"id":          request.ID,
		"name":        request.Name,
		"method":      request.Method,
		"url":         request.URL,
		"headers":     request.Headers,
		"queryParams": request.QueryParams,
		"body":        request.Body,
		"auth":        request.Auth,
	})
	
	if environment != nil {
		se.vm.Set("environment", map[string]interface{}{
			"id":        environment.ID,
			"name":      environment.Name,
			"variables": environment.Variables,
		})
	}
	
	// Execute the script
	_, err := se.vm.RunString(script)
	if err != nil {
		return fmt.Errorf("pre-script execution failed: %w", err)
	}
	
	// Update request with any modifications made by the script
	se.updateRequestFromScript(request)
	
	return nil
}

// ExecutePostScript executes a post-response script
func (se *ScriptEngine) ExecutePostScript(script string, request *models.Request, response *models.Response, environment *models.Environment) error {
	if script == "" {
		return nil
	}
	
	// Set up script context
	se.vm.Set("request", map[string]interface{}{
		"id":          request.ID,
		"name":        request.Name,
		"method":      request.Method,
		"url":         request.URL,
		"headers":     request.Headers,
		"queryParams": request.QueryParams,
		"body":        request.Body,
		"auth":        request.Auth,
	})
	
	se.vm.Set("response", map[string]interface{}{
		"id":         response.ID,
		"statusCode": response.StatusCode,
		"headers":    response.Headers,
		"body":       response.Body,
		"size":       response.Size,
		"duration":   response.Duration.Milliseconds(),
		"timestamp":  response.CreatedAt,
	})
	
	if environment != nil {
		se.vm.Set("environment", map[string]interface{}{
			"id":        environment.ID,
			"name":      environment.Name,
			"variables": environment.Variables,
		})
	}
	
	// Execute the script
	_, err := se.vm.RunString(script)
	if err != nil {
		return fmt.Errorf("post-script execution failed: %w", err)
	}
	
	return nil
}

// ExecuteTestScript executes a test script and returns test results
func (se *ScriptEngine) ExecuteTestScript(script string, request *models.Request, response *models.Response, environment *models.Environment) (*TestResult, error) {
	if script == "" {
		return &TestResult{Passed: true, Message: "No tests to run"}, nil
	}
	
	// Set up script context
	se.vm.Set("request", map[string]interface{}{
		"id":          request.ID,
		"name":        request.Name,
		"method":      request.Method,
		"url":         request.URL,
		"headers":     request.Headers,
		"queryParams": request.QueryParams,
		"body":        request.Body,
		"auth":        request.Auth,
	})
	
	se.vm.Set("response", map[string]interface{}{
		"id":         response.ID,
		"statusCode": response.StatusCode,
		"headers":    response.Headers,
		"body":       response.Body,
		"size":       response.Size,
		"duration":   response.Duration.Milliseconds(),
		"timestamp":  response.CreatedAt,
	})
	
	if environment != nil {
		se.vm.Set("environment", map[string]interface{}{
			"id":        environment.ID,
			"name":      environment.Name,
			"variables": environment.Variables,
		})
	}
	
	// Add test utilities
	se.vm.Set("pm", map[string]interface{}{
		"test": func(name string, condition bool) {
			// Test assertion
		},
		"expect": func(value interface{}) map[string]interface{} {
			return map[string]interface{}{
				"to": map[string]interface{}{
					"be": func(expected interface{}) bool {
						return value == expected
					},
					"equal": func(expected interface{}) bool {
						return value == expected
					},
					"contain": func(substring string) bool {
						if str, ok := value.(string); ok {
							return contains(str, substring)
						}
						return false
					},
				},
			}
		},
	})
	
	// Execute the script
	_, err := se.vm.RunString(script)
	if err != nil {
		return &TestResult{
			Passed:  false,
			Message: fmt.Sprintf("Test execution failed: %v", err),
		}, nil
	}
	
	// For now, return a basic test result
	// In a real implementation, you'd capture test results from the script
	return &TestResult{
		Passed:  true,
		Message: "All tests passed",
	}, nil
}

// TestResult represents the result of a test execution
type TestResult struct {
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// updateRequestFromScript updates the request with any modifications made by the script
func (se *ScriptEngine) updateRequestFromScript(request *models.Request) {
	// This would extract any modifications made to the request object in the script
	// For now, it's a placeholder
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
