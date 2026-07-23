package router

import (
	"net/http"
	"wallets-api-postgres/internal/handlers"
)

func UserRouterRegister(mux *http.ServeMux, userHandler *handlers.UserHandler) {
	mux.HandleFunc("POST /register", userHandler.CreateUser)
	mux.HandleFunc("POST /login", userHandler.Login)
}
