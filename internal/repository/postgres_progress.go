package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryavamsivaggu/goverse/internal/domain"
)

type postgresProgressRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProgressRepository(db *pgxpool.Pool) domain.ProgressRepository {
	return &postgresProgressRepository{db: db}
}

func (r *postgresProgressRepository) MarkCompleted(ctx context.Context, userID, entityType, entityID string) error {
	now := time.Now()
	query := `
		INSERT INTO user_progress (user_id, entity_type, entity_id, status, completed_at)
		VALUES ($1, $2, $3, 'completed', $4)
		ON CONFLICT (user_id, entity_type, entity_id) 
		DO UPDATE SET status = 'completed', completed_at = $4
	`
	_, err := r.db.Exec(ctx, query, userID, entityType, entityID, now)
	return err
}

func (r *postgresProgressRepository) GetProgress(ctx context.Context, userID string) ([]domain.UserProgress, error) {
	rows, err := r.db.Query(ctx, "SELECT user_id, entity_type, entity_id, status, completed_at FROM user_progress WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progress []domain.UserProgress
	for rows.Next() {
		var p domain.UserProgress
		if err := rows.Scan(&p.UserID, &p.EntityType, &p.EntityID, &p.Status, &p.CompletedAt); err != nil {
			return nil, err
		}
		progress = append(progress, p)
	}
	return progress, nil
}
