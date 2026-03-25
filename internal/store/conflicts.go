package store

import "fmt"

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
	return nil
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
