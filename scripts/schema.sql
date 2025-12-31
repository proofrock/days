-- Journal Application Database Schema
-- SQLite database for storing journal entries

-- Table: entries
-- Stores the header information for each journal entry
CREATE TABLE IF NOT EXISTS entries (
    date TEXT PRIMARY KEY,          -- Date of the journal entry (YYYY-MM-DD format)
    timestamp TEXT NOT NULL         -- Last update timestamp (ISO8601 format)
);

-- Table: details
-- Stores the field values for each journal entry
-- This structure allows easy addition of new fields in the future
CREATE TABLE IF NOT EXISTS details (
    date TEXT NOT NULL,             -- Reference to entries.date
    field_id TEXT NOT NULL,         -- Field identifier (POSITION_LON, POSITION_LAT, RATING, etc.)
    value TEXT,                     -- Field value (NULL if not filled)
    PRIMARY KEY (date, field_id),
    FOREIGN KEY (date) REFERENCES entries(date) ON DELETE CASCADE
);

-- Index for faster lookups by date
CREATE INDEX IF NOT EXISTS idx_details_date ON details(date);
