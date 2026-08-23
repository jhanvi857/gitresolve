package history

import (
	"database/sql"
	"fmt"
	"time"
)

// LoadSyncState reads the last processed commit SHA from the database.
// Returns empty string if no sync state exists yet.
func LoadSyncState(conn *sql.DB) (string, error) {
	var sha string
	err := conn.QueryRow("SELECT last_processed_sha FROM history_sync_state WHERE id = 1").Scan(&sha)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("LoadSyncState: %w", err)
	}
	return sha, nil
}

// SaveSyncState upserts the last processed commit SHA.
func SaveSyncState(conn *sql.DB, sha string) error {
	_, err := conn.Exec(`
		INSERT INTO history_sync_state (id, last_processed_sha, updated_at)
		VALUES (1, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(id) DO UPDATE SET
			last_processed_sha = excluded.last_processed_sha,
			updated_at = CURRENT_TIMESTAMP`,
		sha,
	)
	if err != nil {
		return fmt.Errorf("SaveSyncState: %w", err)
	}
	return nil
}

// PersistCoChanges upserts co-change data into the co_changes table.
func PersistCoChanges(conn *sql.DB, changes []CoChange) error {
	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("PersistCoChanges: begin tx: %w", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO co_changes (file_a, file_b, count, strength)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(file_a, file_b) DO UPDATE SET
			count = excluded.count,
			strength = excluded.strength`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("PersistCoChanges: prepare: %w", err)
	}
	defer stmt.Close()

	for _, cc := range changes {
		if _, err := stmt.Exec(cc.FileA, cc.FileB, cc.Count, cc.Strength); err != nil {
			tx.Rollback()
			return fmt.Errorf("PersistCoChanges: exec: %w", err)
		}
	}

	return tx.Commit()
}

// PersistAuthors upserts author contribution data into the file_authors table.
func PersistAuthors(conn *sql.DB, authors map[string][]AuthorContribution) error {
	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("PersistAuthors: begin tx: %w", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO file_authors (file, author_email, weight, last_touched)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(file, author_email) DO UPDATE SET
			weight = excluded.weight,
			last_touched = excluded.last_touched`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("PersistAuthors: prepare: %w", err)
	}
	defer stmt.Close()

	for file, contribs := range authors {
		for _, c := range contribs {
			if _, err := stmt.Exec(file, c.Email, c.Weight, c.LastTouched); err != nil {
				tx.Rollback()
				return fmt.Errorf("PersistAuthors: exec: %w", err)
			}
		}
	}

	return tx.Commit()
}

// LoadCoChanges reads persisted co-change data back into an in-memory map.
func LoadCoChanges(conn *sql.DB) (map[string]map[string]int, error) {
	rows, err := conn.Query("SELECT file_a, file_b, count FROM co_changes")
	if err != nil {
		return nil, fmt.Errorf("LoadCoChanges: %w", err)
	}
	defer rows.Close()

	result := make(map[string]map[string]int)
	for rows.Next() {
		var a, b string
		var count int
		if err := rows.Scan(&a, &b, &count); err != nil {
			return nil, fmt.Errorf("LoadCoChanges: scan: %w", err)
		}
		if result[a] == nil {
			result[a] = make(map[string]int)
		}
		result[a][b] = count
	}

	return result, nil
}

// LoadAuthors reads persisted author data back into an in-memory map.
func LoadAuthors(conn *sql.DB) (map[string][]AuthorContribution, error) {
	rows, err := conn.Query("SELECT file, author_email, weight, last_touched FROM file_authors")
	if err != nil {
		return nil, fmt.Errorf("LoadAuthors: %w", err)
	}
	defer rows.Close()

	result := make(map[string][]AuthorContribution)
	for rows.Next() {
		var file, email string
		var weight float64
		var lastTouched time.Time
		if err := rows.Scan(&file, &email, &weight, &lastTouched); err != nil {
			return nil, fmt.Errorf("LoadAuthors: scan: %w", err)
		}
		result[file] = append(result[file], AuthorContribution{
			Email:       email,
			Weight:      weight,
			LastTouched: lastTouched,
		})
	}

	return result, nil
}
