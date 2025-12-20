package models

import (
	"time"
	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID
	Name      string
	UserID    uuid.UUID
	CreatedAt time.Time
}

