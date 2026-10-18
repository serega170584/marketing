package telemetry

import (
	"context"
	"fmt"
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
	fields map[string]any,
) {
	vs, err := mapToFields(fields)
	if err != nil {
		l.logger.Error("failed to convert to fields", zap.Error(err))

		return
	}

	l.logger.Info(msg, withTraceID(ctx, vs...)...)
}

func (l *Logger) Error(
	ctx context.Context,
	msg string,
	fields map[string]any,
) {
	vs, err := mapToFields(fields)
	if err != nil {
		l.logger.Error("failed to convert to fields", zap.Error(err))

		return
	}

	l.logger.Error(msg, withTraceID(ctx, vs...)...)
}

func (l *Logger) Debug(
	ctx context.Context,
	msg string,
	fields map[string]any,
) {
	vs, err := mapToFields(fields)
	if err != nil {
		l.logger.Error("failed to convert to fields", zap.Error(err))

		return
	}

	l.logger.Debug(msg, withTraceID(ctx, vs...)...)
}

func (l *Logger) Warn(
	ctx context.Context,
	msg string,
	fields map[string]any,
) {
	vs, err := mapToFields(fields)
	if err != nil {
		l.logger.Error("failed to convert to fields", zap.Error(err))

		return
	}

	l.logger.Warn(msg, withTraceID(ctx, vs...)...)
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

func mapToFields(m map[string]any) (
	[]Field,
	error,
) {
	fields := make([]Field, 0, len(m))
	for fk, fv := range m {
		switch v := fv.(type) {
		case string:
			fields = append(fields, zap.String(fk, v))
		case time.Duration:
			fields = append(fields, zap.Duration(fk, v))
		case int:
			fields = append(fields, zap.Int(fk, v))
		case error:
			fields = append(fields, zap.Error(v))
		default:
			return nil, fmt.Errorf("unsupported type %T", fv)
		}
	}

	return fields, nil
}
