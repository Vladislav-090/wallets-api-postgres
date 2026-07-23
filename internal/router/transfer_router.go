package router

import (
	"net/http"
	"wallets-api-postgres/internal/handlers"
	"wallets-api-postgres/internal/middleware"
)

func TransferRouterRegister(mux *http.ServeMux, transferHandler *handlers.TransferHandler, secret string) {
	mux.Handle(
		"POST /transfers",
		middleware.AuthMiddleware(
			secret,
			http.HandlerFunc(transferHandler.CreateTransfer),
		),
	)
}
