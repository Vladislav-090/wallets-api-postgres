package middleware

import (
	"context"
	"net/http"
	"strings"
	"wallets-api-postgres/internal/auth"
	"wallets-api-postgres/internal/response"
)

type contextKey string

const claimsKey contextKey = "claims"

func AuthMiddleware(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.WriteError(w, http.StatusUnauthorized, "authorization header is required")
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			response.WriteError(w, http.StatusUnauthorized, "authorization header must use bearer token")
			return
		}

		tokenString := parts[1]

		claims, err := auth.ParseToken(tokenString, secret)
		if err != nil {
			response.WriteError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func GetClaims(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*auth.Claims)
	return claims, ok

}
