package middleware

import (
	"context"
	"net/http"
	"uuid"
)

func TraceIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-Id")

		if traceID == "" {
			traceID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), traceID, traceID)

		w.Header().Set("X-Trace-Id", traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
