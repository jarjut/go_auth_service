package usecase

import (
	"auth-service/internal/domain"
	"auth-service/internal/repository"
	"auth-service/pkg/jwt"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// AuthUseCase defines the interface for authentication use cases
type AuthUseCase interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (*AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID string) error
	ValidateAccessToken(ctx context.Context, token string) (*jwt.Claims, error)
	GetUserByID(ctx context.Context, userID string) (*UserResponse, error)
}

type authUseCase struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtManager       *jwt.JWTManager
}

// NewAuthUseCase creates a new auth use case
func NewAuthUseCase(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtManager *jwt.JWTManager,
) AuthUseCase {
	return &authUseCase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
	}
}

func (uc *authUseCase) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && err != domain.ErrUserNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	return uc.generateTokens(ctx, user)
}

func (uc *authUseCase) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate tokens
	return uc.generateTokens(ctx, user)
}

func (uc *authUseCase) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*AuthResponse, error) {
	// Validate refresh token from database
	refreshToken, err := uc.refreshTokenRepo.FindByToken(ctx, req.RefreshToken)
	if err != nil {
		if err == domain.ErrRefreshTokenNotFound {
			return nil, domain.ErrInvalidToken
		}
		return nil, err
	}

	// Check if token is valid
	if !refreshToken.IsValid() {
		if refreshToken.IsRevoked {
			return nil, domain.ErrRefreshTokenRevoked
		}
		return nil, domain.ErrRefreshTokenExpired
	}

	// Get user (refresh token already contains user ID)
	user, err := uc.userRepo.FindByID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	// Revoke old refresh token
	if err := uc.refreshTokenRepo.Revoke(ctx, req.RefreshToken); err != nil {
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	// Generate new tokens
	return uc.generateTokens(ctx, user)
}

func (uc *authUseCase) Logout(ctx context.Context, refreshToken string) error {
	return uc.refreshTokenRepo.Revoke(ctx, refreshToken)
}

func (uc *authUseCase) LogoutAll(ctx context.Context, userID string) error {
	return uc.refreshTokenRepo.RevokeAllByUserID(ctx, userID)
}

func (uc *authUseCase) ValidateAccessToken(ctx context.Context, token string) (*jwt.Claims, error) {
	claims, err := uc.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	return claims, nil
}

func (uc *authUseCase) GetUserByID(ctx context.Context, userID string) (*UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

// generateTokens generates access and refresh tokens for a user
func (uc *authUseCase) generateTokens(ctx context.Context, user *domain.User) (*AuthResponse, error) {
	// Generate access token
	accessToken, err := uc.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshTokenString, expiresAt, err := uc.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Save refresh token to database
	refreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: expiresAt,
		IsRevoked: false,
	}

	if err := uc.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int(uc.jwtManager.GetAccessTokenDuration().Seconds()),
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}
