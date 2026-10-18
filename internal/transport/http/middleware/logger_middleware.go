package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func LoggerMiddleware(customlogger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			customlogger.Info(
				r.Context(),
				"request completed",
				map[string]any{
					"method":   r.Method,
					"uri":      r.RequestURI,
					"duration": time.Since(t1),
					"status":   ww.Status(),
				},
			)
		})
	}
}
