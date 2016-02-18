package errlog

import (
	"log"
	"runtime"
	"sync/atomic"
	"unsafe"
)

// -----------------------------------------------------------------------------

// ErrLogger is a global thread-safe error logger with on/off capability.
var ErrLogger ErrLog

// -----------------------------------------------------------------------------

// ErrLog implements a thrad-safe error logger with on/off capability.
type ErrLog uint32

// Toggle enables or disables logging.
func (el *ErrLog) Toggle(toggle bool) {
	el32 := unsafe.Pointer(el)
	if toggle {
		atomic.StoreUint32((*uint32)(el32), 1)
	} else {
		atomic.StoreUint32((*uint32)(el32), 0)
	}
}

// Log logs the given error iff logging is enabled.
func (el *ErrLog) Log(err error) error {
	el32 := unsafe.Pointer(el)
	if err != nil && atomic.LoadUint32((*uint32)(el32)) == 1 {
		log.Printf("%s\n", err)
	}
	return err
}

// LogST logs the given error iff logging is enabled.
//
// The first 4096 characters from the stacktrace of the calling goroutine are
// also printed with the error.
func (el *ErrLog) LogST(err error) error {
	el32 := unsafe.Pointer(el)
	if err != nil && atomic.LoadUint32((*uint32)(el32)) == 1 {
		buf := make([]byte, 4096)
		ellipses := ""
		if runtime.Stack(buf, false) >= 4096 {
			ellipses = "..."
		}
		log.Printf("%s\n%s\n%s", err, buf, ellipses)
	}

	return err
}
