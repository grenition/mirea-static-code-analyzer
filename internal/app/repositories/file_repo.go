package repositories

import (
	"database/sql"
	"time"
	"github.com/google/uuid"
	"webapp/internal/app/models"
)

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) CreateBatch(analysisID uuid.UUID, files []models.AnalysisFile) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`INSERT INTO analysis_files (id, analysis_id, file_path, content, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()
	for _, file := range files {
		id := uuid.New()
		_, err := stmt.Exec(id, analysisID, file.FilePath, file.Content, now)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *FileRepository) GetByAnalysisID(analysisID uuid.UUID) ([]*models.AnalysisFile, error) {
	rows, err := r.db.Query(
		`SELECT id, analysis_id, file_path, content, created_at
		 FROM analysis_files WHERE analysis_id = $1 ORDER BY file_path`,
		analysisID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*models.AnalysisFile
	for rows.Next() {
		var f models.AnalysisFile
		if err := rows.Scan(
			&f.ID, &f.AnalysisID, &f.FilePath, &f.Content, &f.CreatedAt,
		); err != nil {
			return nil, err
		}
		files = append(files, &f)
	}
	return files, rows.Err()
}

func (r *FileRepository) GetByID(id uuid.UUID) (*models.AnalysisFile, error) {
	var f models.AnalysisFile
	err := r.db.QueryRow(
		`SELECT id, analysis_id, file_path, content, created_at
		 FROM analysis_files WHERE id = $1`,
		id,
	).Scan(
		&f.ID, &f.AnalysisID, &f.FilePath, &f.Content, &f.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FileRepository) GetByAnalysisIDAndPath(analysisID uuid.UUID, filePath string) (*models.AnalysisFile, error) {
	var f models.AnalysisFile
	err := r.db.QueryRow(
		`SELECT id, analysis_id, file_path, content, created_at
		 FROM analysis_files WHERE analysis_id = $1 AND file_path = $2`,
		analysisID, filePath,
	).Scan(
		&f.ID, &f.AnalysisID, &f.FilePath, &f.Content, &f.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FileRepository) DeleteByAnalysisID(analysisID uuid.UUID) error {
	_, err := r.db.Exec(
		"DELETE FROM analysis_files WHERE analysis_id = $1",
		analysisID,
	)
	return err
}

