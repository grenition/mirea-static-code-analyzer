package models

import (
	"github.com/google/uuid"
)

type Issue struct {
	ID         uuid.UUID
	AnalysisID uuid.UUID
	Severity   string // "error", "warn", "info"
	RuleCode   string
	Message    string
	FilePath   string
	Line       *int
}

