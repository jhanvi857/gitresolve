package store

import (
	"fmt"
)

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS _schema_version (version INTEGER PRIMARY KEY);`,
	`CREATE TABLE IF NOT EXISTS conflicts (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		repo_path   TEXT NOT NULL,
		file_path   TEXT NOT NULL,
		conflict_type TEXT NOT NULL,
		severity    TEXT NOT NULL,
		resolved_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		strategy    TEXT NOT NULL
	);`,
	`CREATE TABLE IF NOT EXISTS sessions (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		repo_path   TEXT NOT NULL,
		operation   TEXT NOT NULL,
		snapshot_sha TEXT NOT NULL,
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
	);`,
	`CREATE TABLE IF NOT EXISTS decision_logs (
		id              INTEGER PRIMARY KEY AUTOINCREMENT,
		repo_path       TEXT NOT NULL,
		file_path       TEXT NOT NULL,
		operation       TEXT NOT NULL,
		conflict_type   TEXT NOT NULL,
		severity        TEXT NOT NULL,
		action          TEXT NOT NULL,
		reason_code     TEXT NOT NULL,
		reason          TEXT NOT NULL,
		confidence      REAL NOT NULL,
		shadow          INTEGER NOT NULL DEFAULT 0,
		original_hash   TEXT NOT NULL DEFAULT '',
		simulated_hash  TEXT NOT NULL DEFAULT '',
		created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
	);`,
}

func (db *DB) migrate() error {
	var currentVersion int
	err := db.conn.QueryRow("SELECT MAX(version) FROM _schema_version").Scan(&currentVersion)
	if err != nil {
		// Table might not exist yet
		currentVersion = 0
	}

	for i, sql := range migrations {
		version := i + 1
		if version <= currentVersion {
			continue
		}

		tx, err := db.conn.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(sql); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %d failed: %w", version, err)
		}

		if _, err := tx.Exec("INSERT INTO _schema_version (version) VALUES (?)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update schema version to %d: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
