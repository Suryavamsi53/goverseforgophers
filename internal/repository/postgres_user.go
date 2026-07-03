package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryavamsivaggu/goverse/internal/domain"
)

type postgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash, user.Role).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return err
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE username = $1`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) CreateProfile(ctx context.Context, profile *domain.UserProfile) error {
	query := `
		INSERT INTO user_profiles (user_id, avatar_url, bio, github_handle, linkedin_handle, daily_streak, total_score)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query, profile.UserID, profile.AvatarURL, profile.Bio, profile.GithubHandle, profile.LinkedinHandle, profile.DailyStreak, profile.TotalScore)
	return err
}

func (r *postgresUserRepository) GetProfile(ctx context.Context, userID string) (*domain.UserProfile, error) {
	query := `
		SELECT user_id, avatar_url, bio, github_handle, linkedin_handle, daily_streak, total_score 
		FROM user_profiles WHERE user_id = $1
	`
	profile := &domain.UserProfile{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&profile.UserID, &profile.AvatarURL, &profile.Bio, &profile.GithubHandle, &profile.LinkedinHandle, &profile.DailyStreak, &profile.TotalScore,
	)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *postgresUserRepository) UpdateProfile(ctx context.Context, profile *domain.UserProfile) error {
	query := `
		UPDATE user_profiles 
		SET avatar_url = $1, bio = $2, github_handle = $3, linkedin_handle = $4, daily_streak = $5, total_score = $6
		WHERE user_id = $7
	`
	_, err := r.db.Exec(ctx, query, profile.AvatarURL, profile.Bio, profile.GithubHandle, profile.LinkedinHandle, profile.DailyStreak, profile.TotalScore, profile.UserID)
	return err
}

func (r *postgresUserRepository) GetSettings(ctx context.Context, userID string) (*domain.UserSettings, error) {
	query := `SELECT user_id, editor_settings, extensions FROM user_settings WHERE user_id = $1`
	settings := &domain.UserSettings{}
	err := r.db.QueryRow(ctx, query, userID).Scan(&settings.UserID, &settings.EditorSettings, &settings.Extensions)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *postgresUserRepository) UpdateSettings(ctx context.Context, settings *domain.UserSettings) error {
	query := `
		INSERT INTO user_settings (user_id, editor_settings, extensions) 
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) 
		DO UPDATE SET editor_settings = EXCLUDED.editor_settings, extensions = EXCLUDED.extensions, updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(ctx, query, settings.UserID, settings.EditorSettings, settings.Extensions)
	return err
}

func (r *postgresUserRepository) GetLeaderboard(ctx context.Context, limit int) ([]*domain.LeaderboardEntry, error) {
	query := `
		SELECT u.id, u.username, p.avatar_url, p.daily_streak, p.total_score 
		FROM users u
		JOIN user_profiles p ON u.id = p.user_id
		ORDER BY p.total_score DESC, p.daily_streak DESC
		LIMIT $1
	`
	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*domain.LeaderboardEntry
	for rows.Next() {
		e := &domain.LeaderboardEntry{}
		err := rows.Scan(&e.UserID, &e.Username, &e.AvatarURL, &e.DailyStreak, &e.TotalScore)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
