package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"postgirl/internal/models"
)

// SQLiteStorage represents the SQLite storage implementation
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage creates a new SQLite storage instance
func NewSQLiteStorage(path string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &SQLiteStorage{db: db}
	
	// Create tables
	if err := storage.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return storage, nil
}

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

// createTables creates all necessary tables
func (s *SQLiteStorage) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS requests (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			method TEXT NOT NULL,
			url TEXT NOT NULL,
			headers TEXT,
			query_params TEXT,
			body TEXT,
			auth TEXT,
			pre_script TEXT,
			post_script TEXT,
			tests TEXT,
			collection_id TEXT,
			folder_id TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS collections (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			variables TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS environments (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			variables TEXT,
			is_active BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS responses (
			id TEXT PRIMARY KEY,
			request_id TEXT NOT NULL,
			status_code INTEGER,
			headers TEXT,
			body TEXT,
			size INTEGER,
			duration INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}

// SaveRequest saves a request to the database
func (s *SQLiteStorage) SaveRequest(req *models.Request) error {
	headers, _ := json.Marshal(req.Headers)
	queryParams, _ := json.Marshal(req.QueryParams)
	body, _ := json.Marshal(req.Body)
	auth, _ := json.Marshal(req.Auth)
	tests, _ := json.Marshal(req.Tests)

	query := `INSERT OR REPLACE INTO requests 
		(id, name, method, url, headers, query_params, body, auth, pre_script, post_script, tests, collection_id, folder_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query,
		req.ID, req.Name, req.Method, req.URL,
		string(headers), string(queryParams), string(body), string(auth),
		req.PreScript, req.PostScript, string(tests),
		req.CollectionID, req.FolderID, req.CreatedAt, req.UpdatedAt)

	return err
}

// GetRequest retrieves a request by ID
func (s *SQLiteStorage) GetRequest(id string) (*models.Request, error) {
	query := `SELECT id, name, method, url, headers, query_params, body, auth, pre_script, post_script, tests, collection_id, folder_id, created_at, updated_at
		FROM requests WHERE id = ?`

	row := s.db.QueryRow(query, id)
	
	var req models.Request
	var headers, queryParams, body, auth, tests string
	
	err := row.Scan(
		&req.ID, &req.Name, &req.Method, &req.URL,
		&headers, &queryParams, &body, &auth, &tests,
		&req.PreScript, &req.PostScript,
		&req.CollectionID, &req.FolderID, &req.CreatedAt, &req.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// Unmarshal JSON fields
	json.Unmarshal([]byte(headers), &req.Headers)
	json.Unmarshal([]byte(queryParams), &req.QueryParams)
	json.Unmarshal([]byte(body), &req.Body)
	json.Unmarshal([]byte(auth), &req.Auth)
	json.Unmarshal([]byte(tests), &req.Tests)

	return &req, nil
}

// SaveResponse saves a response to the database
func (s *SQLiteStorage) SaveResponse(resp *models.Response) error {
	headers, _ := json.Marshal(resp.Headers)

	query := `INSERT INTO responses 
		(id, request_id, status_code, headers, body, size, duration, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query,
		resp.ID, resp.RequestID, resp.StatusCode,
		string(headers), resp.Body, resp.Size, resp.Duration.Milliseconds(), resp.CreatedAt)

	return err
}

// GetResponses retrieves responses for a request
func (s *SQLiteStorage) GetResponses(requestID string) ([]*models.Response, error) {
	query := `SELECT id, request_id, status_code, headers, body, size, duration, created_at
		FROM responses WHERE request_id = ? ORDER BY created_at DESC`

	rows, err := s.db.Query(query, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []*models.Response
	for rows.Next() {
		var resp models.Response
		var headers string
		var duration int64

		err := rows.Scan(
			&resp.ID, &resp.RequestID, &resp.StatusCode,
			&headers, &resp.Body, &resp.Size, &duration, &resp.CreatedAt)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(headers), &resp.Headers)
		resp.Duration = time.Duration(duration) * time.Millisecond
		responses = append(responses, &resp)
	}

	return responses, nil
}

// ListRequests returns all requests
func (s *SQLiteStorage) ListRequests() ([]*models.Request, error) {
	query := `SELECT id, name, method, url, headers, query_params, body, auth, pre_script, post_script, tests, collection_id, folder_id, created_at, updated_at
		FROM requests ORDER BY updated_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*models.Request
	for rows.Next() {
		var req models.Request
		var headers, queryParams, body, auth, tests string
		
		err := rows.Scan(
			&req.ID, &req.Name, &req.Method, &req.URL,
			&headers, &queryParams, &body, &auth, &tests,
			&req.PreScript, &req.PostScript,
			&req.CollectionID, &req.FolderID, &req.CreatedAt, &req.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSON fields
		json.Unmarshal([]byte(headers), &req.Headers)
		json.Unmarshal([]byte(queryParams), &req.QueryParams)
		json.Unmarshal([]byte(body), &req.Body)
		json.Unmarshal([]byte(auth), &req.Auth)
		json.Unmarshal([]byte(tests), &req.Tests)

		requests = append(requests, &req)
	}

	return requests, nil
}

// SaveCollection saves a collection to the database
func (s *SQLiteStorage) SaveCollection(col *models.Collection) error {
	variables, _ := json.Marshal(col.Variables)

	query := `INSERT OR REPLACE INTO collections 
		(id, name, description, variables, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query,
		col.ID, col.Name, col.Description, string(variables), col.CreatedAt, col.UpdatedAt)

	return err
}

// GetCollection retrieves a collection by ID
func (s *SQLiteStorage) GetCollection(id string) (*models.Collection, error) {
	query := `SELECT id, name, description, variables, created_at, updated_at
		FROM collections WHERE id = ?`

	row := s.db.QueryRow(query, id)
	
	var col models.Collection
	var variables string
	
	err := row.Scan(
		&col.ID, &col.Name, &col.Description, &variables, &col.CreatedAt, &col.UpdatedAt)

	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(variables), &col.Variables)
	return &col, nil
}

// ListCollections returns all collections
func (s *SQLiteStorage) ListCollections() ([]*models.Collection, error) {
	query := `SELECT id, name, description, variables, created_at, updated_at
		FROM collections ORDER BY updated_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []*models.Collection
	for rows.Next() {
		var col models.Collection
		var variables string
		
		err := rows.Scan(
			&col.ID, &col.Name, &col.Description, &variables, &col.CreatedAt, &col.UpdatedAt)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(variables), &col.Variables)
		collections = append(collections, &col)
	}

	return collections, nil
}

// SaveEnvironment saves an environment to the database
func (s *SQLiteStorage) SaveEnvironment(env *models.Environment) error {
	variables, _ := json.Marshal(env.Variables)

	query := `INSERT OR REPLACE INTO environments 
		(id, name, variables, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query,
		env.ID, env.Name, string(variables), env.IsActive, env.CreatedAt, env.UpdatedAt)

	return err
}

// GetEnvironment retrieves an environment by ID
func (s *SQLiteStorage) GetEnvironment(id string) (*models.Environment, error) {
	query := `SELECT id, name, variables, is_active, created_at, updated_at
		FROM environments WHERE id = ?`

	row := s.db.QueryRow(query, id)
	
	var env models.Environment
	var variables string
	
	err := row.Scan(
		&env.ID, &env.Name, &variables, &env.IsActive, &env.CreatedAt, &env.UpdatedAt)

	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(variables), &env.Variables)
	return &env, nil
}

// ListEnvironments returns all environments
func (s *SQLiteStorage) ListEnvironments() ([]*models.Environment, error) {
	query := `SELECT id, name, variables, is_active, created_at, updated_at
		FROM environments ORDER BY updated_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var environments []*models.Environment
	for rows.Next() {
		var env models.Environment
		var variables string
		
		err := rows.Scan(
			&env.ID, &env.Name, &variables, &env.IsActive, &env.CreatedAt, &env.UpdatedAt)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(variables), &env.Variables)
		environments = append(environments, &env)
	}

	return environments, nil
}

// DeleteRequest deletes a request by ID
func (s *SQLiteStorage) DeleteRequest(id string) error {
	query := `DELETE FROM requests WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}

// DeleteCollection deletes a collection by ID
func (s *SQLiteStorage) DeleteCollection(id string) error {
	query := `DELETE FROM collections WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}

// DeleteEnvironment deletes an environment by ID
func (s *SQLiteStorage) DeleteEnvironment(id string) error {
	query := `DELETE FROM environments WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}
