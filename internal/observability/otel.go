// Package observability provides utilities for initializing and configuring OpenTelemetry tracing
// within the application. It sets up an OTLP HTTP exporter and manages the global tracer provider,
// enabling distributed tracing and observability for services.
package observability

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Init initializes the OpenTelemetry tracer provider with an OTLP HTTP exporter.
// It sets the tracer provider globally and returns a shutdown function to clean up resources.
// If the exporter initialization fails, it logs the error and returns a no-op shutdown function.
//
// Parameters:
//
//	ctx - The context for exporter initialization.
//
// Returns:
//
//	A function that shuts down the tracer provider when called.
func Init(ctx context.Context) func(context.Context) error {
	exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint("localhost:4318"))
	if err != nil {
		log.Printf("otel exporter init: %v", err)
		return func(context.Context) error { return nil }
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp))
	otel.SetTracerProvider(tp)
	return tp.Shutdown
}
