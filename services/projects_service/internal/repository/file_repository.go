package repository

import (
	"database/sql"
	"fmt"

	"projects_service/internal/models"
)

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(file *models.File) error {
	query := `INSERT INTO files (project_id, path, content) VALUES ($1, $2, $3) 
		ON CONFLICT (project_id, path) DO UPDATE SET content = $3, updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(query, file.ProjectID, file.Path, file.Content).Scan(&file.ID, &file.CreatedAt, &file.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	return nil
}

func (r *FileRepository) FindByID(id int) (*models.File, error) {
	file := &models.File{}
	query := `SELECT id, project_id, path, content, created_at, updated_at FROM files WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&file.ID, &file.ProjectID, &file.Path, &file.Content, &file.CreatedAt, &file.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("file not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find file: %w", err)
	}
	return file, nil
}

func (r *FileRepository) FindByProjectID(projectID int) ([]*models.File, error) {
	query := `SELECT id, project_id, path, content, created_at, updated_at FROM files WHERE project_id = $1 ORDER BY path`
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var files []*models.File
	for rows.Next() {
		file := &models.File{}
		if err := rows.Scan(&file.ID, &file.ProjectID, &file.Path, &file.Content, &file.CreatedAt, &file.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	return files, nil
}

func (r *FileRepository) Update(file *models.File) error {
	query := `UPDATE files SET path = $1, content = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING updated_at`
	err := r.db.QueryRow(query, file.Path, file.Content, file.ID).Scan(&file.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}
	return nil
}

func (r *FileRepository) Delete(id int) error {
	query := `DELETE FROM files WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("file not found")
	}
	return nil
}

func (r *FileRepository) DeleteByProjectID(projectID int) error {
	query := `DELETE FROM files WHERE project_id = $1`
	_, err := r.db.Exec(query, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete files: %w", err)
	}
	return nil
}

