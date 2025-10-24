package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"postgirl/internal/app"
	"postgirl/internal/models"
)

// Server represents the web server
type Server struct {
	app    *app.Service
	server *http.Server
}

// NewServer creates a new web server
func NewServer(app *app.Service, port int, assets embed.FS) *Server {
	mux := http.NewServeMux()
	
	// Create server instance
	s := &Server{
		app: app,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
	
	// Setup routes
	s.setupRoutes(mux, assets)
	
	return s
}

// setupRoutes configures all the HTTP routes
func (s *Server) setupRoutes(serveMux *http.ServeMux, assets embed.FS) {
	// Create gorilla/mux router
	router := mux.NewRouter()
	
	// API routes (must come before static files)
	api := router.PathPrefix("/api").Subrouter()
	
	// Request routes
	api.HandleFunc("/requests", s.handleRequests).Methods("GET", "POST")
	api.HandleFunc("/requests/{id}", s.handleRequest).Methods("GET", "PUT", "DELETE")
	api.HandleFunc("/requests/{id}/execute", s.handleExecuteRequest).Methods("POST")
	
	// Response routes
	api.HandleFunc("/requests/{id}/responses", s.handleResponses).Methods("GET")
	
	// Collection routes
	api.HandleFunc("/collections", s.handleCollections).Methods("GET", "POST")
	api.HandleFunc("/collections/{id}", s.handleCollection).Methods("GET", "PUT", "DELETE")
	
	// Environment routes
	api.HandleFunc("/environments", s.handleEnvironments).Methods("GET", "POST")
	api.HandleFunc("/environments/{id}", s.handleEnvironment).Methods("GET", "PUT", "DELETE")
	
	// Health check
	router.HandleFunc("/health", s.handleHealth).Methods("GET")
	
	// Static files (catch-all for non-API routes)
	// Create a sub-filesystem for the static files
	staticFS, err := fs.Sub(assets, "web/static")
	if err != nil {
		log.Printf("Error creating static filesystem: %v", err)
		// Fallback to serving a simple message
		router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Static files not found", http.StatusNotFound)
		})
	} else {
		router.PathPrefix("/").Handler(http.FileServer(http.FS(staticFS)))
	}
	
	// Mount router
	serveMux.Handle("/", router)
}

// Start starts the web server
func (s *Server) Start() error {
	log.Printf("Starting web server on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

// Stop stops the web server
func (s *Server) Stop() error {
	return s.server.Shutdown(nil)
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// handleRequests handles request list and creation
func (s *Server) handleRequests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.getRequests(w, r)
	case "POST":
		s.createRequest(w, r)
	}
}

// handleRequest handles individual request operations
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	switch r.Method {
	case "GET":
		s.getRequest(w, r, id)
	case "PUT":
		s.updateRequest(w, r, id)
	case "DELETE":
		s.deleteRequest(w, r, id)
	}
}

// handleExecuteRequest handles request execution
func (s *Server) handleExecuteRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	// Get the request
	req, err := s.app.GetRequest(id)
	if err != nil {
		http.Error(w, "Request not found", http.StatusNotFound)
		return
	}
	
	// Execute the request
	resp, err := s.app.ExecuteRequest(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// getRequests returns all requests
func (s *Server) getRequests(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement request listing
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]models.Request{})
}

// createRequest creates a new request
func (s *Server) createRequest(w http.ResponseWriter, r *http.Request) {
	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Set default values
	req.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	
	// Save the request
	if err := s.app.SaveRequest(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

// getRequest returns a specific request
func (s *Server) getRequest(w http.ResponseWriter, r *http.Request, id string) {
	req, err := s.app.GetRequest(id)
	if err != nil {
		http.Error(w, "Request not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

// updateRequest updates a request
func (s *Server) updateRequest(w http.ResponseWriter, r *http.Request, id string) {
	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	req.ID = id
	req.UpdatedAt = time.Now()
	
	if err := s.app.SaveRequest(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

// deleteRequest deletes a request
func (s *Server) deleteRequest(w http.ResponseWriter, r *http.Request, id string) {
	// TODO: Implement request deletion
	w.WriteHeader(http.StatusNoContent)
}

// handleResponses returns responses for a request
func (s *Server) handleResponses(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	responses, err := s.app.GetResponses(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// handleCollections handles collection operations
func (s *Server) handleCollections(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement collection operations
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]models.Collection{})
}

// handleCollection handles individual collection operations
func (s *Server) handleCollection(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement collection operations
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// handleEnvironments handles environment operations
func (s *Server) handleEnvironments(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement environment operations
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]models.Environment{})
}

// handleEnvironment handles individual environment operations
func (s *Server) handleEnvironment(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement environment operations
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
