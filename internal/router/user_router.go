package router

import (
	"net/http"
	"wallets-api-postgres/internal/handlers"
	"wallets-api-postgres/internal/middleware"
)

func UserRouterRegister(mux *http.ServeMux, userHandler *handlers.UserHandler, secret string) {
	mux.HandleFunc("POST /register", userHandler.CreateUser)
	mux.HandleFunc("POST /login", userHandler.Login)
	mux.Handle(
		"GET /admin/users",
		middleware.AuthMiddleware(
			secret,
			middleware.AdminMiddleware(
				http.HandlerFunc(userHandler.GetUsers),
			),
		),
	)
	mux.Handle(
		"GET /admin/users/{id}",
		middleware.AuthMiddleware(
			secret,
			middleware.AdminMiddleware(
				http.HandlerFunc(userHandler.GetUserByID),
			),
		),
	)
}
