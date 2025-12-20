package repositories

import (
	"database/sql"
	"time"
	"github.com/google/uuid"
	"webapp/internal/app/models"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(name string) (*models.Project, error) {
	id := uuid.New()
	now := time.Now()

	_, err := r.db.Exec(
		"INSERT INTO projects (id, name, created_at) VALUES ($1, $2, $3)",
		id, name, now,
	)
	if err != nil {
		return nil, err
	}

	return &models.Project{
		ID:        id,
		Name:      name,
		CreatedAt: now,
	}, nil
}

func (r *ProjectRepository) GetByID(id uuid.UUID) (*models.Project, error) {
	var p models.Project
	err := r.db.QueryRow(
		"SELECT id, name, created_at FROM projects WHERE id = $1",
		id,
	).Scan(&p.ID, &p.Name, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepository) ListAll() ([]*models.Project, error) {
	rows, err := r.db.Query("SELECT id, name, created_at FROM projects ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}
	return projects, rows.Err()
}

