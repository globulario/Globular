package observability

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

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
