package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mano/days/internal/models"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	// API endpoints
	mux.HandleFunc("/api/entries/", makeEntryHandler(db))
	mux.HandleFunc("/api/entries/month/", makeMonthHandler(db))
}

// makeEntryHandler creates a handler for individual entry operations
// GET /api/entries/{date} - Get entry
// POST /api/entries/{date} - Create/update entry
// DELETE /api/entries/{date} - Delete entry
func makeEntryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract date from path
		date := r.URL.Path[len("/api/entries/"):]
		if date == "" {
			http.Error(w, "Date required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetEntry(w, r, db, date)
		case http.MethodPost:
			handleSaveEntry(w, r, db, date)
		case http.MethodDelete:
			handleDeleteEntry(w, r, db, date)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// makeMonthHandler creates a handler for getting entries by month
// GET /api/entries/month/{year}/{month} - Get all entry dates for a month
func makeMonthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract year and month from path
		path := r.URL.Path[len("/api/entries/month/"):]
		var year, month int
		if _, err := fmt.Sscanf(path, "%d/%d", &year, &month); err != nil {
			http.Error(w, "Invalid year/month format", http.StatusBadRequest)
			return
		}

		dates, err := models.GetEntriesByMonth(db, year, month)
		if err != nil {
			log.Printf("Error getting entries by month: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		respondJSON(w, dates)
	}
}

// handleGetEntry retrieves a single entry
func handleGetEntry(w http.ResponseWriter, r *http.Request, db *sql.DB, date string) {
	entry, err := models.GetEntry(db, date)
	if err != nil {
		log.Printf("Error getting entry: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if entry == nil {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	respondJSON(w, entry)
}

// handleSaveEntry creates or updates an entry
func handleSaveEntry(w http.ResponseWriter, r *http.Request, db *sql.DB, date string) {
	var entry models.Entry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Ensure date matches URL
	entry.Date = date

	if err := models.SaveEntry(db, &entry); err != nil {
		log.Printf("Error saving entry: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, entry)
}

// handleDeleteEntry removes an entry
func handleDeleteEntry(w http.ResponseWriter, r *http.Request, db *sql.DB, date string) {
	if err := models.DeleteEntry(db, date); err != nil {
		log.Printf("Error deleting entry: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}
