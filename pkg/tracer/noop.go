package tracer

import (
	"context"

	oteltrace "go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// disable tracer

type NoopProvider struct {
	tp     *noop.TracerProvider
	tracer oteltrace.Tracer
}

// NewNoopProvider dummyのtracer
func NewNoopProvider() *NoopProvider {
	tp := noop.NewTracerProvider()
	return &NoopProvider{
		tp:     &tp,
		tracer: tp.Tracer("noop"),
	}
}

func (n *NoopProvider) NewSpan(
	ctx context.Context,
	name string,
	_ ...oteltrace.SpanStartOption,
) (context.Context, oteltrace.Span) {
	return n.tracer.Start(ctx, name, nil)
}

func (n *NoopProvider) Close(_ context.Context) error {
	return nil
}
