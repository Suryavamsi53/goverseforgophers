package domain

import (
	"context"
	"time"
)

type UserProgress struct {
	UserID      string     `json:"user_id"`
	EntityType  string     `json:"entity_type"`
	EntityID    string     `json:"entity_id"`
	Status      string     `json:"status"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type ProgressRepository interface {
	MarkCompleted(ctx context.Context, userID, entityType, entityID string) error
	GetProgress(ctx context.Context, userID string) ([]UserProgress, error)
}
