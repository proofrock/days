package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Initialize creates and initializes the SQLite database
func Initialize(dbPath string) (*sql.DB, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, err
	}

	// Create entries table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS entries (
			date TEXT PRIMARY KEY,
			timestamp TEXT NOT NULL
		)
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	// Create details table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS details (
			date TEXT NOT NULL,
			field_id TEXT NOT NULL,
			value TEXT,
			PRIMARY KEY (date, field_id),
			FOREIGN KEY (date) REFERENCES entries(date) ON DELETE CASCADE
		)
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	// Create index
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_details_date ON details(date)
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
