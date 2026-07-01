package domain

import (
	"context"
	"time"
)

type Workspace struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Type      string            `json:"type"`   // e.g., "practice", "project"
	RefID     string            `json:"ref_id"` // e.g., "default", "concurrent-scraper"
	Files     map[string]string `json:"files"`  // filename -> content
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type WorkspaceRepository interface {
	Save(ctx context.Context, workspace *Workspace) error
	Get(ctx context.Context, userID, wsType, refID string) (*Workspace, error)
}
