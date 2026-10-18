package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
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
	serviceName string,
	tracer trace.Tracer,
) (*Tracer, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
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
		serviceName: serviceName,
		provider:    tp,
		tracer:      tracer,
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
