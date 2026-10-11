package middleware

import (
	"marketing/internal/infrastructure/telemetry"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func LoggerMiddleware(customlogger *telemetry.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			customlogger.Info(
				r.Context(),
				"request completed",
				telemetry.String("method", r.Method),
				telemetry.String("uri", r.RequestURI),
				telemetry.Duration("duration", time.Since(t1)),
				telemetry.Int("status", ww.Status()),
			)
		})
	}
}
