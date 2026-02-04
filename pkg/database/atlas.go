package database

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

// MigrateWithAtlas checks if Atlas migrations table exists
// This helps transition from GORM AutoMigrate to Atlas migrations
func MigrateWithAtlas(db *gorm.DB) error {
	// Check if atlas_schema_revisions table exists
	var exists bool
	err := db.Raw(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'atlas_schema_revisions'
		)
	`).Scan(&exists).Error

	if err != nil {
		return fmt.Errorf("failed to check atlas migrations table: %w", err)
	}

	if !exists {
		log.Println("⚠️  Atlas migrations not detected. Run 'make atlas-generate' to create initial migration.")
		log.Println("💡 Using GORM AutoMigrate as fallback for development...")
		return nil
	}

	log.Println("✅ Atlas migrations detected. Skipping GORM AutoMigrate.")
	log.Println("💡 Run 'make atlas-status' to check migration status.")
	return nil
}

// InitializeAtlas provides instructions for setting up Atlas
func InitializeAtlas() {
	separator := strings.Repeat("=", 70)
	log.Println("\n" + separator)
	log.Println("📦 Atlas Migration Setup")
	log.Println(separator)
	log.Println("To use Atlas for database migrations:")
	log.Println("")
	log.Println("1. Install Atlas:")
	log.Println("   make install-tools")
	log.Println("")
	log.Println("2. Generate initial migration:")
	log.Println("   make atlas-generate")
	log.Println("")
	log.Println("3. Apply migrations:")
	log.Println("   make atlas-apply")
	log.Println("")
	log.Println("4. Check status:")
	log.Println("   make atlas-status")
	log.Println("")
	log.Println("📚 See docs/MIGRATIONS.md for detailed guide")
	log.Println(separator + "\n")
}
