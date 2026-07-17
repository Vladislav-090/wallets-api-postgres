package router

import (
	"net/http"
	"wallets-api-postgres/internal/handlers"
)

func New() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.HealthCheck)
	return mux
}
