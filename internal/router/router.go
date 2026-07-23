package router

import (
	"net/http"
	"wallets-api-postgres/internal/handlers"
)

func New(userHandler *handlers.UserHandler,
	walletHandler *handlers.WalletHandler,
	secret string,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.HealthCheck)

	UserRouterRegister(mux, userHandler)
	WalletRouterRegister(mux, walletHandler, secret)

	return mux
}
