package conflict

const DefaultMaxConflictFileBytes int64 = 10 * 1024 * 1024

type ResolverConfig struct {
	MaxFileBytes int64
}

func (c ResolverConfig) FileTooLarge(size int64) bool {
	if c.MaxFileBytes < 0 {
		return false
	}
	return size > c.MaxFileBytes
}
