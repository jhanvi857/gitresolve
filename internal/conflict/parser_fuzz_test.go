package conflict

import (
	"strings"
	"testing"
)

func FuzzParseAndResolveOracle(f *testing.F) {
	f.Add("package main\n<<<<<<< HEAD\nfunc a() {}\n=======\nfunc b() {}\n>>>>>>> branch\n")
	f.Add("plain text without markers")
	f.Add("<<<<<<< HEAD\nA\n=======\nB\n>>>>>>> x\n")

	f.Fuzz(func(t *testing.T, input string) {
		content := []byte(input)
		blocks := ParseFile("fuzz.go", content)

		for _, b := range blocks {
			if b.StartIndex < 0 || b.EndIndex < b.StartIndex {
				t.Fatalf("invalid block indexes: start=%d end=%d", b.StartIndex, b.EndIndex)
			}
			if b.StartLine <= 0 || b.EndLine < b.StartLine {
				t.Fatalf("invalid block lines: start=%d end=%d", b.StartLine, b.EndLine)
			}
		}

		if len(blocks) == 0 {
			return
		}

		for _, b := range blocks {
			_, _ = Resolve(b, StrategyOurs, ResolveOptions{NonInteractive: true})
		}
		out := CompileResolution(content, blocks)

		if strings.Contains(out, "<<<<<<<") || strings.Contains(out, "=======") || strings.Contains(out, ">>>>>>>") {
			t.Fatalf("resolved output still contains markers")
		}

		inLen := len(input)
		outLen := len(out)
		if inLen > 0 {
			if outLen > inLen*3 {
				t.Fatalf("oracle size guard failed: output too large (%d > %d)", outLen, inLen*3)
			}
			if inLen > 80 && outLen < int(float64(inLen)*0.3) {
				t.Fatalf("oracle size guard failed: output too small (%d < %d)", outLen, int(float64(inLen)*0.3))
			}
		}
	})
}
