package telemetry

import "go.uber.org/zap"

type Field = zap.Field

func Strings(key string, ss []string) Field {
	return zap.Strings(key, ss)
}

func Error(err error) Field {
	return zap.Error(err)
}
