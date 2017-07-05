package tracing

import (
	ot "github.com/opentracing/opentracing-go"
	ot_ext "github.com/opentracing/opentracing-go/ext"
	ot_log "github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// -----------------------------------------------------------------------------

// LogError logs on error both in the specified `span` as well as through the
// given `zap.Logger`.
//
// TODO(cmc): tests
func LogError(err error,
	logger *zap.Logger, span ot.Span, fields ...ot_log.Field,
) {
	/* opentracing */
	ot_ext.Error.Set(span, true)
	span.LogFields(ot_log.Error(err))
	span.LogFields(fields...)

	/* zap */
	zapFields := make([]zapcore.Field, 0, len(fields))
	for _, f := range fields {
		switch v := f.Value().(type) {
		case error:
			zapFields = append(zapFields, zap.Error(v))
		case string:
			zapFields = append(zapFields, zap.String(f.Key(), v))
		case bool:
			zapFields = append(zapFields, zap.Bool(f.Key(), v))
		case int:
			zapFields = append(zapFields, zap.Int(f.Key(), v))
		case int32:
			zapFields = append(zapFields, zap.Int32(f.Key(), v))
		case int64:
			zapFields = append(zapFields, zap.Int64(f.Key(), v))
		case uint:
			zapFields = append(zapFields, zap.Uint(f.Key(), v))
		case uint32:
			zapFields = append(zapFields, zap.Uint32(f.Key(), v))
		case uint64:
			zapFields = append(zapFields, zap.Uint64(f.Key(), v))
		case float32:
			zapFields = append(zapFields, zap.Float32(f.Key(), v))
		case float64:
			zapFields = append(zapFields, zap.Float64(f.Key(), v))
		}
	}
	logger.Error(err.Error(), zapFields...)
}
