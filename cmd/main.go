package main

import (
	"auth-service/internal/delivery/http"
	"auth-service/internal/repository"
	"auth-service/internal/usecase"
	"auth-service/pkg/config"
	"auth-service/pkg/database"
	"auth-service/pkg/jwt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize JWT manager
	jwtManager, err := jwt.NewJWTManager(
		cfg.JWT.PrivateKeyPath,
		cfg.JWT.PublicKeyPath,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
	)
	if err != nil {
		log.Fatalf("Failed to initialize JWT manager: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, refreshTokenRepo, jwtManager)

	// Initialize dependency container
	container := http.NewContainer(authUseCase, jwtManager)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Setup routes
	http.SetupRoutes(app, container)

	// Start server
	log.Printf("Starting server on port %s...", cfg.Server.Port)
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
