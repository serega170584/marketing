package middleware

import (
	"context"
	"marketing/internal/auth"
	"marketing/internal/infrastructure/telemetry/logger"
	"net/http"
	"strings"
)

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			requestCtx := r.Context()

			publicKey := auth.GetPublicKey()

			authHeader := r.Header.Get("Authorization")

			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 ||
				!strings.EqualFold(parts[0], "Bearer") ||
				parts[1] == "" {
				logger.Error(requestCtx, "invalid authorization header")

				http.Error(
					w,
					"Unauthorized",
					http.StatusUnauthorized,
				)

				return
			}

			claims, err := auth.ValidateJWT(
				parts[1],
				publicKey,
			)
			if err != nil {
				logger.Error(
					requestCtx,
					"invalid token",
					logger.ErrorField(err),
				)

				http.Error(
					w,
					"Invalid token",
					http.StatusUnauthorized,
				)

				return
			}

			ctx := context.WithValue(
				requestCtx,
				auth.ClaimsContextKey{},
				claims,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
