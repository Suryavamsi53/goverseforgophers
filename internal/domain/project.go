package domain

import "context"

type ProjectDifficulty string

const (
	DifficultyBeginner     ProjectDifficulty = "Beginner"
	DifficultyIntermediate ProjectDifficulty = "Intermediate"
	DifficultyAdvanced     ProjectDifficulty = "Advanced"
)

type Project struct {
	ID           string            `json:"id"`
	Slug         string            `json:"slug"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	Difficulty   ProjectDifficulty `json:"difficulty"`
	Tags         []string          `json:"tags"`
	Icon         string            `json:"icon"`
	Color        string            `json:"color"` // Used for the UI glow (e.g., "blue", "purple", "emerald")
	Requirements []string          `json:"requirements"`
	StarterCode  string            `json:"starter_code"`
	TestFile     string            `json:"test_file"`
}

type ProjectRepository interface {
	GetAll(ctx context.Context) ([]Project, error)
	GetBySlug(ctx context.Context, slug string) (*Project, error)
}
