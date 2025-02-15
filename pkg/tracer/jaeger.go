package tracer

import (
	"context"
	"errors"

	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type JaegerProvider struct {
	tp     *sdktrace.TracerProvider
	tracer oteltrace.Tracer
}

// Jaeger Provider (HTTPS通信)
// - host: localhost:4318 (HTTPS用: portは4318)
func NewJaegerHTTPProvider(
	host, serviceName, tracerName, version string,
	sampler sdktrace.Sampler,
) (*JaegerProvider, error) {
	// define Exporter
	headers := map[string]string{
		"content-type": "application/json",
	}
	traceExporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(host), // Jaegerへのエンドポイントを設定
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, errors.New("failed to create Jaeger trace exporter")
	}

	// define tracer provider
	tp := tracerProvider(traceExporter, serviceName, version, sampler)

	// set OpenTelemetry as global telemetry
	jaegerPropagator := jaeger.Jaeger{}
	setGlobalTelemtetry(tp, jaegerPropagator)

	return &JaegerProvider{
		tp:     tp,
		tracer: tp.Tracer(tracerName),
	}, nil
}

// Jaeger Provider (for gRPC)
// - host: localhost:4317
func NewJaegerGRPCProvider(
	host, serviceName, tracerName, version string,
	sampler sdktrace.Sampler,
) (*JaegerProvider, error) {
	// define Exporter
	traceExporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithEndpoint(host), // endpoint of Jaeger
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, errors.New("failed to create Jaeger trace exporter")
	}

	// define tracer provider
	tp := tracerProvider(traceExporter, serviceName, version, sampler)

	// set OpenTelemetry as global telemetry
	jaegerPropagator := jaeger.Jaeger{}
	setGlobalTelemtetry(tp, jaegerPropagator)

	return &JaegerProvider{
		tp:     tp,
		tracer: tp.Tracer(tracerName),
	}, nil
}

func (j *JaegerProvider) NewSpan(
	ctx context.Context,
	name string,
	opts ...oteltrace.SpanStartOption,
) (context.Context, oteltrace.Span) {
	return j.tracer.Start(ctx, name, opts...)
}

// Note: must be called before main is done
func (j *JaegerProvider) Close(ctx context.Context) error {
	return j.tp.Shutdown(ctx)
}
