package store

import "fmt"

type SessionRecord struct {
	Operation   string
	SnapshotSHA string
}

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

func (db *DB) GetRecentSessions(repoPath string, limit int) ([]SessionRecord, error) {
	rows, err := db.conn.Query(`
		SELECT operation, snapshot_sha
		FROM sessions
		WHERE repo_path = ?
		ORDER BY created_at DESC, id DESC
		LIMIT ?`,
		repoPath, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("GetRecentSessions: %w", err)
	}
	defer rows.Close()

	var sessions []SessionRecord
	for rows.Next() {
		var s SessionRecord
		if err := rows.Scan(&s.Operation, &s.SnapshotSHA); err != nil {
			return nil, fmt.Errorf("GetRecentSessions: scanning row: %w", err)
		}
		sessions = append(sessions, s)
	}

	return sessions, nil
}

func (db *DB) DeleteRecentSessions(repoPath string, limit int) error {
	_, err := db.conn.Exec(`
		DELETE FROM sessions
		WHERE id IN (
			SELECT id
			FROM sessions
			WHERE repo_path = ?
			ORDER BY created_at DESC, id DESC
			LIMIT ?
		)`,
		repoPath, limit,
	)
	if err != nil {
		return fmt.Errorf("DeleteRecentSessions: %w", err)
	}
	return nil
}
