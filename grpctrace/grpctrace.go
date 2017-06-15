package grpctrace

import (
	"context"

	ot "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/metadata"
)

// -----------------------------------------------------------------------------

// InjectSpan injects the specified `span` into `ctx` using the metadata
// features from the gRPC library.
//
// On success, a new context is returned with all the necessary information
// for the current trace to cross the boundaries of the remote call.
// On failure, `ctx` is returned untouched.
func InjectSpan(ctx context.Context, span ot.Span) context.Context {
	tmc := ot.TextMapCarrier{}
	if err := span.Tracer().Inject(span.Context(), ot.TextMap, tmc); err == nil {
		ctx = metadata.NewOutgoingContext(ctx, metadata.New(tmc))
	}
	return ctx
}

// ExtractSpan extracts the tracing metadata stored in the specified `ctx` and
// starts a new child span with the given `opName`.
//
// If the extraction fails for any reason, a new span is started without any
// parental links.
// The metadata stored in the context must have been set via `InjectSpan`.
func ExtractSpan(ctx context.Context, opName string) (ot.Span, context.Context) {
	var span ot.Span
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ot.StartSpanFromContext(ctx, opName)
	}
	tmc := make(ot.TextMapCarrier, md.Len())
	for k, vs := range md {
		for _, v := range vs {
			tmc[k] = v
			break
		}
	}
	sc, err := ot.GlobalTracer().Extract(ot.TextMap, tmc)
	if err != nil {
		return ot.StartSpanFromContext(ctx, opName)
	}
	span = ot.StartSpan(opName, ot.ChildOf(sc))
	return span, ot.ContextWithSpan(ctx, span)
}
