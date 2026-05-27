package store

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

const defaultConflictRetentionCap = 1000

var retentionWarningOnce sync.Once

func conflictRetentionCap() int {
	v := strings.TrimSpace(os.Getenv("GITRESOLVE_DB_CONFLICT_CAP"))
	if v == "" {
		return defaultConflictRetentionCap
	}

	capValue, err := strconv.Atoi(v)
	if err != nil || capValue < 0 {
		return defaultConflictRetentionCap
	}

	return capValue
}

type ConflictRecord struct {
	RepoPath     string
	FilePath     string
	ConflictType string
	Severity     string
	Strategy     string
}

func (db *DB) SaveConflict(r ConflictRecord) error {
	_, err := db.conn.Exec(`
		INSERT INTO conflicts (repo_path, file_path, conflict_type, severity, strategy)
		VALUES (?, ?, ?, ?, ?)`,
		r.RepoPath, r.FilePath, r.ConflictType, r.Severity, r.Strategy,
	)
	if err != nil {
		return fmt.Errorf("SaveConflict: %w", err)
	}

	capValue := conflictRetentionCap()
	if capValue == 0 {
		return nil
	}

	// Housekeeping: Cap per-repo history to prevent unbounded growth.
	pruneResult, err := db.conn.Exec(`
		DELETE FROM conflicts 
		WHERE id IN (
			SELECT id FROM conflicts 
			WHERE repo_path = ? 
			ORDER BY resolved_at DESC 
			LIMIT -1 OFFSET ?
		)`, r.RepoPath, capValue)
	if err != nil {
		return fmt.Errorf("SaveConflict prune: %w", err)
	}

	if pruneResult != nil {
		if prunedCount, rowsErr := pruneResult.RowsAffected(); rowsErr == nil && prunedCount > 0 {
			retentionWarningOnce.Do(func() {
				fmt.Fprintf(os.Stderr, "Warning: pruned %d conflict history rows (retention cap=%d, override with GITRESOLVE_DB_CONFLICT_CAP).\n", prunedCount, capValue)
			})
		}
	}

	return nil
}

type Pattern struct {
	Label string
	Count int
}

func (db *DB) GetPatterns(repoPath string) ([]Pattern, error) {
	rows, err := db.conn.Query(`
		SELECT conflict_type, COUNT(*) as c
		FROM conflicts
		WHERE repo_path = ?
		GROUP BY conflict_type
		ORDER BY c DESC`,
		repoPath,
	)
	if err != nil {
		return nil, fmt.Errorf("GetPatterns: %w", err)
	}
	defer rows.Close()

	var patterns []Pattern
	for rows.Next() {
		var p Pattern
		if err := rows.Scan(&p.Label, &p.Count); err != nil {
			return nil, err
		}
		patterns = append(patterns, p)
	}
	return patterns, nil
}

func (db *DB) GetHistory(repoPath string) ([]ConflictRecord, error) {
	rows, err := db.conn.Query(`
		SELECT repo_path, file_path, conflict_type, severity, strategy
		FROM conflicts
		WHERE repo_path = ?
		ORDER BY resolved_at DESC`,
		repoPath,
	)
	if err != nil {
		return nil, fmt.Errorf("GetHistory: %w", err)
	}
	defer rows.Close()

	var records []ConflictRecord
	for rows.Next() {
		var r ConflictRecord
		if err := rows.Scan(&r.RepoPath, &r.FilePath, &r.ConflictType, &r.Severity, &r.Strategy); err != nil {
			return nil, fmt.Errorf("GetHistory: scanning row: %w", err)
		}
		records = append(records, r)
	}

	return records, nil
}
