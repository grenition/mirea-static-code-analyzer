package models

import (
	"time"
	"github.com/google/uuid"
)

type AnalysisFile struct {
	ID         uuid.UUID
	AnalysisID uuid.UUID
	FilePath   string
	Content    []byte
	CreatedAt  time.Time
}

