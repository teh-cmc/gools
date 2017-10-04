package tracing

import (
	"os"
	"path"
	"reflect"
	"sync"
	"time"
	"unsafe"

	ot "github.com/opentracing/opentracing-go"
	"github.com/segmentio/stats"
	jaeger "github.com/uber/jaeger-client-go"
)

// -----------------------------------------------------------------------------

// HistogramBuckets define the histogram buckets that will be used to
// compute quantiles.
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
var HistogramBuckets = []interface{}{
	0e3, 2e3, 4e3, 6e3, 8e3,
	10e3, 20e3, 30e3, 40e3, 50e3, 60e3, 70e3, 80e3, 90e3,
	100e3, 200e3, 300e3, 400e3, 500e3, 600e3, 700e3, 800e3, 900e3,
	1e6, 2e6, 3e6, 4e6, 5e6, 6e6, 7e6, 8e6, 9e6,
}

type TracingProbes struct {
	Ops struct {
		Time time.Duration `metric:"/tracing/ops/time" type:"histogram"`
		Err  string        `tag:"err"`
		Name string        `tag:"name"`
	}
}

var _probesPool = sync.Pool{New: func() interface{} { return &TracingProbes{} }}

func init() {
	for i, hb := range HistogramBuckets {
		HistogramBuckets[i] = time.Duration(hb.(float64)) * time.Microsecond
	}
	bin := path.Base(os.Args[0])
	stats.Buckets.Set(bin+":/tracing/ops/time", HistogramBuckets...)
}

// -----------------------------------------------------------------------------

const _startTimeField = "startTime" // jaeger's span private StartTime field

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
	if s, ok := span.(*jaeger.Span); ok {
		op := s.OperationName()

		prb := _probesPool.Get().(*TracingProbes)
		prb.Ops.Name = op
		if err != nil {
			prb.Ops.Err = "true"
		}
		prb.Ops.Err = "false"

		var startTime *time.Time
		startPtr := reflect.ValueOf(s).Elem().FieldByName(_startTimeField).UnsafeAddr()
		if startPtr == 0 { // should not happen, unless struct schema changes
			span.Finish()
			_probesPool.Put(prb)
			return
		}

		startTime = (*time.Time)(unsafe.Pointer(startPtr))
		prb.Ops.Time = time.Since(*startTime)
		stats.Report(prb)

		_probesPool.Put(prb)
	}
	span.Finish()
}
