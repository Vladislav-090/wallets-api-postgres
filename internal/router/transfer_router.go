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

	mux.Handle(
		"GET /transfers",
		middleware.AuthMiddleware(
			secret,
			http.HandlerFunc(transferHandler.GetTransfers),
		),
	)

	mux.Handle(
		"GET /transfers/{id}",
		middleware.AuthMiddleware(
			secret,
			http.HandlerFunc(transferHandler.GetTransferByID),
		),
	)

	mux.Handle(
		"GET /admin/transfers",
		middleware.AuthMiddleware(
			secret,
			middleware.AdminMiddleware(
				http.HandlerFunc(transferHandler.GetAllTransfers),
			),
		),
	)

	mux.Handle(
		"GET /admin/transfers/{id}",
		middleware.AuthMiddleware(
			secret,
			middleware.AdminMiddleware(
				http.HandlerFunc(transferHandler.GetTransferByIDForAdmin),
			),
		),
	)
}
