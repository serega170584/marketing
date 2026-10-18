package middleware

import "context"

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
