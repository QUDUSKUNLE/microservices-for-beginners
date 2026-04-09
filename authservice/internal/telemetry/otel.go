package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var tracer trace.Tracer

func Init(service string) func() {
	exp, _ := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint("otel-collector:4317"),
		otlptracegrpc.WithInsecure(),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
	)

	otel.SetTracerProvider(tp)

	// IMPORTANT: initialize tracer AFTER provider
	tracer = otel.Tracer(service)

	return func() { _ = tp.Shutdown(context.Background()) }
}

func Tracer() trace.Tracer {
	return tracer
}
