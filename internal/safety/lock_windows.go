//go:build windows

package safety

import (
	"os"
	"syscall"
)

func checkProcessAlive(p *os.Process) bool {
	// PROCESS_QUERY_LIMITED_INFORMATION (0x1000) is usually enough to check exit code
	h, err := syscall.OpenProcess(0x1000, false, uint32(p.Pid))
	if err != nil {
		// If we can't open the process, it might not exist
		return false
	}
	defer syscall.CloseHandle(h)

	var exitCode uint32
	err = syscall.GetExitCodeProcess(h, &exitCode)
	if err != nil {
		return false
	}
	// 259 is STILL_ACTIVE
	return exitCode == 259
}
