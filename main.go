package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/mano/days/internal/api"
	"github.com/mano/days/internal/database"
)

// Version information (set by build flags)
var Version = "dev"

//go:embed frontend/dist/*
var frontendFS embed.FS

func main() {
	// Print version
	log.Printf("Days Journal v%s starting...", Version)

	// Initialize database
	dbPath := getEnv("DB_PATH", "./data/journal.db")
	db, err := database.Initialize(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Set up API routes
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, db)

	// Serve embedded frontend
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Fatalf("Failed to access embedded frontend: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(distFS)))

	// Start server
	port := getEnv("PORT", "8080")
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
