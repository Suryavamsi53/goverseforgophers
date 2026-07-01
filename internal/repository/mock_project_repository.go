package repository

import (
	"context"
	"errors"

	"github.com/suryavamsivaggu/goverse/internal/domain"
)

type MockProjectRepository struct {
	projects []domain.Project
}

func NewMockProjectRepository() *MockProjectRepository {
	return &MockProjectRepository{
		projects: []domain.Project{
			{
				ID:          "proj-1",
				Slug:        "url-shortener",
				Title:       "URL Shortener API",
				Description: "Design and build a production-ready REST API for shortening URLs. Includes PostgreSQL integration, Redis caching, and rate limiting.",
				Difficulty:  domain.DifficultyBeginner,
				Tags:        []string{"REST API", "PostgreSQL", "Redis"},
				Icon:        "🔗",
				Color:       "emerald",
				Requirements: []string{
					"Create an endpoint to submit a long URL and receive a short code.",
					"Create an endpoint to redirect from a short code to the original URL.",
					"Use PostgreSQL to persist the URL mappings.",
					"Implement a Redis cache layer to speed up redirections.",
					"Add basic rate limiting middleware to prevent abuse.",
				},
				StarterCode: ``,
				TestFile: ``,
			},
			{
				ID:          "proj-2",
				Slug:        "concurrent-scraper",
				Title:       "Concurrent Web Scraper",
				Description: "Build a high-performance web scraper using goroutines, channels, and worker pools. Learn to manage rate limits and graceful shutdown.",
				Difficulty:  domain.DifficultyIntermediate,
				Tags:        []string{"Goroutines", "Channels", "HTTP Client"},
				Icon:        "🕷️",
				Color:       "blue",
				Requirements: []string{
					"Implement a worker pool pattern to fetch multiple URLs concurrently.",
					"Use channels to safely pass URLs to workers and collect results.",
					"Implement rate limiting to avoid overwhelming target servers.",
					"Add graceful shutdown using context.Context.",
					"Extract title and meta tags from the fetched HTML.",
				},
				StarterCode: ``,
				TestFile: ``,
			},
			{
				ID:          "proj-3",
				Slug:        "distributed-kv",
				Title:       "Distributed KV Store",
				Description: "Implement a highly available distributed Key-Value store. Master the Raft consensus algorithm, gRPC communication, and WAL logging.",
				Difficulty:  domain.DifficultyAdvanced,
				Tags:        []string{"gRPC", "Raft Consensus", "Mutexes"},
				Icon:        "🗄️",
				Color:       "purple",
				Requirements: []string{
					"Define a gRPC service for Put, Get, and Delete operations.",
					"Implement a thread-safe in-memory map using sync.RWMutex.",
					"Write an append-only Write-Ahead Log (WAL) to disk for durability.",
					"Implement basic leader election using a simplified Raft concept.",
					"Replicate write operations to follower nodes.",
				},
				StarterCode: ``,
				TestFile: ``,
			},
		},
	}
}

func (r *MockProjectRepository) GetAll(ctx context.Context) ([]domain.Project, error) {
	return r.projects, nil
}

func (r *MockProjectRepository) GetBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	for _, p := range r.projects {
		if p.Slug == slug {
			return &p, nil
		}
	}
	return nil, errors.New("project not found")
}
