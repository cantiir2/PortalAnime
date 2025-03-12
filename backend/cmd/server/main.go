package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/username/anime-streaming/internal/api/routes"
	"github.com/username/anime-streaming/internal/config"
	"github.com/username/anime-streaming/internal/db"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize database connection
	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := db.Migrate(database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Buat direktori media
	mediaPath := os.Getenv("MEDIA_PATH")
	if mediaPath == "" {
		mediaPath = "./media"
	}

	if err := os.MkdirAll(filepath.Join(mediaPath, "thumbnails", "content"), 0755); err != nil {
		log.Printf("Warning: Failed to create media directories: %v", err)
	}

	// Initialize router
	router := routes.SetupRouter(database, cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	log.Printf("Environment: %s", os.Getenv("ENV"))

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
