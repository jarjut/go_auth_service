package database

import (
	"auth-service/internal/domain"
	"auth-service/pkg/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect creates a database connection
func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	return db, nil
}

// AutoMigrate runs database migrations
// Deprecated: Use Atlas migrations instead (make atlas-generate && make atlas-apply)
// This is kept for backward compatibility and development
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Check if Atlas migrations are being used
	if err := MigrateWithAtlas(db); err != nil {
		return err
	}

	// Fallback to GORM AutoMigrate for development
	err := db.AutoMigrate(
		&domain.User{},
		&domain.RefreshToken{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	log.Println("💡 For production, use Atlas migrations: make atlas-help")
	return nil
}
