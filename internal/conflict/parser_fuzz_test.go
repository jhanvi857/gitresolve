package conflict

import (
	"testing"
)

func FuzzParseFile(f *testing.F) {
	// Seed with valid conflict examples
	f.Add("<<<<<<< OURS\nhello\n=======\nworld\n>>>>>>> THEIRS")
	f.Add("<<<<<<< OURS\n{\n  \"a\": 1\n}\n=======\n{\n  \"a\": 2\n}\n>>>>>>> THEIRS")
	f.Add("<<<<<<< OURS\n||||||| BASE\nold\n=======\nnew\n>>>>>>> THEIRS")

	f.Fuzz(func(t *testing.T, content string) {
		// Just ensure it doesn't panic
		_ = ParseFile("fuzz.txt", []byte(content))
	})
}
