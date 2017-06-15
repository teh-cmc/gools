package grpctrace

import (
	"context"
	"os"
	"testing"

	ot "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	jaeger "github.com/uber/jaeger-client-go"
)

// -----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	reporter := jaeger.NewNullReporter()
	sampler := jaeger.NewConstSampler(false)
	t, closer := jaeger.NewTracer("grpctrace", sampler, reporter)
	ot.SetGlobalTracer(t)

	ret := m.Run()

	closer.Close()
	os.Exit(ret)
}

func TestGRPCTrace_InjectAndExtract(t *testing.T) {
	var spanAID, spanBID string

	/* client */
	ctx := func() context.Context {
		var spanA ot.Span
		ctx := context.Background()
		spanA, ctx = ot.StartSpanFromContext(ctx, "A")
		assert.NotNil(t, spanA)
		assert.IsType(t, &jaeger.Span{}, spanA)
		spanAID = spanA.Context().(jaeger.SpanContext).SpanID().String()
		return InjectSpan(ctx, spanA)
	}()

	/* server */
	func(ctx context.Context) {
		var spanB ot.Span
		spanB, ctx = ExtractSpan(ctx, "B")
		assert.NotNil(t, spanB)
		assert.IsType(t, &jaeger.Span{}, spanB)
		spanBID = spanB.Context().(jaeger.SpanContext).ParentID().String()
		assert.Equal(t, spanAID, spanBID)
	}(ctx)
}
