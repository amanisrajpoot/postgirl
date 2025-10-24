package storage

import (
	"sync"
	"time"

	"postgirl/internal/models"
)

// MemoryStorage is an in-memory storage implementation
type MemoryStorage struct {
	requests    map[string]*models.Request
	responses   map[string]*models.Response
	collections map[string]*models.Collection
	environments map[string]*models.Environment
	mutex       sync.RWMutex
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		requests:     make(map[string]*models.Request),
		responses:    make(map[string]*models.Response),
		collections:  make(map[string]*models.Collection),
		environments: make(map[string]*models.Environment),
	}
}

// SaveRequest saves a request to memory
func (m *MemoryStorage) SaveRequest(req *models.Request) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	req.UpdatedAt = time.Now()
	m.requests[req.ID] = req
	return nil
}

// GetRequest retrieves a request by ID
func (m *MemoryStorage) GetRequest(id string) (*models.Request, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	req, exists := m.requests[id]
	if !exists {
		return nil, nil
	}
	return req, nil
}

// GetAllRequests returns all requests
func (m *MemoryStorage) GetAllRequests() ([]*models.Request, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	requests := make([]*models.Request, 0, len(m.requests))
	for _, req := range m.requests {
		requests = append(requests, req)
	}
	return requests, nil
}

// DeleteRequest deletes a request
func (m *MemoryStorage) DeleteRequest(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	delete(m.requests, id)
	return nil
}

// SaveResponse saves a response to memory
func (m *MemoryStorage) SaveResponse(resp *models.Response) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	resp.CreatedAt = time.Now()
	m.responses[resp.ID] = resp
	return nil
}

// GetResponsesForRequest returns all responses for a request
func (m *MemoryStorage) GetResponsesForRequest(requestID string) ([]*models.Response, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	responses := make([]*models.Response, 0)
	for _, resp := range m.responses {
		if resp.RequestID == requestID {
			responses = append(responses, resp)
		}
	}
	return responses, nil
}

// SaveCollection saves a collection to memory
func (m *MemoryStorage) SaveCollection(coll *models.Collection) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	coll.UpdatedAt = time.Now()
	m.collections[coll.ID] = coll
	return nil
}

// GetCollection retrieves a collection by ID
func (m *MemoryStorage) GetCollection(id string) (*models.Collection, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	coll, exists := m.collections[id]
	if !exists {
		return nil, nil
	}
	return coll, nil
}

// GetAllCollections returns all collections
func (m *MemoryStorage) GetAllCollections() ([]*models.Collection, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	collections := make([]*models.Collection, 0, len(m.collections))
	for _, coll := range m.collections {
		collections = append(collections, coll)
	}
	return collections, nil
}

// DeleteCollection deletes a collection
func (m *MemoryStorage) DeleteCollection(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	delete(m.collections, id)
	return nil
}

// SaveEnvironment saves an environment to memory
func (m *MemoryStorage) SaveEnvironment(env *models.Environment) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	env.UpdatedAt = time.Now()
	m.environments[env.ID] = env
	return nil
}

// GetEnvironment retrieves an environment by ID
func (m *MemoryStorage) GetEnvironment(id string) (*models.Environment, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	env, exists := m.environments[id]
	if !exists {
		return nil, nil
	}
	return env, nil
}

// GetAllEnvironments returns all environments
func (m *MemoryStorage) GetAllEnvironments() ([]*models.Environment, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	environments := make([]*models.Environment, 0, len(m.environments))
	for _, env := range m.environments {
		environments = append(environments, env)
	}
	return environments, nil
}

// DeleteEnvironment deletes an environment
func (m *MemoryStorage) DeleteEnvironment(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	delete(m.environments, id)
	return nil
}
