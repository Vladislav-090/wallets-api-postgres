package router

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"wallets-api-postgres/internal/handlers"
	"wallets-api-postgres/internal/middleware"
)

func New(userHandler *handlers.UserHandler,
	walletHandler *handlers.WalletHandler,
	transferHandler *handlers.TransferHandler,
	secret string,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.HealthCheck)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	UserRouterRegister(mux, userHandler, secret)
	WalletRouterRegister(mux, walletHandler, secret)
	TransferRouterRegister(mux, transferHandler, secret)

	return middleware.LoggingMiddleware(
		middleware.RecoveryMiddleware(mux),
	)
}
