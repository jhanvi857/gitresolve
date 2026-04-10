package store

import "fmt"

type DecisionRecord struct {
	RepoPath      string
	FilePath      string
	Operation     string
	ConflictType  string
	Severity      string
	Action        string
	ReasonCode    string
	Reason        string
	Confidence    float64
	Shadow        bool
	OriginalHash  string
	SimulatedHash string
}

func (db *DB) SaveDecision(r DecisionRecord) error {
	shadow := 0
	if r.Shadow {
		shadow = 1
	}
	_, err := db.conn.Exec(`
		INSERT INTO decision_logs (
			repo_path, file_path, operation, conflict_type, severity, action,
			reason_code, reason, confidence, shadow, original_hash, simulated_hash
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.RepoPath,
		r.FilePath,
		r.Operation,
		r.ConflictType,
		r.Severity,
		r.Action,
		r.ReasonCode,
		r.Reason,
		r.Confidence,
		shadow,
		r.OriginalHash,
		r.SimulatedHash,
	)
	if err != nil {
		return fmt.Errorf("SaveDecision: %w", err)
	}
	return nil
}
