package store

import "fmt"

func (db *DB) SaveSession(repoPath, operation, snapshotSHA string) error {
	_, err := db.conn.Exec(`
		INSERT INTO sessions (repo_path, operation, snapshot_sha)
		VALUES (?, ?, ?)`,
		repoPath, operation, snapshotSHA,
	)
	if err != nil {
		return fmt.Errorf("SaveSession: %w", err)
	}
	return nil
}

func (db *DB) GetLastSession(repoPath string) (string, string, error) {
	var operation, snapshotSHA string
	err := db.conn.QueryRow(`
		SELECT operation, snapshot_sha
		FROM sessions
		WHERE repo_path = ?
		ORDER BY created_at DESC
		LIMIT 1`,
		repoPath,
	).Scan(&operation, &snapshotSHA)
	if err != nil {
		return "", "", fmt.Errorf("GetLastSession: %w", err)
	}
	return operation, snapshotSHA, nil
}
