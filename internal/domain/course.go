package domain

import (
	"context"
	"time"
)

type Course struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Difficulty  string    `json:"difficulty"`
	CreatedAt   time.Time `json:"created_at"`
}

type Lesson struct {
	ID         string    `json:"id"`
	CourseID   string    `json:"course_id"`
	Slug       string    `json:"slug"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	OrderIndex int       `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
}

type CourseRepository interface {
	GetAll(ctx context.Context) ([]Course, error)
	GetBySlug(ctx context.Context, slug string) (*Course, error)
	GetLessonsByCourseID(ctx context.Context, courseID string) ([]Lesson, error)
	GetLessonBySlug(ctx context.Context, slug string) (*Lesson, error)
}
