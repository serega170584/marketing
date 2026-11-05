package telemetry

import (
	"context"
	"fmt"
	"marketing/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	serviceName string
	provider    *sdktrace.TracerProvider
	tracer      trace.Tracer
}

func NewTracer(
	ctx context.Context,
	cfg config.TracerConfig,
) (
	*Tracer,
	error,
) {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(cfg.Address),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			resource.Default().SchemaURL(),
			semconv.ServiceNameKey.String(cfg.ServiceName),
			attribute.String("environment", "development"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return &Tracer{
		serviceName: cfg.ServiceName,
		provider:    tp,
		tracer:      otel.Tracer(cfg.ServiceName),
	}, nil
}

func (tracer *Tracer) Start(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (
	context.Context,
	trace.Span,
) {
	return tracer.tracer.Start(
		ctx,
		spanName,
		opts...,
	)
}

func (tracer *Tracer) Shutdown(
	ctx context.Context,
) error {
	err := tracer.provider.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("tracer shutdown failed: %w", err)
	}

	return nil
}
