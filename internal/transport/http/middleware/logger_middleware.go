package middleware

import (
	"marketing/internal/infrastructure/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func LoggerMiddleware(customlogger *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			customlogger.Info(
				r.Context(),
				"request completed",
				logger.String("method", r.Method),
				logger.String("uri", r.RequestURI),
				logger.Duration("duration", time.Since(t1)),
				logger.Int("status", ww.Status()),
			)
		})
	}
}
