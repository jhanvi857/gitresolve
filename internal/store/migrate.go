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
	// Migration 5: co-change coupling data from git history
	`CREATE TABLE IF NOT EXISTS co_changes (
		file_a    TEXT NOT NULL,
		file_b    TEXT NOT NULL,
		count     INTEGER NOT NULL DEFAULT 0,
		strength  REAL NOT NULL DEFAULT 0.0,
		PRIMARY KEY (file_a, file_b)
	);`,
	// Migration 6: per-file author contribution data from git history
	`CREATE TABLE IF NOT EXISTS file_authors (
		file         TEXT NOT NULL,
		author_email TEXT NOT NULL,
		weight       REAL NOT NULL DEFAULT 0.0,
		last_touched DATETIME,
		PRIMARY KEY (file, author_email)
	);`,
	// Migration 7: incremental sync state for history index
	`CREATE TABLE IF NOT EXISTS history_sync_state (
		id                  INTEGER PRIMARY KEY CHECK (id = 1),
		last_processed_sha  TEXT NOT NULL,
		updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP
	);`,
	// Migration 8: escalation message + suggested command on decision logs
	`ALTER TABLE decision_logs ADD COLUMN escalation_message TEXT NOT NULL DEFAULT '';
	 ALTER TABLE decision_logs ADD COLUMN suggested_command TEXT NOT NULL DEFAULT '';`,
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
			if rollErr := tx.Rollback(); rollErr != nil {
				return fmt.Errorf("migration %d failed: %w (rollback failed: %v)", version, err, rollErr)
			}
			return fmt.Errorf("migration %d failed: %w", version, err)
		}

		if _, err := tx.Exec("INSERT INTO _schema_version (version) VALUES (?)", version); err != nil {
			if rollErr := tx.Rollback(); rollErr != nil {
				return fmt.Errorf("failed to update schema version to %d: %w (rollback failed: %v)", version, err, rollErr)
			}
			return fmt.Errorf("failed to update schema version to %d: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
