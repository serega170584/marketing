package middleware

import (
	"context"
	"crypto/ecdsa"
	"marketing/internal/auth"
	"net/http"
	"strings"
)

func AuthMiddleware(
	publicKey *ecdsa.PublicKey,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			authHeader := r.Header.Get("Authorization")

			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 ||
				!strings.EqualFold(parts[0], "Bearer") ||
				parts[1] == "" {
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
				http.Error(
					w,
					"Invalid token",
					http.StatusUnauthorized,
				)

				return
			}

			ctx := context.WithValue(
				r.Context(),
				claimsContextKey{},
				claims,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type claimsContextKey struct{}

func GetClaims(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(*auth.Claims)
	return claims, ok
}
