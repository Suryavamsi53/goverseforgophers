package domain

import (
	"context"
	"time"
)

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserProfile struct {
	UserID         string `json:"user_id"`
	AvatarURL      string `json:"avatar_url"`
	Bio            string `json:"bio"`
	GithubHandle   string `json:"github_handle"`
	LinkedinHandle string `json:"linkedin_handle"`
	DailyStreak    int    `json:"daily_streak"`
	TotalScore     int    `json:"total_score"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	
	CreateProfile(ctx context.Context, profile *UserProfile) error
	GetProfile(ctx context.Context, userID string) (*UserProfile, error)
	UpdateProfile(ctx context.Context, profile *UserProfile) error
}

type AuthUseCase interface {
	Register(ctx context.Context, username, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (string, error)
}
