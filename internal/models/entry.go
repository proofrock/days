package models

import (
	"database/sql"
	"time"
)

// Entry represents a journal entry with all its fields
type Entry struct {
	Date      string            `json:"date"`      // Date in YYYY-MM-DD format
	Timestamp string            `json:"timestamp"` // ISO8601 timestamp of last update
	Fields    map[string]string `json:"fields"`    // All field values (field_id -> value)
}

// Field IDs as constants
const (
	FieldPositionLon  = "POSITION_LON"
	FieldPositionLat  = "POSITION_LAT"
	FieldPositionName = "POSITION_NAME"
	FieldRating       = "RATING"
	FieldGeneral      = "GENERAL"
	FieldMood         = "MOOD"
	FieldMoodText     = "MOOD_TXT"
	FieldLunch        = "LUNCH"
	FieldDinner       = "DINNER"
	FieldTV           = "TV"
	FieldSleep        = "SLEEP"
	FieldSleepText    = "SLEEP_TXT"
)

// AllFieldIDs returns all valid field IDs in order
func AllFieldIDs() []string {
	return []string{
		FieldPositionLon,
		FieldPositionLat,
		FieldPositionName,
		FieldRating,
		FieldGeneral,
		FieldMood,
		FieldMoodText,
		FieldLunch,
		FieldDinner,
		FieldTV,
		FieldSleep,
		FieldSleepText,
	}
}

// GetEntry retrieves a single entry by date
func GetEntry(db *sql.DB, date string) (*Entry, error) {
	// Check if entry exists
	var timestamp string
	err := db.QueryRow("SELECT timestamp FROM entries WHERE date = ?", date).Scan(&timestamp)
	if err == sql.ErrNoRows {
		return nil, nil // Entry doesn't exist
	}
	if err != nil {
		return nil, err
	}

	// Get all details
	rows, err := db.Query("SELECT field_id, value FROM details WHERE date = ?", date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fields := make(map[string]string)
	for rows.Next() {
		var fieldID string
		var value sql.NullString
		if err := rows.Scan(&fieldID, &value); err != nil {
			return nil, err
		}
		// Store empty string if value is NULL
		if value.Valid {
			fields[fieldID] = value.String
		} else {
			fields[fieldID] = ""
		}
	}

	return &Entry{
		Date:      date,
		Timestamp: timestamp,
		Fields:    fields,
	}, nil
}

// GetEntriesByMonth retrieves all entry dates for a given year and month
func GetEntriesByMonth(db *sql.DB, year, month int) ([]string, error) {
	// SQLite date format: YYYY-MM-DD
	pattern := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Format("2006-01")

	rows, err := db.Query("SELECT date FROM entries WHERE date LIKE ? ORDER BY date", pattern+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}

	return dates, nil
}

// SaveEntry creates or updates an entry
func SaveEntry(db *sql.DB, entry *Entry) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update timestamp to current time
	entry.Timestamp = time.Now().UTC().Format(time.RFC3339)

	// Insert or update entry header
	_, err = tx.Exec(`
		INSERT INTO entries (date, timestamp) VALUES (?, ?)
		ON CONFLICT(date) DO UPDATE SET timestamp = ?
	`, entry.Date, entry.Timestamp, entry.Timestamp)
	if err != nil {
		return err
	}

	// Delete existing details for this date
	_, err = tx.Exec("DELETE FROM details WHERE date = ?", entry.Date)
	if err != nil {
		return err
	}

	// Insert all fields (including NULL values)
	for _, fieldID := range AllFieldIDs() {
		value := entry.Fields[fieldID]
		// Store empty strings as NULL
		var valuePtr *string
		if value != "" {
			valuePtr = &value
		}
		_, err = tx.Exec("INSERT INTO details (date, field_id, value) VALUES (?, ?, ?)",
			entry.Date, fieldID, valuePtr)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeleteEntry removes an entry and all its details
func DeleteEntry(db *sql.DB, date string) error {
	// Foreign key cascade will delete details automatically
	_, err := db.Exec("DELETE FROM entries WHERE date = ?", date)
	return err
}
