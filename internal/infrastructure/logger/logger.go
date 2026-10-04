package logger

import (
	"context"
	"fmt"
	"marketing/internal/config"
	"marketing/pkg/contextutil"

	"go.uber.org/zap"
)

const traceID = "trace_id"

type Logger struct {
	logger *zap.Logger
}

func NewLogger(
	cfg config.LoggerConfig,
) (*Logger, error) {
	var logger *zap.Logger

	var err error

	if cfg.Env == "dev" {
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("zap development logger initialization failed: %w", err)
		}
	}

	if cfg.Env == "prod" {
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("zap production logger initialization failed: %w", err)
		}
	}

	return &Logger{
		logger: logger,
	}, nil
}

func (l *Logger) Info(
	ctx context.Context,
	msg string,
	fields ...Field,
) {
	l.logger.Info(msg, withTraceID(ctx, fields...)...)
}

func (l *Logger) Error(
	ctx context.Context,
	msg string,
	fields ...Field,
) {
	l.logger.Error(msg, withTraceID(ctx, fields...)...)
}

func (l *Logger) Debug(
	ctx context.Context,
	msg string,
	fields ...Field,
) {
	l.logger.Debug(msg, withTraceID(ctx, fields...)...)
}

func (l *Logger) Warn(
	ctx context.Context,
	msg string,
	fields ...Field,
) {
	l.logger.Warn(msg, withTraceID(ctx, fields...)...)
}

func (l *Logger) With(fields ...Field) *Logger {
	l.logger = l.logger.With(fields...)

	return l
}

func withTraceID(
	ctx context.Context,
	fields ...Field,
) []Field {
	traceIDValue := contextutil.GetTraceID(ctx)

	if traceIDValue != "" {
		return append([]Field{zap.String(traceID, traceIDValue)}, fields...)
	}

	return fields
}
