package logger

import (
	"context"
	"fmt"
	"log"
	"marketing/internal/config"
	"marketing/pkg/contextutil"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const traceID = "trace_id"

func InitLogger(cfg config.LoggerConfig) error {
	var logger *zap.Logger

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	switch cfg.Env {
	case "dev", "local":
		devConfig := zap.NewDevelopmentEncoderConfig()
		devConfig.TimeKey = "time"
		devConfig.LevelKey = "level"
		devConfig.CallerKey = "caller"
		devConfig.MessageKey = "msg"
		devConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		devConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")

		logger = zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(devConfig),
			zapcore.AddSync(os.Stdout),
			zap.DebugLevel,
		), zap.AddCaller(), zap.AddCallerSkip(1))
	case "prod":
		logger = zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.AddSync(os.Stdout),
			zap.InfoLevel,
		), zap.AddCaller(), zap.AddCallerSkip(1))
	default:
		return fmt.Errorf("invalid environment: %s", cfg.Env)
	}

	zap.ReplaceGlobals(logger)

	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Printf("invalid sync: %v", err)
		}
	}()

	return nil
}

func Info(
	ctx context.Context,
	msg string,
	fields ...Field,
) {
	zap.L().Info(msg, withTraceID(ctx, fields...)...)
}

func Error(
	ctx context.Context,
	msg string,
	fields ...Field,
) {
	zap.L().Error(msg, withTraceID(ctx, fields...)...)
}

func Debug(
	ctx context.Context,
	msg string,
	fields ...Field,
) {
	zap.L().Debug(msg, withTraceID(ctx, fields...)...)
}

func Warn(
	ctx context.Context,
	msg string,
	fields ...Field,
) {
	zap.L().Warn(msg, withTraceID(ctx, fields...)...)
}

func ErrorField(err error) zap.Field {
	return zap.Error(err)
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

func String(key string, value string) Field {
	return zap.String(key, value)
}

func Duration(key string, value time.Duration) Field {
	return zap.Duration(key, value)
}

func Int(key string, value int) Field {
	return zap.Int(key, value)
}

func Strings(key string, ss []string) Field {
	return zap.Strings(key, ss)
}
