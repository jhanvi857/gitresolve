package safepath

import (
	"errors"
	"runtime"
	"sync/atomic"
)

var ErrUnsafePath = errors.New("path escapes repository root")

type UnsupportedPlatformError struct {
	Platform string
}

func (e UnsupportedPlatformError) Error() string {
	platform := e.Platform
	if platform == "" {
		platform = runtime.GOOS
	}
	return "safepath is unsupported on platform " + platform
}

var ErrUnsupportedPlatform error = UnsupportedPlatformError{}

var forceUnsupported uint32

func SetForceAllowUnsupported(force bool) {
	if force {
		atomic.StoreUint32(&forceUnsupported, 1)
		return
	}
	atomic.StoreUint32(&forceUnsupported, 0)
}

func IsForceAllowed() bool {
	return atomic.LoadUint32(&forceUnsupported) == 1
}

func UnsupportedPlatformErr() error {
	return UnsupportedPlatformError{Platform: runtime.GOOS}
}
