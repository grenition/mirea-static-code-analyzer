package models

import "time"

type Project struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	UserID    int       `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type File struct {
	ID        int       `json:"id" db:"id"`
	ProjectID int       `json:"project_id" db:"project_id"`
	Path      string    `json:"path" db:"path"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateProjectRequest struct {
	Name string `json:"name"`
}

type UpdateProjectRequest struct {
	Name string `json:"name"`
}

type CreateFileRequest struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type UpdateFileRequest struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type ProjectWithFiles struct {
	Project
	Files []File `json:"files"`
}

