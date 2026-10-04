package contextutil

import "context"

type ctxKey string

const traceIDKey ctxKey = "trace_id"

func WithTraceID(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		return ctx
	}

	return context.WithValue(ctx, traceIDKey, traceID)
}

func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if id, ok := ctx.Value(traceIDKey).(string); ok {
		return id
	}

	return ""
}
