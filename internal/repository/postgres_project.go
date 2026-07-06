package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryavamsivaggu/goverse/internal/domain"
)

type PostgresProjectRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProjectRepository(db *pgxpool.Pool) *PostgresProjectRepository {
	return &PostgresProjectRepository{db: db}
}

func (r *PostgresProjectRepository) GetAll(ctx context.Context) ([]domain.Project, error) {
	query := `
		SELECT id, slug, title, description, scenario, difficulty, icon, color, 
		       tags, requirements, tips, starter_code, test_file
		FROM projects
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		var scenario *string
		var tagsJSON, reqsJSON, tipsJSON []byte

		err := rows.Scan(
			&p.ID, &p.Slug, &p.Title, &p.Description, &scenario, &p.Difficulty, &p.Icon, &p.Color,
			&tagsJSON, &reqsJSON, &tipsJSON, &p.StarterCode, &p.TestFile,
		)
		if err != nil {
			return nil, err
		}

		if scenario != nil {
			p.Scenario = *scenario
		}

		if len(tagsJSON) > 0 {
			json.Unmarshal(tagsJSON, &p.Tags)
		}
		if len(reqsJSON) > 0 {
			json.Unmarshal(reqsJSON, &p.Requirements)
		}
		if len(tipsJSON) > 0 {
			json.Unmarshal(tipsJSON, &p.Tips)
		}

		projects = append(projects, p)
	}

	return projects, rows.Err()
}

func (r *PostgresProjectRepository) GetBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	query := `
		SELECT id, slug, title, description, scenario, difficulty, icon, color, 
		       tags, requirements, tips, starter_code, test_file
		FROM projects
		WHERE slug = $1
	`
	row := r.db.QueryRow(ctx, query, slug)

	var p domain.Project
	var scenario *string
	var tagsJSON, reqsJSON, tipsJSON []byte

	err := row.Scan(
		&p.ID, &p.Slug, &p.Title, &p.Description, &scenario, &p.Difficulty, &p.Icon, &p.Color,
		&tagsJSON, &reqsJSON, &tipsJSON, &p.StarterCode, &p.TestFile,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if scenario != nil {
		p.Scenario = *scenario
	}

	if len(tagsJSON) > 0 {
		json.Unmarshal(tagsJSON, &p.Tags)
	}
	if len(reqsJSON) > 0 {
		json.Unmarshal(reqsJSON, &p.Requirements)
	}
	if len(tipsJSON) > 0 {
		json.Unmarshal(tipsJSON, &p.Tips)
	}

	return &p, nil
}
