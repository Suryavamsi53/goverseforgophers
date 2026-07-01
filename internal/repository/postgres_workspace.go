package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryavamsivaggu/goverse/internal/domain"
)

type postgresWorkspaceRepository struct {
	db *pgxpool.Pool
}

func NewPostgresWorkspaceRepository(db *pgxpool.Pool) domain.WorkspaceRepository {
	return &postgresWorkspaceRepository{db: db}
}

func (r *postgresWorkspaceRepository) Save(ctx context.Context, ws *domain.Workspace) error {
	filesJSON, err := json.Marshal(ws.Files)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO user_workspaces (id, user_id, type, ref_id, files, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, type, ref_id)
		DO UPDATE SET files = EXCLUDED.files, updated_at = EXCLUDED.updated_at
	`

	if ws.ID == "" {
		ws.ID = uuid.New().String()
	}
	
	now := time.Now()
	if ws.CreatedAt.IsZero() {
		ws.CreatedAt = now
	}
	ws.UpdatedAt = now

	_, err = r.db.Exec(ctx, query, ws.ID, ws.UserID, ws.Type, ws.RefID, filesJSON, ws.CreatedAt, ws.UpdatedAt)
	return err
}

func (r *postgresWorkspaceRepository) Get(ctx context.Context, userID, wsType, refID string) (*domain.Workspace, error) {
	query := `
		SELECT id, user_id, type, ref_id, files, created_at, updated_at
		FROM user_workspaces
		WHERE user_id = $1 AND type = $2 AND ref_id = $3
	`

	ws := &domain.Workspace{}
	var filesJSON []byte

	err := r.db.QueryRow(ctx, query, userID, wsType, refID).Scan(
		&ws.ID,
		&ws.UserID,
		&ws.Type,
		&ws.RefID,
		&filesJSON,
		&ws.CreatedAt,
		&ws.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(filesJSON, &ws.Files); err != nil {
		return nil, err
	}

	return ws, nil
}
