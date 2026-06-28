package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryavamsivaggu/goverse/internal/domain"
)

type postgresCourseRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCourseRepository(db *pgxpool.Pool) domain.CourseRepository {
	return &postgresCourseRepository{db: db}
}

func (r *postgresCourseRepository) GetAll(ctx context.Context) ([]domain.Course, error) {
	// Not implemented for brevity
	return nil, nil
}

func (r *postgresCourseRepository) GetBySlug(ctx context.Context, slug string) (*domain.Course, error) {
	var c domain.Course
	err := r.db.QueryRow(ctx, "SELECT id, slug, title, description, difficulty, created_at FROM courses WHERE slug = $1", slug).Scan(
		&c.ID, &c.Slug, &c.Title, &c.Description, &c.Difficulty, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *postgresCourseRepository) GetLessonsByCourseID(ctx context.Context, courseID string) ([]domain.Lesson, error) {
	rows, err := r.db.Query(ctx, "SELECT id, course_id, slug, title, content, order_index, created_at FROM lessons WHERE course_id = $1 ORDER BY order_index ASC", courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []domain.Lesson
	for rows.Next() {
		var l domain.Lesson
		if err := rows.Scan(&l.ID, &l.CourseID, &l.Slug, &l.Title, &l.Content, &l.OrderIndex, &l.CreatedAt); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}
	return lessons, nil
}

func (r *postgresCourseRepository) GetLessonBySlug(ctx context.Context, slug string) (*domain.Lesson, error) {
	var l domain.Lesson
	err := r.db.QueryRow(ctx, "SELECT id, course_id, slug, title, content, order_index, created_at FROM lessons WHERE slug = $1", slug).Scan(
		&l.ID, &l.CourseID, &l.Slug, &l.Title, &l.Content, &l.OrderIndex, &l.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &l, nil
}
