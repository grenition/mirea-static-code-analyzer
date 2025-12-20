package repositories

import (
	"database/sql"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"webapp/internal/app/models"
)

type AnalysisRepository struct {
	db *sql.DB
}

func NewAnalysisRepository(db *sql.DB) *AnalysisRepository {
	return &AnalysisRepository{db: db}
}

func (r *AnalysisRepository) Create(projectID, userID uuid.UUID, inputType string, inputMeta json.RawMessage) (*models.AnalysisRun, error) {
	id := uuid.New()
	now := time.Now()

	_, err := r.db.Exec(
		`INSERT INTO analysis_runs (id, project_id, user_id, status, created_at, input_type, input_meta)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		id, projectID, userID, "Created", now, inputType, inputMeta,
	)
	if err != nil {
		return nil, err
	}

	return &models.AnalysisRun{
		ID:        id,
		ProjectID: projectID,
		UserID:    userID,
		Status:    "Created",
		CreatedAt: now,
		InputType: inputType,
		InputMeta: inputMeta,
	}, nil
}

func (r *AnalysisRepository) UpdateStatus(id uuid.UUID, status string) error {
	now := time.Now()
	var query string
	var args []interface{}

	if status == "Running" {
		query = "UPDATE analysis_runs SET status = $1, started_at = $2 WHERE id = $3"
		args = []interface{}{status, now, id}
	} else if status == "Done" || status == "Failed" {
		query = "UPDATE analysis_runs SET status = $1, finished_at = $2 WHERE id = $3"
		args = []interface{}{status, now, id}
	} else {
		query = "UPDATE analysis_runs SET status = $1 WHERE id = $2"
		args = []interface{}{status, id}
	}

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *AnalysisRepository) UpdateSummary(id uuid.UUID, summary json.RawMessage) error {
	_, err := r.db.Exec(
		"UPDATE analysis_runs SET summary_json = $1 WHERE id = $2",
		summary, id,
	)
	return err
}

func (r *AnalysisRepository) GetByID(id uuid.UUID) (*models.AnalysisRun, error) {
	var a models.AnalysisRun
	var startedAt, finishedAt sql.NullTime
	var inputMeta, summaryJSON sql.NullString

	err := r.db.QueryRow(
		`SELECT id, project_id, user_id, status, created_at, started_at, finished_at, 
		 input_type, input_meta, summary_json
		 FROM analysis_runs WHERE id = $1`,
		id,
	).Scan(
		&a.ID, &a.ProjectID, &a.UserID, &a.Status, &a.CreatedAt,
		&startedAt, &finishedAt, &a.InputType, &inputMeta, &summaryJSON,
	)
	if err != nil {
		return nil, err
	}

	if startedAt.Valid {
		a.StartedAt = &startedAt.Time
	}
	if finishedAt.Valid {
		a.FinishedAt = &finishedAt.Time
	}
	if inputMeta.Valid {
		a.InputMeta = json.RawMessage(inputMeta.String)
	}
	if summaryJSON.Valid {
		a.SummaryJSON = json.RawMessage(summaryJSON.String)
	}

	return &a, nil
}

func (r *AnalysisRepository) ListAll() ([]*models.AnalysisRun, error) {
	rows, err := r.db.Query(
		`SELECT id, project_id, user_id, status, created_at, started_at, finished_at,
		 input_type, input_meta, summary_json
		 FROM analysis_runs ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []*models.AnalysisRun
	for rows.Next() {
		var a models.AnalysisRun
		var startedAt, finishedAt sql.NullTime
		var inputMeta, summaryJSON sql.NullString

		if err := rows.Scan(
			&a.ID, &a.ProjectID, &a.UserID, &a.Status, &a.CreatedAt,
			&startedAt, &finishedAt, &a.InputType, &inputMeta, &summaryJSON,
		); err != nil {
			return nil, err
		}

		if startedAt.Valid {
			a.StartedAt = &startedAt.Time
		}
		if finishedAt.Valid {
			a.FinishedAt = &finishedAt.Time
		}
		if inputMeta.Valid {
			a.InputMeta = json.RawMessage(inputMeta.String)
		}
		if summaryJSON.Valid {
			a.SummaryJSON = json.RawMessage(summaryJSON.String)
		}

		analyses = append(analyses, &a)
	}
	return analyses, rows.Err()
}

func (r *AnalysisRepository) ListByUserID(userID uuid.UUID) ([]*models.AnalysisRun, error) {
	rows, err := r.db.Query(
		`SELECT id, project_id, user_id, status, created_at, started_at, finished_at,
		 input_type, input_meta, summary_json
		 FROM analysis_runs WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []*models.AnalysisRun
	for rows.Next() {
		var a models.AnalysisRun
		var startedAt, finishedAt sql.NullTime
		var inputMeta, summaryJSON sql.NullString

		if err := rows.Scan(
			&a.ID, &a.ProjectID, &a.UserID, &a.Status, &a.CreatedAt,
			&startedAt, &finishedAt, &a.InputType, &inputMeta, &summaryJSON,
		); err != nil {
			return nil, err
		}

		if startedAt.Valid {
			a.StartedAt = &startedAt.Time
		}
		if finishedAt.Valid {
			a.FinishedAt = &finishedAt.Time
		}
		if inputMeta.Valid {
			a.InputMeta = json.RawMessage(inputMeta.String)
		}
		if summaryJSON.Valid {
			a.SummaryJSON = json.RawMessage(summaryJSON.String)
		}

		analyses = append(analyses, &a)
	}
	return analyses, rows.Err()
}

func (r *AnalysisRepository) Update(id uuid.UUID, inputMeta json.RawMessage, summary json.RawMessage) error {
	_, err := r.db.Exec(
		`UPDATE analysis_runs SET input_meta = $1, summary_json = $2 WHERE id = $3`,
		inputMeta, summary, id,
	)
	return err
}

func (r *AnalysisRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec(
		"DELETE FROM analysis_runs WHERE id = $1",
		id,
	)
	return err
}

