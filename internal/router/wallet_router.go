package router

import (
	"net/http"
	"wallets-api-postgres/internal/handlers"
	"wallets-api-postgres/internal/middleware"
)

func WalletRouterRegister(mux *http.ServeMux, walletHandler *handlers.WalletHandler, secret string) {
	mux.Handle(
		"POST /wallets",
		middleware.AuthMiddleware(
			secret,
			http.HandlerFunc(walletHandler.CreateWallet),
		),
	)

	mux.Handle(
		"GET /wallets",
		middleware.AuthMiddleware(
			secret,
			http.HandlerFunc(walletHandler.GetWallets),
		),
	)

	mux.Handle(
		"GET /wallets/{id}",
		middleware.AuthMiddleware(
			secret,
			http.HandlerFunc(walletHandler.GetWalletByID),
		),
	)

	mux.Handle(
		"PATCH /wallets/{id}",
		middleware.AuthMiddleware(
			secret,
			http.HandlerFunc(walletHandler.UpdateWallet),
		),
	)

	mux.Handle(
		"DELETE /wallets/{id}",
		middleware.AuthMiddleware(
			secret,
			http.HandlerFunc(walletHandler.DeleteWallet),
		),
	)
}
