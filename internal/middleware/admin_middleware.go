package middleware

import (
	"net/http"
	"wallets-api-postgres/internal/response"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetClaims(r.Context())
		if !ok {
			response.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		if claims.UserRole != "admin" {
			response.WriteError(w, http.StatusForbidden, "you must have admin role")
			return
		}
		next.ServeHTTP(w, r)
	})
}
