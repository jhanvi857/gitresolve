package store

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

var ErrDBCorrupt = errors.New("decision log database is corrupt")

type DB struct {
	conn *sql.DB
}

// Open opens the SQLite database and runs migrations
func Open(dbPath string) (*DB, error) {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("store.Open: mkdir: %w", err)
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("store.Open: %w", err)
	}

	db := &DB{conn: conn}

	// Init sequence
	initSQL := `
		PRAGMA journal_mode=WAL;
		PRAGMA synchronous=NORMAL;
		PRAGMA foreign_keys=ON;
		PRAGMA busy_timeout=5000;
	`
	if _, err := db.conn.Exec(initSQL); err != nil {
		db.conn.Close()
		// ModernC SQLite error for corrupt/non-db file often manifests here
		if strings.Contains(err.Error(), "file is not a database") {
			return nil, fmt.Errorf("%w: %v", ErrDBCorrupt, err)
		}
		return nil, fmt.Errorf("init sequence: %w", err)
	}

	// Integrity check
	var checkResult string
	if err := db.conn.QueryRow("PRAGMA integrity_check;").Scan(&checkResult); err != nil {
		db.conn.Close()
		if strings.Contains(err.Error(), "file is not a database") {
			return nil, fmt.Errorf("%w: %v", ErrDBCorrupt, err)
		}
		return nil, fmt.Errorf("integrity check failed: %w", err)
	}
	if checkResult != "ok" {
		db.conn.Close()
		return nil, fmt.Errorf("%w: %s", ErrDBCorrupt, checkResult)
	}

	if err := db.migrate(); err != nil {
		db.conn.Close()
		return nil, fmt.Errorf("store.Open: migration: %w", err)
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}
