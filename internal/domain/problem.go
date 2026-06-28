package domain

import (
	"context"
	"time"
)

type TestCase struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
}

type PracticeProblem struct {
	ID          string     `json:"id"`
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Category    string     `json:"category"`
	Difficulty  string     `json:"difficulty"`
	Description string     `json:"description"`
	Hints       []string   `json:"hints"`
	StarterCode string     `json:"starter_code"`
	TestCases   []TestCase `json:"test_cases"`
	CreatedAt   time.Time  `json:"created_at"`
}

type ProblemRepository interface {
	GetAll(ctx context.Context, category string) ([]PracticeProblem, error)
	GetBySlug(ctx context.Context, slug string) (*PracticeProblem, error)
}
