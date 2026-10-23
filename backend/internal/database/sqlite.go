package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Open() (*sql.DB, error) {
	dbPath := os.Getenv("JOBTRACK_DB")
	if dbPath == "" {
		dbPath = "jobtrack.db"
	}

	if dir := filepath.Dir(dbPath); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("create database directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS jobs (
			id TEXT PRIMARY KEY,
			company TEXT NOT NULL,
			job_title TEXT NOT NULL,
			job_url TEXT,
			date_applied TEXT,
			deadline TEXT,
			status TEXT NOT NULL DEFAULT 'Saved'
				CHECK (status IN ('Saved', 'Applied', 'Interview', 'Offer', 'Rejected')),
			notes TEXT,
			job_description TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)
	`); err != nil {
		db.Close()
		return nil, fmt.Errorf("initialize jobs table: %w", err)
	}

	_, _ = db.Exec(`ALTER TABLE jobs ADD COLUMN job_description TEXT`)

	return db, nil
}
