package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// SetupRoutes sets up all HTTP routes
func SetupRoutes(app *fiber.App, container *Container) {
	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Initialize handlers
	authHandler := NewAuthHandler(container.AuthUseCase, container.JWTManager)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "auth-service",
		})
	})

	// JWKS endpoint (public key for token verification)
	app.Get("/.well-known/jwks.json", authHandler.GetJWKS)

	// Auth routes
	auth := app.Group("/auth")
	{
		// Public routes
		auth.Post("/register", authHandler.Register)
		auth.Post("/login", authHandler.Login)
		auth.Post("/refresh", authHandler.RefreshToken)
		auth.Post("/logout", authHandler.Logout)

		// Protected routes (require authentication)
		protected := auth.Group("", AuthMiddleware(container.AuthUseCase))
		protected.Get("/profile", authHandler.GetProfile)
		protected.Post("/logout-all", authHandler.LogoutAll)
	}
}
