package tracing

import (
	"reflect"
	"time"
	"unsafe"

	ot "github.com/opentracing/opentracing-go"
	"github.com/segmentio/stats"
	jaeger "github.com/uber/jaeger-client-go"
)

// -----------------------------------------------------------------------------

// TODO(cmc): doc
// TODO(cmc): tests
func Finish(span ot.Span) {
	defer span.Finish()
	if s, ok := span.(*jaeger.Span); ok {
		op := s.OperationName()

		if _, ok := stats.DefaultEngine.HistogramBuckets()[op]; !ok {
			stats.DefaultEngine.SetHistogramBuckets(s.OperationName(),
				0, 2e3, 4e3, 6e3, 8e3,
				10e3, 20e3, 30e3, 40e3, 50e3, 60e3, 70e3, 80e3, 90e3,
				100e3, 200e3, 300e3, 400e3, 500e3, 600e3, 700e3, 800e3, 900e3,
			)
		}

		var startTime *time.Time
		startPtr := reflect.ValueOf(s).Elem().FieldByName("startTime").UnsafeAddr()
		if startPtr == 0 { // should not happen, unless struct schema changes
			return
		}

		startTime = (*time.Time)(unsafe.Pointer(startPtr))
		d := time.Since(*startTime) / time.Microsecond
		stats.Observe(s.OperationName(), float64(d))
	}
}
