//go:build !windows

package safety

import (
	"os"
	"syscall"
)

func checkProcessAlive(p *os.Process) bool {
	err := p.Signal(syscall.Signal(0))
	if err == nil {
		return true
	}
	if err == syscall.EPERM {
		return true
	}
	return false
}
