package conflict

import "testing"

func TestConflictFingerprint_DedupEquivalentBlocks(t *testing.T) {
	a := &ConflictBlock{
		FilePath: "pkg/file.go",
		OursLines: []string{
			" const X = 1 ",
		},
		TheirsLines: []string{
			"const X = 2",
		},
	}
	b := &ConflictBlock{
		FilePath: "PKG/file.go",
		OursLines: []string{
			"const X = 1",
		},
		TheirsLines: []string{
			" const X = 2 ",
		},
	}

	if ConflictFingerprint(a) != ConflictFingerprint(b) {
		t.Fatal("equivalent conflict blocks should produce identical fingerprints")
	}
}

func TestDedupByFingerprint(t *testing.T) {
	blocks := []*ConflictBlock{
		{FilePath: "a.go", OursLines: []string{"x"}, TheirsLines: []string{"y"}},
		{FilePath: "a.go", OursLines: []string{" x "}, TheirsLines: []string{" y "}},
		{FilePath: "b.go", OursLines: []string{"x"}, TheirsLines: []string{"z"}},
	}

	deduped := DedupByFingerprint(blocks)
	if len(deduped) != 2 {
		t.Fatalf("expected 2 deduplicated blocks, got %d", len(deduped))
	}
}
