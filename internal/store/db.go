package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

// Open opens the SQLite database and runs migrations
func Open(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, fmt.Errorf("store.Open: %w", err)
	}

	db := &DB{conn: conn}

	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("store.Open: migration: %w", err)
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) migrate() error {
	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS conflicts (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			repo_path   TEXT NOT NULL,
			file_path   TEXT NOT NULL,
			conflict_type TEXT NOT NULL,
			severity    TEXT NOT NULL,
			resolved_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			strategy    TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS sessions (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			repo_path   TEXT NOT NULL,
			operation   TEXT NOT NULL,
			snapshot_sha TEXT NOT NULL,
			created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}
