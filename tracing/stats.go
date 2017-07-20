package tracing

import (
	"reflect"
	"sync"
	"time"
	"unsafe"

	ot "github.com/opentracing/opentracing-go"
	"github.com/segmentio/stats"
	jaeger "github.com/uber/jaeger-client-go"
)

// -----------------------------------------------------------------------------

// Buckets define the histogram buckets that will be used to compute quantiles.
//
// Modifying this variable is safe as long as you do it before the calls to
// `Finish` start coming in.
// Bucket ranges are defined in microseconds.
//
// By default, histograms will use the following buckets:
//  0ms   2ms   4ms   6ms   8ms
//  10ms  20ms  30ms  40ms  50ms  60ms  70ms  80ms  90ms
//  100ms 200ms 300ms 400ms 500ms 600ms 700ms 800ms 900ms
//  1s    2s    3s    4s    5s    6s    7s    8s    9s
var Buckets = []float64{
	0, 2e3, 4e3, 6e3, 8e3,
	10e3, 20e3, 30e3, 40e3, 50e3, 60e3, 70e3, 80e3, 90e3,
	100e3, 200e3, 300e3, 400e3, 500e3, 600e3, 700e3, 800e3, 900e3,
	1e6, 2e6, 3e6, 4e6, 5e6, 6e6, 7e6, 8e6, 9e6,
}

// TagErr returns a metric tag that's either 'err=false' or 'err=true'
// depending on the specified error.
func TagErr(err error) stats.Tag {
	errTag := stats.Tag{Name: "err", Value: "false"}
	if err != nil {
		errTag.Value = "true"
	}
	return errTag
}

// -----------------------------------------------------------------------------

const _startTimeField = "startTime" // jaeger's span private StartTime field

var (
	_bucketedOps     = map[string]struct{}{} // already bucketed histograms
	_bucketedOpsLock = &sync.RWMutex{}
)

// Finish finishes the specified `span`.
//
// If and only if the span is implemented by a `jaeger.Span`, a metric with
// the name and duration of the span will be sent out using the 'stats' library
// from Segment ('github.com/segmentio/stats').
// If the specified error is not nil, the metric will be tagged with an
// 'err=true' KV pair.
//
// TODO(cmc): tests
func Finish(span ot.Span, err error) {
	defer span.Finish()
	if s, ok := span.(*jaeger.Span); ok {
		op := s.OperationName()

		_bucketedOpsLock.RLock()
		_, ok = _bucketedOps[op]
		_bucketedOpsLock.RUnlock()
		if !ok {
			_bucketedOpsLock.Lock()
			stats.DefaultEngine.SetHistogramBuckets(op, Buckets...)
			_bucketedOpsLock.Unlock()
		}

		var startTime *time.Time
		startPtr := reflect.ValueOf(s).Elem().FieldByName(_startTimeField).UnsafeAddr()
		if startPtr == 0 { // should not happen, unless struct schema changes
			return
		}

		startTime = (*time.Time)(unsafe.Pointer(startPtr))
		d := time.Since(*startTime) / time.Microsecond
		stats.Observe(op, float64(d), TagErr(err))
	}
}
