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

type DecisionActionCount struct {
	Action string
	Count  int
}

type DecisionReasonCount struct {
	ReasonCode string
	Count      int
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

func (db *DB) GetDecisionActionCounts(repoPath, operation string) ([]DecisionActionCount, error) {
	query := `
		SELECT action, COUNT(*) as c
		FROM decision_logs
		WHERE repo_path = ?`
	args := []interface{}{repoPath}
	if operation != "" && operation != "all" {
		query += ` AND operation = ?`
		args = append(args, operation)
	}
	query += ` GROUP BY action ORDER BY c DESC`

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetDecisionActionCounts: %w", err)
	}
	defer rows.Close()

	var out []DecisionActionCount
	for rows.Next() {
		var item DecisionActionCount
		if err := rows.Scan(&item.Action, &item.Count); err != nil {
			return nil, fmt.Errorf("GetDecisionActionCounts: scan: %w", err)
		}
		out = append(out, item)
	}
	return out, nil
}

func (db *DB) GetTopDecisionReasons(repoPath, operation string, limit int) ([]DecisionReasonCount, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT reason_code, COUNT(*) as c
		FROM decision_logs
		WHERE repo_path = ?`
	args := []interface{}{repoPath}
	if operation != "" && operation != "all" {
		query += ` AND operation = ?`
		args = append(args, operation)
	}
	query += ` GROUP BY reason_code ORDER BY c DESC LIMIT ?`
	args = append(args, limit)

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetTopDecisionReasons: %w", err)
	}
	defer rows.Close()

	var out []DecisionReasonCount
	for rows.Next() {
		var item DecisionReasonCount
		if err := rows.Scan(&item.ReasonCode, &item.Count); err != nil {
			return nil, fmt.Errorf("GetTopDecisionReasons: scan: %w", err)
		}
		out = append(out, item)
	}
	return out, nil
}
