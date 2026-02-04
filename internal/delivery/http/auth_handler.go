package http

import (
	"auth-service/internal/domain"
	"auth-service/internal/usecase"
	"auth-service/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authUseCase usecase.AuthUseCase
	jwtManager  *jwt.JWTManager
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUseCase usecase.AuthUseCase, jwtManager *jwt.JWTManager) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		jwtManager:  jwtManager,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email, password, and name
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecase.RegisterRequest true "Registration request"
// @Success 201 {object} usecase.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req usecase.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Basic validation
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email, password, and name are required",
		})
	}

	if len(req.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password must be at least 8 characters",
		})
	}

	resp, err := h.authUseCase.Register(c.Context(), req)
	if err != nil {
		if err == domain.ErrUserAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "user already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// Login handles user login
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecase.LoginRequest true "Login request"
// @Success 200 {object} usecase.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req usecase.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email and password are required",
		})
	}

	resp, err := h.authUseCase.Login(c.Context(), req)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid credentials",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to login",
		})
	}

	return c.JSON(resp)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Get a new access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecase.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} usecase.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req usecase.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "refresh token is required",
		})
	}

	resp, err := h.authUseCase.RefreshToken(c.Context(), req)
	if err != nil {
		if err == domain.ErrInvalidToken || err == domain.ErrRefreshTokenExpired || err == domain.ErrRefreshTokenRevoked {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to refresh token",
		})
	}

	return c.JSON(resp)
}

// Logout handles user logout (revokes refresh token)
// @Summary Logout user
// @Description Revoke a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecase.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req usecase.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "refresh token is required",
		})
	}

	err := h.authUseCase.Logout(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to logout",
		})
	}

	return c.JSON(fiber.Map{
		"message": "successfully logged out",
	})
}

// LogoutAll handles logout from all devices (revokes all refresh tokens for a user)
// @Summary Logout from all devices
// @Description Revoke all refresh tokens for the authenticated user
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/logout-all [post]
func (h *AuthHandler) LogoutAll(c *fiber.Ctx) error {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	err := h.authUseCase.LogoutAll(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to logout from all devices",
		})
	}

	return c.JSON(fiber.Map{
		"message": "successfully logged out from all devices",
	})
}

// GetProfile retrieves the authenticated user's profile
// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} usecase.UserResponse
// @Failure 401 {object} map[string]interface{}
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	user, err := h.authUseCase.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get user profile",
		})
	}

	return c.JSON(user)
}

// GetJWKS returns the JSON Web Key Set
// @Summary Get JWKS
// @Description Get the JSON Web Key Set for token validation
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /.well-known/jwks.json [get]
func (h *AuthHandler) GetJWKS(c *fiber.Ctx) error {
	jwks, err := h.jwtManager.GetJWKS()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get JWKS",
		})
	}

	return c.JSON(jwks)
}
