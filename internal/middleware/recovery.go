package middleware

import (
	"log/slog"
	"net/http"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				slog.Error(
					"panic recovered",
					"error", recovered,
					"method", r.Method,
					"path", r.URL.Path,
				)

				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
