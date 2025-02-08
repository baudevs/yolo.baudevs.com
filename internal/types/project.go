package types

import "time"

// Project represents a YOLO project
type Project struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Epics       []Epic    `json:"epics"`
}

// Epic represents a large body of work that can be broken down into features
type Epic struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Status       string    `json:"status"` // planned, in-progress, completed
	Features     []Feature `json:"features"`
	Dependencies []string  `json:"dependencies"` // IDs of other epics this depends on
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Feature represents a specific functionality within an epic
type Feature struct {
	ID           string    `json:"id"`
	EpicID       string    `json:"epic_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Status       string    `json:"status"` // planned, in-progress, completed
	Tasks        []Task    `json:"tasks"`
	Dependencies []string  `json:"dependencies"` // IDs of other features this depends on
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Task represents a specific unit of work within a feature
type Task struct {
	ID           string    `json:"id"`
	FeatureID    string    `json:"feature_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`       // planned, in-progress, completed
	Dependencies []string  `json:"dependencies"` // IDs of other tasks this depends on
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ProjectVersion represents the development stage of a project
type ProjectVersion string

const (
	VersionMVP      ProjectVersion = "mvp"
	VersionV1       ProjectVersion = "v1"
	VersionComplete ProjectVersion = "complete"
)
