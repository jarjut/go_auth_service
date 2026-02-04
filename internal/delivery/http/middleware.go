package http

import (
	"auth-service/internal/usecase"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT access token
func AuthMiddleware(authUseCase usecase.AuthUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		token := parts[1]

		// Validate token
		claims, err := authUseCase.ValidateAccessToken(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// Store user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}

// GetUserIDFromContext retrieves the user ID from the context
func GetUserIDFromContext(c *fiber.Ctx) (string, bool) {
	userID, ok := c.Locals("userID").(string)
	return userID, ok
}

// GetEmailFromContext retrieves the email from the context
func GetEmailFromContext(c *fiber.Ctx) (string, bool) {
	email, ok := c.Locals("email").(string)
	return email, ok
}
