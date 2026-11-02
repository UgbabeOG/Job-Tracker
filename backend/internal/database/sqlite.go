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

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS profiles (
			id TEXT PRIMARY KEY,
			full_name TEXT,
			email TEXT,
			phone TEXT,
			linkedin_url TEXT,
			portfolio_url TEXT,
			cv_text TEXT,
			cv_path TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS job_sources (
			id TEXT PRIMARY KEY,
			job_id TEXT,
			source_url TEXT,
			source_platform TEXT,
			parsed_json TEXT,
			created_at TEXT NOT NULL,
			FOREIGN KEY(job_id) REFERENCES jobs(id) ON DELETE SET NULL
		);

		CREATE TABLE IF NOT EXISTS apply_suggestions (
			id TEXT PRIMARY KEY,
			job_id TEXT,
			suggestions_json TEXT,
			match_score INTEGER,
			created_at TEXT NOT NULL,
			FOREIGN KEY(job_id) REFERENCES jobs(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS assistant_audit (
			id TEXT PRIMARY KEY,
			action TEXT NOT NULL,
			details TEXT,
			created_at TEXT NOT NULL
		);
	`); err != nil {
		db.Close()
		return nil, fmt.Errorf("initialize assistant tables: %w", err)
	}

	return db, nil
}
