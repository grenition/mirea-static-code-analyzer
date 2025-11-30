package repository

import (
	"database/sql"
	"fmt"

	"projects_service/internal/models"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *models.Project) error {
	query := `INSERT INTO projects (name, user_id) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(query, project.Name, project.UserID).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}

func (r *ProjectRepository) FindByID(id int) (*models.Project, error) {
	project := &models.Project{}
	query := `SELECT id, name, user_id, created_at, updated_at FROM projects WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&project.ID, &project.Name, &project.UserID, &project.CreatedAt, &project.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find project: %w", err)
	}
	return project, nil
}

func (r *ProjectRepository) FindByUserID(userID int) ([]*models.Project, error) {
	query := `SELECT id, name, user_id, created_at, updated_at FROM projects WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		if err := rows.Scan(&project.ID, &project.Name, &project.UserID, &project.CreatedAt, &project.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (r *ProjectRepository) Update(project *models.Project) error {
	query := `UPDATE projects SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 RETURNING updated_at`
	err := r.db.QueryRow(query, project.Name, project.ID).Scan(&project.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}
	return nil
}

func (r *ProjectRepository) Delete(id int) error {
	query := `DELETE FROM projects WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}
	return nil
}

