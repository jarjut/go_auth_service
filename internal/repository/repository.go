package repository

import (
	"auth-service/internal/domain"
	"context"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

// RefreshTokenRepository defines the interface for refresh token data access
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *domain.RefreshToken) error
	FindByToken(ctx context.Context, tokenString string) (*domain.RefreshToken, error)
	FindByUserID(ctx context.Context, userID string) ([]*domain.RefreshToken, error)
	Revoke(ctx context.Context, tokenString string) error
	RevokeAllByUserID(ctx context.Context, userID string) error
	DeleteExpired(ctx context.Context) error
}
