package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	NewSpan(
		ctx context.Context,
		name string,
		opts ...oteltrace.SpanStartOption,
	) (context.Context, oteltrace.Span)
	Close(ctx context.Context) error
}

//
// Tracer Utility
//

type TracerMode int

const (
	TracerModeNOOP       TracerMode = iota // no tracer
	TracerModeJaegerHTTP                   // Jaeger with HTTP
	TracerModeJaegerGRPC                   // Jaeger with gRPC
	TracerModeDataDog                      // Datadog
)

func ValidateTracerEnv(value string) TracerMode {
	switch value {
	case "jaeger_http":
		return TracerModeJaegerHTTP
	case "jaeger_grpc":
		return TracerModeJaegerGRPC
	case "datadog":
		return TracerModeDataDog
	default:
		return TracerModeNOOP
	}
}

// define tracer provider
// sampler allows nil
func tracerProvider(
	traceExporter sdktrace.SpanExporter,
	serviceName, version string,
	sampler sdktrace.Sampler,
) *sdktrace.TracerProvider {
	if sampler == nil {
		sampler = sdktrace.AlwaysSample()
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler), // set sampling
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(version)),
		),
	)
}

// set OpenTelemetry as global telemetry
func setGlobalTelemtetry(tp *sdktrace.TracerProvider, propagator propagation.TextMapPropagator) {
	otel.SetTracerProvider(tp)
	if propagator == nil {
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.Baggage{},
			propagation.TraceContext{},
		))
	} else {
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagator,
			propagation.Baggage{},
			propagation.TraceContext{},
		))
	}
}
