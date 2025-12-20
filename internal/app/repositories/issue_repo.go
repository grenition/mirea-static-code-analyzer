package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"webapp/internal/app/models"
)

type IssueRepository struct {
	db *sql.DB
}

func NewIssueRepository(db *sql.DB) *IssueRepository {
	return &IssueRepository{db: db}
}

func (r *IssueRepository) CreateBatch(issues []*models.Issue) error {
	if len(issues) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`INSERT INTO issues (id, analysis_id, severity, rule_code, message, file_path, line)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, issue := range issues {
		_, err := stmt.Exec(
			issue.ID, issue.AnalysisID, issue.Severity, issue.RuleCode,
			issue.Message, issue.FilePath, issue.Line,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *IssueRepository) GetByAnalysisID(analysisID uuid.UUID) ([]*models.Issue, error) {
	rows, err := r.db.Query(
		`SELECT id, analysis_id, severity, rule_code, message, file_path, line
		 FROM issues WHERE analysis_id = $1 ORDER BY file_path, line`,
		analysisID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []*models.Issue
	for rows.Next() {
		var issue models.Issue
		var line sql.NullInt64

		if err := rows.Scan(
			&issue.ID, &issue.AnalysisID, &issue.Severity, &issue.RuleCode,
			&issue.Message, &issue.FilePath, &line,
		); err != nil {
			return nil, err
		}

		if line.Valid {
			lineNum := int(line.Int64)
			issue.Line = &lineNum
		}

		issues = append(issues, &issue)
	}
	return issues, rows.Err()
}

