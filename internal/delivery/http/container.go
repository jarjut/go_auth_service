package http

import (
	"auth-service/internal/usecase"
	"auth-service/pkg/jwt"
)

// Container holds all dependencies for HTTP handlers
type Container struct {
	// Use cases
	AuthUseCase usecase.AuthUseCase
	// Add more use cases here as your application grows
	// UserUseCase usecase.UserUseCase
	// ProductUseCase usecase.ProductUseCase
	// etc.

	// Utilities
	JWTManager *jwt.JWTManager
	// Add more utilities here
	// EmailService *email.Service
	// StorageService *storage.Service
	// etc.
}

// NewContainer creates a new dependency container
func NewContainer(
	authUseCase usecase.AuthUseCase,
	jwtManager *jwt.JWTManager,
) *Container {
	return &Container{
		AuthUseCase: authUseCase,
		JWTManager:  jwtManager,
	}
}
