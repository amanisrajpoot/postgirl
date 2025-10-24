package storage

import "postgirl/internal/models"

// Storage defines the interface for data persistence
type Storage interface {
	// Request methods
	SaveRequest(req *models.Request) error
	GetRequest(id string) (*models.Request, error)
	GetAllRequests() ([]*models.Request, error)
	DeleteRequest(id string) error

	// Response methods
	SaveResponse(resp *models.Response) error
	GetResponsesForRequest(requestID string) ([]*models.Response, error)

	// Collection methods
	SaveCollection(coll *models.Collection) error
	GetCollection(id string) (*models.Collection, error)
	GetAllCollections() ([]*models.Collection, error)
	DeleteCollection(id string) error

	// Environment methods
	SaveEnvironment(env *models.Environment) error
	GetEnvironment(id string) (*models.Environment, error)
	GetAllEnvironments() ([]*models.Environment, error)
	DeleteEnvironment(id string) error
}
