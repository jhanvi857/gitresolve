package conflict

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jhanvi857/gitresolve/pkg/logger"
	"github.com/rs/zerolog"
)

func TestLogHardening_NoContentLeak(t *testing.T) {
	var buf bytes.Buffer

	// Set log level to Warn as per requirement
	logger.InitWithLevelAndOutput(zerolog.WarnLevel, &buf)

	// Simulate a full resolve (partially)
	c := &ConflictBlock{
		OursLines:   []string{"OUR SENSITIVE DATA"},
		TheirsLines: []string{"THEIR SENSITIVE DATA"},
		FilePath:    "sensitive.go",
	}

	// Call Resolve which now contains Debug logs with hashes
	opts := ResolveOptions{NonInteractive: true}
	_, _ = Resolve(c, StrategyOurs, opts)

	// Capture output
	output := buf.String()

	// Assert no "content=" field appears (zerolog fields are usually "field":value)
	if strings.Contains(output, "\"content\"") || strings.Contains(output, "SENSITIVE") {
		t.Errorf("Security breach: sensitive content found in logs at Warn level: %s", output)
	}

	// Also check that even at Debug level, we use hashes instead of raw content
	buf.Reset()
	logger.InitWithLevelAndOutput(zerolog.DebugLevel, &buf)
	_, _ = Resolve(c, StrategyOurs, opts)
	output = buf.String()

	if strings.Contains(output, "SENSITIVE") {
		t.Errorf("Security breach: raw content found in Debug logs: %s", output)
	}
	if !strings.Contains(output, "content_hash") {
		t.Error("Expected content_hash in Debug logs, but not found")
	}
}
