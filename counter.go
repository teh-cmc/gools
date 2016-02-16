package counter

import "sync/atomic"

// -----------------------------------------------------------------------------

// Counter implements an atomic counter with a threshold.
//
// Once it reaches its threshold, the Counter is decremented by the value of
// said threshold, and a set of user-defined callbacks is executed.
type Counter func() int64

// Callback is executed once a Counter reaches its threshold.
type Callback func(int64)

// NewCounter returns a new Counter with the given `start` and
// `threshold` values.
//
// Once it reaches its threshold, the Counter is decremented by the value of
// said threshold, and a set of user-defined callbacks is executed.
// Each callback is executed in its own goroutine.
func NewCounter(start, threshold int64, callbacks ...Callback) Counter {
	cnt := start
	return func() int64 {
		x := atomic.SwapInt64(&cnt, atomic.AddInt64(&cnt, 1))
		if x == threshold {
			for _, cb := range callbacks {
				go cb(x)
			}
		}
		return x
	}
}
