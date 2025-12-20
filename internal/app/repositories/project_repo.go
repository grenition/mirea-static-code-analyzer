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

func (r *ProjectRepository) Create(name string, userID uuid.UUID) (*models.Project, error) {
	id := uuid.New()
	now := time.Now()

	_, err := r.db.Exec(
		"INSERT INTO projects (id, name, user_id, created_at) VALUES ($1, $2, $3, $4)",
		id, name, userID, now,
	)
	if err != nil {
		return nil, err
	}

	return &models.Project{
		ID:        id,
		Name:      name,
		UserID:    userID,
		CreatedAt: now,
	}, nil
}

func (r *ProjectRepository) GetByID(id uuid.UUID) (*models.Project, error) {
	var p models.Project
	err := r.db.QueryRow(
		"SELECT id, name, user_id, created_at FROM projects WHERE id = $1",
		id,
	).Scan(&p.ID, &p.Name, &p.UserID, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepository) ListByUserID(userID uuid.UUID) ([]*models.Project, error) {
	rows, err := r.db.Query(
		"SELECT id, name, created_at FROM projects WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
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
		p.UserID = userID
		projects = append(projects, &p)
	}
	return projects, rows.Err()
}

func (r *ProjectRepository) ListAll() ([]*models.Project, error) {
	rows, err := r.db.Query("SELECT id, name, user_id, created_at FROM projects ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.UserID, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}
	return projects, rows.Err()
}

