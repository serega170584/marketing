package telemetry

import (
	"context"
	"marketing/internal/config"
	"marketing/pkg/contextutil"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const traceID = "trace_id"

type Logger struct {
	logger *zap.Logger
}

func NewLogger(cfg config.LoggerConfig) *Logger {
	var logger *zap.Logger

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	if cfg.Env == "dev" || cfg.Env == "local" {
		devConfig := zap.NewDevelopmentEncoderConfig()
		devConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		devConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")

		logger = zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(devConfig),
			zapcore.AddSync(os.Stdout),
			zap.DebugLevel,
		), zap.AddCaller(), zap.AddCallerSkip(1))
	}

	if cfg.Env == "prod" {
		logger = zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.AddSync(os.Stdout),
			zap.InfoLevel,
		), zap.AddCaller(), zap.AddCallerSkip(1))
	}

	return &Logger{
		logger: logger,
	}
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
		return append([]Field{String(traceID, traceIDValue)}, fields...)
	}

	return fields
}

func String(key string, value string) Field {
	return zap.String(key, value)
}

func Duration(key string, value time.Duration) Field {
	return zap.Duration(key, value)
}

func Int(key string, value int) Field {
	return zap.Int(key, value)
}
