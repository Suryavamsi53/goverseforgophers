package usecase

import (
	"context"
	"errors"

	"github.com/suryavamsivaggu/goverse/internal/domain"
	"github.com/suryavamsivaggu/goverse/pkg/auth"
)

type authUseCase struct {
	userRepo   domain.UserRepository
	jwtManager *auth.JWTManager
}

func NewAuthUseCase(userRepo domain.UserRepository, jwtManager *auth.JWTManager) domain.AuthUseCase {
	return &authUseCase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (u *authUseCase) Register(ctx context.Context, username, email, password string) (*domain.User, error) {
	// Check if user exists
	existingUser, _ := u.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("email already in use")
	}

	existingUser, _ = u.userRepo.GetByUsername(ctx, username)
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         "user",
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Create empty profile
	profile := &domain.UserProfile{
		UserID: user.ID,
	}
	_ = u.userRepo.CreateProfile(ctx, profile)

	return user, nil
}

func (u *authUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = auth.VerifyPassword(user.PasswordHash, password)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := u.jwtManager.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
