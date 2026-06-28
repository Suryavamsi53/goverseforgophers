package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryavamsivaggu/goverse/internal/domain"
)

type postgresProblemRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProblemRepository(db *pgxpool.Pool) domain.ProblemRepository {
	return &postgresProblemRepository{db: db}
}

func (r *postgresProblemRepository) GetAll(ctx context.Context, category string) ([]domain.PracticeProblem, error) {
	// Not implemented for brevity
	return nil, nil
}

func (r *postgresProblemRepository) GetBySlug(ctx context.Context, slug string) (*domain.PracticeProblem, error) {
	var p domain.PracticeProblem
	var hintsBytes, testCasesBytes []byte

	err := r.db.QueryRow(ctx, "SELECT id, slug, title, category, difficulty, description, hints, starter_code, test_cases, created_at FROM practice_problems WHERE slug = $1", slug).Scan(
		&p.ID, &p.Slug, &p.Title, &p.Category, &p.Difficulty, &p.Description, &hintsBytes, &p.StarterCode, &testCasesBytes, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if len(hintsBytes) > 0 {
		_ = json.Unmarshal(hintsBytes, &p.Hints)
	}
	if len(testCasesBytes) > 0 {
		_ = json.Unmarshal(testCasesBytes, &p.TestCases)
	}

	return &p, nil
}
