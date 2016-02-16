package counter

import "sync/atomic"

// -----------------------------------------------------------------------------

// Counter implements an atomic counter with a threshold.
//
// Once it reaches its threshold, the Counter is decremented by the value of
// said threshold, and the associated set of user-defined callbacks is executed.
type Counter func() int64

// Callback is executed once a Counter reaches its threshold.
type Callback func(int64)

// NewCounter returns a new Counter with the given `start` and
// `threshold` values.
//
// Once it reaches its threshold, the Counter is decremented by the value of
// said threshold, and the associated set of user-defined callbacks is executed.
func NewCounter(start, threshold int64, callbacks ...Callback) Counter {
	if threshold < 1 {
		panic("threshold < 1")
	}

	cnt := start
	return func() int64 {
		x := atomic.AddInt64(&cnt, 1)
		if x >= threshold && x%threshold == 0 {
			atomic.AddInt64(&cnt, -threshold)
			for _, cb := range callbacks {
				cb(threshold)
			}
		}
		return x
	}
}
