package models

import (
	"time"
)

// Collection represents a collection of requests
type Collection struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Requests    []Request `json:"requests"`
	Folders     []Folder  `json:"folders"`
	Variables   map[string]string `json:"variables"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Folder represents a folder within a collection
type Folder struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Requests []Request `json:"requests"`
	Folders  []Folder  `json:"folders"`
}

// CollectionSummary represents a summary of a collection
type CollectionSummary struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	RequestCount int      `json:"request_count"`
	FolderCount  int      `json:"folder_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
