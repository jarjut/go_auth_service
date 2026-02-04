package repository

import (
	"auth-service/internal/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	var refreshToken domain.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrRefreshTokenNotFound
		}
		return nil, err
	}
	return &refreshToken, nil
}

func (r *refreshTokenRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.RefreshToken, error) {
	var tokens []*domain.RefreshToken
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_revoked = ?", userID, false).Find(&tokens).Error
	return tokens, err
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, tokenString string) error {
	return r.db.WithContext(ctx).Model(&domain.RefreshToken{}).
		Where("token = ?", tokenString).
		Update("is_revoked", true).Error
}

func (r *refreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Model(&domain.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true).Error
}

func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&domain.RefreshToken{}).Error
}
