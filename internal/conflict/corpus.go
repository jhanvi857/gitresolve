package conflict

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func NormalizeConflictLines(lines []string) []string {
	normalized := make([]string, 0, len(lines))
	for _, line := range lines {
		normalized = append(normalized, strings.TrimSpace(line))
	}
	return normalized
}

func ConflictFingerprint(c *ConflictBlock) string {
	if c == nil {
		return ""
	}
	parts := []string{
		strings.ToLower(c.FilePath),
		strings.Join(NormalizeConflictLines(c.OursLines), "\n"),
		strings.Join(NormalizeConflictLines(c.TheirsLines), "\n"),
		strings.Join(NormalizeConflictLines(c.BaseLines), "\n"),
	}
	sum := sha256.Sum256([]byte(strings.Join(parts, "\n---\n")))
	return hex.EncodeToString(sum[:])
}

func DedupByFingerprint(blocks []*ConflictBlock) []*ConflictBlock {
	seen := make(map[string]struct{})
	result := make([]*ConflictBlock, 0, len(blocks))
	for _, b := range blocks {
		fp := ConflictFingerprint(b)
		if fp == "" {
			continue
		}
		if _, ok := seen[fp]; ok {
			continue
		}
		seen[fp] = struct{}{}
		result = append(result, b)
	}
	return result
}
