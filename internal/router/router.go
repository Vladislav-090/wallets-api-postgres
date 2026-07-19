package router

import (
	"net/http"
	"wallets-api-postgres/internal/handlers"
)

func New(userHandler *handlers.UserHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("/users", userHandler.CreateUser)

	return mux
}
