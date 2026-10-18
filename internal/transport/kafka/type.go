package kafka

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type Logger interface {
	Debug(
		ctx context.Context,
		msg string,
		fields map[string]any,
	)
	Info(
		ctx context.Context,
		msg string,
		fields map[string]any,
	)
	Warn(
		ctx context.Context,
		msg string,
		fields map[string]any,
	)
	Error(
		ctx context.Context,
		msg string,
		fields map[string]any,
	)
}

type Tracer interface {
	Start(
		ctx context.Context,
		spanName string,
		opts ...trace.SpanStartOption,
	) (
		context.Context,
		trace.Span,
	)

	Shutdown(
		ctx context.Context,
	) error
}
