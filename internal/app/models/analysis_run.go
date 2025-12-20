package models

import (
	"time"
	"encoding/json"
	"github.com/google/uuid"
)

type AnalysisRun struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	UserID      uuid.UUID
	Status      string
	CreatedAt   time.Time
	StartedAt   *time.Time
	FinishedAt  *time.Time
	InputType   string // "zip" or "snippet"
	InputMeta   json.RawMessage
	SummaryJSON json.RawMessage
}

type AnalysisSummary struct {
	ErrorCount int `json:"error_count"`
	WarnCount  int `json:"warn_count"`
	InfoCount  int `json:"info_count"`
}

